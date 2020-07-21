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

type mockListService struct {
	anotherPageToken string
	recommendations  []*gcloudRecommendation
	callsDone        int
	numberOfCalls    int
	token            string
	mutex            sync.Mutex
}

func (s *mockListService) ListRecommendations() {
	s.mutex.Lock()
	s.numberOfCalls = rand.Int() % 100
	s.callsDone = 0
	s.anotherPageToken = ksuid.New().String()
	s.mutex.Unlock()
	for s.callsDone < s.numberOfCalls {
		time.Sleep(time.Duration(rand.Int()%200) * time.Millisecond)
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
		time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)
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
		return "IN PROGRESS"
	}
	if s.err != nil {
		return "FAILED"
	}
	return "APPLIED"
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

func main() {
	var cachedCallsMutex sync.Mutex
	cachedCalls := make(map[string][]*gcloudRecommendation) // the key is anotherPageToken
	var listRequestsMutex sync.Mutex
	listRequestsInProcess := make(map[string]*mockListService) // the key is AccessToken, but in this version, token is always ""
	newListRequests := make(chan *mockListService)

	var applyRequestsMutex sync.Mutex
	applyRequestsInProcess := make(map[string]*mockApplyService) // the key is recommendation name
	newApplyRequests := make(chan *mockApplyService)

	go func() { // one goroutine proccessing requests in background
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
			result, ok := cachedCalls[anotherPageToken]
			cachedCallsMutex.Unlock()
			if ok {
				c.JSON(http.StatusOK, NewListRecommendationsResponse(
					anotherPageToken, pageIndex, pageSize, result))
				return
			}
		}

		token := "" // no authentication in this fake service
		listRequestsMutex.Lock()
		service, ok := listRequestsInProcess[token]
		if !ok {
			service = &mockListService{token: token, recommendations: recommendations, callsDone: 0, numberOfCalls: 1}
			newListRequests <- service
			listRequestsInProcess[token] = service
		}
		listRequestsMutex.Unlock()
		done, all := service.GetProgress()
		if done < all {
			c.JSON(http.StatusOK, ListRecommendationsProgressResponse{done, all})
		} else {
			result, anotherPageToken := service.GetResult()
			listRequestsMutex.Lock()
			delete(listRequestsInProcess, token)
			listRequestsMutex.Unlock()
			cachedCallsMutex.Lock()
			cachedCalls[anotherPageToken] = result
			cachedCallsMutex.Unlock()
			c.JSON(http.StatusOK, NewListRecommendationsResponse(
				anotherPageToken, pageIndex, pageSize, result))
		}
		return
	})

	router.POST("/recommendations/:name/apply", func(c *gin.Context) {
		name := c.Param("name")
		applyRequestsMutex.Lock()
		defer applyRequestsMutex.Unlock()
		service, ok := applyRequestsInProcess[name]
		if ok {
			status := service.GetStatus()
			if status != "RUNNING" && status != "IN PROGRESS" {
				ok = false
			}
		}
		if !ok {
			service := &mockApplyService{name: name, callsDone: 0, numberOfCalls: 1}
			newApplyRequests <- service
			applyRequestsInProcess[name] = service
		}
		return
	})

	router.GET("/recommendations/:name/checkStatus", func(c *gin.Context) {
		name := c.Param("name")
		applyRequestsMutex.Lock()
		service, ok := applyRequestsInProcess[name]
		applyRequestsMutex.Unlock()
		status := "NOT APPLIED"
		if ok {
			status = service.GetStatus()
		}
		c.JSON(http.StatusOK, gin.H{"status": status})
	})

	router.Run(":8080")

}
