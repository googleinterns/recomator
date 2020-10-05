/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/recommender/v1"
)

type gcloudStateInfo = recommender.GoogleCloudRecommenderV1RecommendationStateInfo
type gcloudContent = recommender.GoogleCloudRecommenderV1RecommendationContent
type gcloudOperationGroup = recommender.GoogleCloudRecommenderV1OperationGroup

type mockAuth struct {
	tokens        map[string]string
	users         map[string]User
	mutex         sync.Mutex
	redirectCalls int
}

var errorAuthCodeAlreadyUsed = &googleapi.Error{Code: http.StatusBadRequest, Message: "oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\": \"invalid_grant\",\n  \"error_description\": \"Bad Request\"\n}"}

func getToken(code string) string {
	return code + "-token"
}

func (s *mockAuth) AuthCodeURL(options ...oauth2.AuthCodeOption) string {
	s.mutex.Lock()
	s.redirectCalls++
	s.mutex.Unlock()
	return ""
}

func (s *mockAuth) CreateUser(authCode string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	token, ok := s.tokens[authCode]
	if ok {
		return "", errorAuthCodeAlreadyUsed
	}
	token = getToken(authCode)
	s.tokens[authCode] = token
	user := User{email: authCode, service: &mockGoogleService{}}
	s.users[token] = user
	return token, nil
}

var errorUnauthorized = &googleapi.Error{Code: http.StatusUnauthorized, Message: "your token is invalid"}

func (s *mockAuth) Verify(token string) (string, error) {
	if strings.HasSuffix(token, "-token") {
		return token[:len(token)-len("-token")], nil
	}
	return "", errorUnauthorized
}

func (s *mockAuth) GetUser(email string) (User, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	user, ok := s.users[getToken(email)]
	return user, ok
}

type mockGoogleService struct {
	automation.GoogleService
}

var emptyRecommendation = recommender.GoogleCloudRecommenderV1Recommendation{
	Content: &gcloudContent{
		OperationGroups: []*gcloudOperationGroup{},
	},
	Etag:      "\"40204a1000e5befe\"",
	Name:      "name",
	StateInfo: &gcloudStateInfo{State: "Active"},
}

var projects = []string{"project"}

func (s *mockGoogleService) ListProjects() ([]string, error) {
	return projects, nil
}

func (s *mockGoogleService) ListZonesNames(project string) ([]string, error) {
	return []string{}, nil
}

func (s *mockGoogleService) ListRegionsNames(project string) ([]string, error) {
	return []string{}, nil
}

func (s *mockGoogleService) ListAPIRequirements(project string, apis []string) ([]*automation.Requirement, error) {
	return []*automation.Requirement{}, nil
}

func (s *mockGoogleService) ListPermissionRequirements(project string, permissions [][]string) ([]*automation.Requirement, error) {
	return []*automation.Requirement{}, nil
}

func (s *mockGoogleService) GetRecommendation(name string) (*recommender.GoogleCloudRecommenderV1Recommendation, error) {
	rec := emptyRecommendation
	rec.Name = name
	return &rec, nil
}

func (s *mockGoogleService) MarkRecommendationClaimed(name, etag string) (*recommender.GoogleCloudRecommenderV1Recommendation, error) {
	return s.GetRecommendation(name)
}

func (s *mockGoogleService) MarkRecommendationSucceeded(name, etag string) (*recommender.GoogleCloudRecommenderV1Recommendation, error) {
	return s.GetRecommendation(name)
}

func newMockShared() *SharedService {
	var service SharedService
	auth := &mockAuth{}
	auth.users = make(map[string]User)
	auth.tokens = make(map[string]string)
	service.auth = auth
	service.requests = NewRequestsMap()
	return &service
}

func TestAuth(t *testing.T) {
	router := SetUpRouter(newMockShared())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth?code=code", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var actualToken idToken
	err := json.Unmarshal(w.Body.Bytes(), &actualToken)
	assert.NoError(t, err, "This should be json")
	token := idToken{Token: getToken("code")}
	assert.Equal(t, token, actualToken)
}

func createUser(code string, router *gin.Engine) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth?code="+code, nil)
	router.ServeHTTP(w, req)
}

func newDecoder(data []byte) *json.Decoder {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec
}

func TestErrorAuth(t *testing.T) {
	router := SetUpRouter(newMockShared())
	createUser("code", router)
	req, _ := http.NewRequest("GET", "/api/auth?code=code", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, errorAuthCodeAlreadyUsed.Code, w.Code, "Another response code expected")
	var resp ErrorResponse
	err := newDecoder(w.Body.Bytes()).Decode(&resp)
	assert.NoError(t, err, "There should be ErrorResponse in body")
	assert.NotEmpty(t, resp.ErrorMessage, "ErrorMessage should not be empty")
}

func TestWrongAuthFormat(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/recommendations/apply?name=name", nil)
	req.Header.Add("Authorization", "NotBearer "+getToken(code))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Status should be BadRequest")
	var resp ErrorResponse
	err := newDecoder(w.Body.Bytes()).Decode(&resp)
	assert.NoError(t, err, "There should be ErrorResponse in body")
	assert.NotEmpty(t, resp.ErrorMessage, "ErrorMessage should not be empty")
}

