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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/segmentio/ksuid"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/recommender/v1"

	"github.com/gin-gonic/gin"
)

const (
	numberOfFakeRecommendations = 60
	expectedListTime            = time.Second * 10
	expectedApplyTime           = time.Second * 10
	probablityOfErrorList       = 1.0
	probablityOfErrorApply      = 0.3
)

type gcloudRecommendation recommender.GoogleCloudRecommenderV1Recommendation

// ListRecommendationsResponse is response to list/recommendations method
type ListRecommendationsResponse struct {
	AnotherPageToken string                  `json:"anotherPageToken"`
	NumberOfPages    int                     `json:"numberOfPages"`
	PageIndex        int                     `json:"pageIndex"`
	PageSize         int                     `json:"pageSize"`
	Recommendations  []*gcloudRecommendation `json:"recommendations"`
}

// NewListRecommendationsResponse creates new ListRecommendationsResponse
func NewListRecommendationsResponse(anotherPageToken string, pageIndex, pageSize int, recommendations []*gcloudRecommendation) ListRecommendationsResponse {
	numberOfPages := (len(recommendations) + pageSize - 1) / pageSize
	var page []*gcloudRecommendation
	if pageIndex < numberOfPages {
		end := (pageIndex + 1) * pageSize
		if end > len(recommendations) {
			end = len(recommendations)
		}
		page = recommendations[pageIndex*pageSize : end]
	}

	return ListRecommendationsResponse{
		AnotherPageToken: anotherPageToken,
		NumberOfPages:    numberOfPages,
		PageIndex:        pageIndex,
		PageSize:         pageSize,
		Recommendations:  page,
	}
}

// ListRecommendationsProgressResponse is response to list/recommendations method
// if not all recommendations have been processed yet.
type ListRecommendationsProgressResponse struct {
	BatchesProcessed int `json:"batchesProcessed"`
	NumberOfBatches  int `json:"numberOfBatches"`
}

type mockListService struct {
	anotherPageToken string
	recommendations  []*gcloudRecommendation
	callsDone        int
	numberOfCalls    int
	token            string
	mutex            sync.Mutex
	err              error
}

type listResult struct {
	anotherPageToken string
	recommendations  []*gcloudRecommendation
}

func randomTime(average time.Duration) time.Duration {
	ns := average.Nanoseconds()
	return time.Duration(float64(ns)*rand.Float64()*2) * time.Nanosecond
}

func (s *mockListService) ListRecommendations() {
	s.mutex.Lock()
	s.numberOfCalls = rand.Int() % 100
	s.callsDone = 0
	s.anotherPageToken = ksuid.New().String()
	willFail := rand.Float64() <= probablityOfErrorList
	if willFail {
		s.err = &googleapi.Error{Code: 403,
			Message: "Request is missing required authentication credential. Expected OAuth 2 access token, login cookie or other valid authentication credential. See https://developers.google.com/identity/sign-in/web/devconsole-project.",
		}
	}
	s.mutex.Unlock()
	for s.callsDone < s.numberOfCalls {
		sleep := randomTime(expectedListTime) / time.Duration(s.numberOfCalls)
		time.Sleep(sleep)
		s.mutex.Lock()
		s.callsDone++
		s.mutex.Unlock()
	}
}

func (s *mockListService) GetProgress() (int, int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.callsDone, s.numberOfCalls
}

func (s *mockListService) GetResult() (*listResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &listResult{anotherPageToken: s.anotherPageToken, recommendations: s.recommendations}, nil
}

const (
	notAppliedStatus = "NOT APPLIED"
	inProgressStatus = "IN PROGRESS"
	failedStatus     = "FAILED"
	succeededStatus  = "SUCCEEDED"
)

// CheckStatusResponse is the response to recommendations/name/checkStatus method
type CheckStatusResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type mockApplyService struct {
	name          string
	err           error
	callsDone     int
	numberOfCalls int
	mutex         sync.Mutex
}

func (s *mockApplyService) Apply() {
	s.mutex.Lock()
	s.numberOfCalls = rand.Int()%2 + 2
	s.callsDone = 0
	s.mutex.Unlock()
	willFail := rand.Float64() <= probablityOfErrorApply
	for s.callsDone < s.numberOfCalls {
		sleep := randomTime(expectedApplyTime) / time.Duration(s.numberOfCalls)
		time.Sleep(sleep)
		s.mutex.Lock()
		if willFail && rand.Int()%(s.numberOfCalls-s.callsDone) == 0 {
			s.err = fmt.Errorf("applying recommendation failed: error happened on step %d", s.callsDone)
			s.callsDone = s.numberOfCalls
		} else {
			s.callsDone++
		}
		s.mutex.Unlock()
	}
}

func (s *mockApplyService) GetStatus() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.callsDone < s.numberOfCalls {
		return inProgressStatus
	}
	if s.err != nil {
		return failedStatus
	}
	return succeededStatus
}

func (s *mockApplyService) GetErrorMessage() string {
	return s.err.Error()
}

