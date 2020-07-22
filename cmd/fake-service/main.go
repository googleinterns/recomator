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

	"google.golang.org/api/recommender/v1"

	"github.com/gin-gonic/gin"
)

const (
	numberOfFakeRecommendations = 60
	expectedListTime            = time.Second * 10
	expectedApplyTime           = time.Second * 2
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

func (s *mockListService) GetResult() ([]*gcloudRecommendation, string) {
	return s.recommendations, s.anotherPageToken
}

const (
	notAppliedStatus = "NOT APPLIED"
	inProgressStatus = "IN PROGRESS"
	failedStatus     = "FAILED"
	succeededStatus  = "SUCCEEDED"
)

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
	for s.callsDone < s.numberOfCalls {
		sleep := randomTime(expectedApplyTime) / time.Duration(s.numberOfCalls)
		time.Sleep(sleep)
		s.mutex.Lock()
		if rand.Int()%10 == 0 {
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

func main() {
	cachedCalls := recommendationsMap{data: make(map[string][]*gcloudRecommendation)} // the key is anotherPageToken
	listRequestsInProcess := listRequestsMap{data: make(map[string]*mockListService)} // the key is AccessToken, but in this version, token is always ""
	bufferSize := 100
	newListRequests := make(chan *mockListService, bufferSize)

	applyRequestsInProcess := applyRequestsMap{data: make(map[string]*mockApplyService)} // the key is recommendation name
	newApplyRequests := make(chan *mockApplyService, bufferSize)

	numWorkers := 10
	for i := 0; i < numWorkers; i++ { // goroutines proccessing requests in background
		go func() {
			for {
				select {
				case s := <-newListRequests:
					s.ListRecommendations()
				case s := <-newApplyRequests:
					s.Apply()
				default:
					break
				}
			}
		}()
	}

	recommendations := getFakeRecommendations()

	router := gin.Default()

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
			newListRequests <- service
		}
		done, all := service.GetProgress()
		if done < all {
			c.JSON(http.StatusOK, ListRecommendationsProgressResponse{done, all})
		} else {
			result, anotherPageToken := service.GetResult()
			listRequestsInProcess.Delete(token)
			cachedCalls.Store(anotherPageToken, result)
			c.JSON(http.StatusOK, NewListRecommendationsResponse(
				anotherPageToken, pageIndex, pageSize, result))
		}
		return
	})

	router.POST("/recommendations/:name/apply", func(c *gin.Context) {
		name := c.Param("name")
		service, loaded := applyRequestsInProcess.LoadOrStore(name, &mockApplyService{name: name, callsDone: 0, numberOfCalls: 1})
		if !loaded {
			newApplyRequests <- service
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
		c.JSON(http.StatusOK, gin.H{"status": status})
	})

	router.Run(":8080")

}