func TestNoSuchUser(t *testing.T) {
	router := SetUpRouter(newMockShared())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/recommendations/apply?name=name", nil)
	req.Header.Add("Authorization", "Bearer "+getToken("code"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Status should be NotFound")
	var resp ErrorResponse
	err := newDecoder(w.Body.Bytes()).Decode(&resp)
	assert.NoError(t, err, "There should be ErrorResponse in body")
	assert.NotEmpty(t, resp.ErrorMessage, "ErrorMessage should not be empty")
}

func TestInvalidToken(t *testing.T) {
	router := SetUpRouter(newMockShared())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/recommendations/apply?name=name", nil)
	req.Header.Add("Authorization", "Bearer "+"nottoken")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Status should be Unauthorized")
	var resp ErrorResponse
	err := newDecoder(w.Body.Bytes()).Decode(&resp)
	assert.NoError(t, err, "There should be ErrorResponse in body")
	assert.NotEmpty(t, resp.ErrorMessage, "ErrorMessage should not be empty")
}

func TestRedirect(t *testing.T) {
	var service SharedService
	auth := &mockAuth{}
	service.auth = auth
	router := SetUpRouter(&service)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/redirect", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code, "Status should be SeeOther")
	assert.Equal(t, 1, auth.redirectCalls, "AuthCodeURL should be called once")
}

func TestList(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)
	w := httptest.NewRecorder()
	request := ListRequest{Projects: projects}
	bytes, err := json.Marshal(request)
	assert.NoError(t, err, "Should be no error")
	req, _ := http.NewRequest("POST", "/api/recommendations", strings.NewReader(string(bytes)))
	req.Header.Add("Authorization", "Bearer "+getToken(code))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Wrong response code")
	requestID := w.Body.String()

	for {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/recommendations?request_id="+requestID, nil)
		req.Header.Add("Authorization", "Bearer "+getToken(code))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Wrong response code")
		var resp Progress
		err := newDecoder(w.Body.Bytes()).Decode(&resp)
		if err == nil {
			assert.True(t, resp.BatchesProcessed < resp.NumberOfBatches, "Should not be done now")
			continue
		}
		var finalResp ListRecommendationsResponse
		err = newDecoder(w.Body.Bytes()).Decode(&finalResp)
		assert.NoError(t, err, "No error expected")
		assert.Equal(t, 0, len(finalResp.Recommendations), "No recommendations expected")
		break
	}
}

func checkApplySuceeded(t *testing.T, router *gin.Engine, token, recommendationName string) {
	for {
		w := httptest.NewRecorder()
		reqStatus, _ := http.NewRequest("GET", "/api/recommendations/checkStatus?name="+recommendationName, nil)
		reqStatus.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, reqStatus)
		if !assert.Equal(t, http.StatusOK, w.Code, "Wrong response code") {
			break
		}
		var resp CheckStatusResponse
		err := newDecoder(w.Body.Bytes()).Decode(&resp)
		assert.NoError(t, err, "No error expected")
		if resp.Status == succeededStatus {
			break
		}
		assert.Equal(t, inProgressStatus, resp.Status, "Should be in progress if not done")
	}
}

func TestApply(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/recommendations/apply?name=name", nil)
	req.Header.Add("Authorization", "Bearer "+getToken(code))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusCreated, w.Code, "Response code should be OK") {
		return
	}
	checkApplySuceeded(t, router, getToken(code), "name")
}

func TestCheckingStatusOnGCP(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)

	w := httptest.NewRecorder()
	reqStatus, _ := http.NewRequest("GET", "/api/recommendations/checkStatus?name=name", nil)
	reqStatus.Header.Add("Authorization", "Bearer "+getToken(code))
	router.ServeHTTP(w, reqStatus)
	if assert.Equal(t, http.StatusOK, w.Code, "Wrong response code") {
		var resp CheckStatusResponse
		err := newDecoder(w.Body.Bytes()).Decode(&resp)
		assert.NoError(t, err, "No error expected")
		assert.Equal(t, emptyRecommendation.StateInfo.State, resp.Status, "Status should be the same as in GetRecommendation")
	}
}

func TestMultipleApplies(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)

	names := []string{"name1", "name2", "name3"}
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/recommendations/apply?name="+name, nil)
			req.Header.Add("Authorization", "Bearer "+getToken(code))
			router.ServeHTTP(w, req)
			if assert.Equal(t, http.StatusCreated, w.Code, "Should be created") {
				checkApplySuceeded(t, router, getToken(code), name)
			}
		}(name)
	}

	wg.Wait()
}

func TestListingProjects(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/projects", nil)
	req.Header.Add("Authorization", "Bearer "+getToken(code))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Wrong response code")
	var resp ProjectsResponse
	err := newDecoder(w.Body.Bytes()).Decode(&resp)
	assert.NoError(t, err, "No error expected")
	assert.ElementsMatch(t, projects, resp.Projects, "Other projects expected")
}

func TestRequirements(t *testing.T) {
	code := "authcode"
	router := SetUpRouter(newMockShared())
	createUser(code, router)
	w := httptest.NewRecorder()
	request := ListRequest{Projects: projects}
	bytes, err := json.Marshal(request)
	assert.NoError(t, err, "Should be no error")
	req, _ := http.NewRequest("POST", "/api/requirements", strings.NewReader(string(bytes)))
	req.Header.Add("Authorization", "Bearer "+getToken(code))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Wrong response code")
	requestID := w.Body.String()

	for {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/requirements?request_id="+requestID, nil)
		req.Header.Add("Authorization", "Bearer "+getToken(code))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Wrong response code")
		var resp Progress
		err := newDecoder(w.Body.Bytes()).Decode(&resp)
		if err == nil {
			assert.True(t, resp.BatchesProcessed < resp.NumberOfBatches, "Should not be done now")
			continue
		}
		var finalResp CheckRequirementsResponse
		err = newDecoder(w.Body.Bytes()).Decode(&finalResp)
		assert.NoError(t, err, "No error expected")
		assert.Equal(t, len(projects), len(finalResp.ProjectsRequirements), "Requirements for all projects expected")
		break
	}
}