func getFakeRecommendations() []*gcloudRecommendation {
	var result []*gcloudRecommendation
	for _, rec := range recommendationsJSON {
		var recommendation gcloudRecommendation
		err := json.Unmarshal([]byte(rec), &recommendation)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, &recommendation)
	}

	uniqueRecommendations := result
	for i := 0; len(result) < numberOfFakeRecommendations; i++ {
		// add some almost indentical recommendations to get more recommendations
		for _, rec := range uniqueRecommendations {
			var newRec gcloudRecommendation
			copier.Copy(&newRec, rec)
			newRec.Name += fmt.Sprintf("-%d", i)
			result = append(result, &newRec)
		}
	}
	return result[:numberOfFakeRecommendations]
}

type recommendationsMap struct {
	data  map[string][]*gcloudRecommendation
	mutex sync.Mutex
}

func (m *recommendationsMap) Store(pageToken string, recommendations []*gcloudRecommendation) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[pageToken] = recommendations

}

func (m *recommendationsMap) Load(pageToken string) ([]*gcloudRecommendation, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	res, ok := m.data[pageToken]
	return res, ok
}

type listRequestsMap struct {
	data  map[string]*mockListService
	mutex sync.Mutex
}

func (m *listRequestsMap) LoadOrStore(name string, service *mockListService) (*mockListService, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	s, ok := m.data[name]
	if !ok {
		m.data[name] = service
		s = service
	}
	return s, ok
}

func (m *listRequestsMap) Delete(name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, name)
}

type applyRequestsMap struct {
	data  map[string]*mockApplyService
	mutex sync.Mutex
}

func (m *applyRequestsMap) Load(name string) (*mockApplyService, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	s, ok := m.data[name]
	return s, ok
}

func (m *applyRequestsMap) LoadOrStore(name string, service *mockApplyService) (*mockApplyService, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	s, ok := m.data[name]
	ok = ok && s.GetStatus() != failedStatus
	if !ok {
		m.data[name] = service
		s = service
	}
	return s, ok
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	cachedCalls := recommendationsMap{data: make(map[string][]*gcloudRecommendation)} // the key is anotherPageToken
	listRequestsInProcess := listRequestsMap{data: make(map[string]*mockListService)} // the key is AccessToken, but in this version, token is always ""

	applyRequestsInProcess := applyRequestsMap{data: make(map[string]*mockApplyService)} // the key is recommendation names

	recommendations := getFakeRecommendations()

	router := gin.Default()
	router.Use(corsMiddleware())

	router.GET("/recommendations", func(c *gin.Context) {
		defaultPageSize := 50
		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprint(defaultPageSize)))
		if err != nil {
			pageSize = defaultPageSize
		}
		defaultPageIndex := 0
		pageIndex, err := strconv.Atoi(c.DefaultQuery("pageIndex", fmt.Sprint(defaultPageIndex)))
		if err != nil {
			pageIndex = defaultPageIndex
		}
		anotherPageToken := c.Query("pageToken")
		if anotherPageToken != "" {
			result, ok := cachedCalls.Load(anotherPageToken)
			if ok {
				c.JSON(http.StatusOK, NewListRecommendationsResponse(
					anotherPageToken, pageIndex, pageSize, result))
				return
			}
		}

		token := "" // no authentication in this fake service
		service, loaded := listRequestsInProcess.LoadOrStore(token, &mockListService{token: token, recommendations: recommendations, callsDone: 0, numberOfCalls: 1})
		if !loaded {
			go service.ListRecommendations()
		}
		done, all := service.GetProgress()
		if done < all {
			c.JSON(http.StatusOK, ListRecommendationsProgressResponse{done, all})
		} else {
			listRequestsInProcess.Delete(token)
			result, err := service.GetResult()
			if err != nil {
				apiErr := err.(*googleapi.Error)
				c.String(apiErr.Code, apiErr.Message)
				return
			}
			cachedCalls.Store(result.anotherPageToken, result.recommendations)
			c.JSON(http.StatusOK, NewListRecommendationsResponse(
				result.anotherPageToken, pageIndex, pageSize, result.recommendations))
		}
		return
	})

	router.POST("/recommendations/:name/apply", func(c *gin.Context) {
		name := c.Param("name")
		service, loaded := applyRequestsInProcess.LoadOrStore(name, &mockApplyService{name: name, callsDone: 0, numberOfCalls: 1})
		if !loaded {
			go service.Apply()
		}
		return
	})

	router.GET("/recommendations/:name/checkStatus", func(c *gin.Context) {
		name := c.Param("name")
		service, ok := applyRequestsInProcess.Load(name)
		status := notAppliedStatus
		if ok {
			status = service.GetStatus()
		}
		response := CheckStatusResponse{Status: status}
		if status == failedStatus {
			response.ErrorMessage = service.GetErrorMessage()
		}
		c.JSON(http.StatusOK, response)
	})

	router.Run(":8080")

}
