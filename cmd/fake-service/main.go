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

type mockService struct {
	anotherPageToken string
	recommendations  []*gcloudRecommendation
	callsDone        int
	numberOfCalls    int
	mutex            sync.Mutex
}

func (s *mockService) ListRecommendations() {
	s.mutex.Lock()
	s.numberOfCalls = rand.Int() % 100
	s.callsDone = 0
	s.anotherPageToken = ksuid.New().String()
	s.mutex.Unlock()
	for i := 0; i < s.numberOfCalls; i++ {
		time.Sleep(time.Duration(rand.Int()%200) * time.Millisecond)
		s.mutex.Lock()
		s.callsDone++
		s.mutex.Unlock()
	}
}

func (s *mockService) GetProgress() (int, int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.callsDone, s.numberOfCalls
}

func (s *mockService) GetResult() ([]*gcloudRecommendation, string) {
	return s.recommendations, s.anotherPageToken
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
	for i := 1; i < 10; i++ {
		// add some almost indentical recommendations to get more recommendations
		for _, rec := range uniqueRecommendations {
			var newRec gcloudRecommendation
			copier.Copy(&newRec, rec)
			newRec.Name += fmt.Sprintf("-%d", i)
			result = append(result, &newRec)
		}
	}
	return result
}

type request struct {
	AccessToken string
	service     *mockService
}

func main() {
	var cachedCallsMutex sync.Mutex
	cachedCalls := make(map[string][]*gcloudRecommendation) // the key is anotherPageToken
	var requestsInProcessMutex sync.Mutex
	requestsInProcess := make(map[string]*mockService) // the key is AccessToken, but in this version, token is always ""
	newRequests := make(chan request)

	go func() { // one goroutine proccessing ListRecommendations in background
		for {
			select {
			case r := <-newRequests:
				requestsInProcessMutex.Lock()
				requestsInProcess[r.AccessToken] = r.service
				requestsInProcessMutex.Unlock()
				r.service.ListRecommendations()
			default:
				break
			}
		}
	}()

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
			cachedCallsMutex.Lock()
			defer cachedCallsMutex.Unlock()
			result, ok := cachedCalls[anotherPageToken]
			if ok {
				c.JSON(http.StatusOK, NewListRecommendationsResponse(
					anotherPageToken, pageIndex, pageSize, result))
				return
			}
		}

		token := "" // no authentication in this fake service
		requestsInProcessMutex.Lock()
		defer requestsInProcessMutex.Unlock()
		service, ok := requestsInProcess[token]
		if !ok {
			service = &mockService{recommendations: recommendations, callsDone: 0, numberOfCalls: 1}
			newRequests <- request{token, service}
		}
		done, all := service.GetProgress()
		if done < all {
			c.JSON(http.StatusOK, ListRecommendationsProgressResponse{done, all})
		} else {
			result, anotherPageToken := service.GetResult()
			delete(requestsInProcess, token)
			cachedCallsMutex.Lock()
			cachedCalls[anotherPageToken] = result
			cachedCallsMutex.Unlock()
			c.JSON(http.StatusOK, NewListRecommendationsResponse(
				anotherPageToken, pageIndex, pageSize, result))
		}
		return
	})
	router.Run(":8080")

}
