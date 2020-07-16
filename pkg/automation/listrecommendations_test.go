package automation

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/recommender/v1"
)

type MockService struct {
	GoogleService
	mutex                                 sync.Mutex
	numberOfTimesListRecommendationsCalls int
	zones                                 []string
	zonesCalled                           []string
}

func (s *MockService) ListRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	s.mutex.Lock()
	s.numberOfTimesListRecommendationsCalls++
	s.zonesCalled = append(s.zonesCalled, location)
	s.mutex.Unlock()
	return []*recommender.GoogleCloudRecommenderV1Recommendation{}
}

func (s *MockService) ListZonesNames(project string) []string {
	return s.zones
}

func TestListRecommendations(t *testing.T) {
	for numConcurrentCalls := 0; numConcurrentCalls <= 5; numConcurrentCalls++ {
		zones := []string{"zone1", "zone2", "zone3"}
		mock := &MockService{zones: zones}
		result := ListRecommendations(mock, "", "", numConcurrentCalls)

		assert.Equal(t, 0, len(result), "No recommedations expected")
		assert.Equal(t, len(zones), mock.numberOfTimesListRecommendationsCalls, "Wrong number of listZoneRecommendaions calls")
		assert.ElementsMatch(t, mock.zones, mock.zonesCalled, "listZoneRecommendations was called for different zones")
	}
}

type BenchmarkService struct {
	GoogleService
}

func (s *BenchmarkService) ListRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	time.Sleep(time.Second)
	return []*recommender.GoogleCloudRecommenderV1Recommendation{}
}

func (s *BenchmarkService) ListZonesNames(project string) []string {
	zones := []string{}
	for i := 0; i < 100; i++ {
		zones = append(zones, fmt.Sprintf("zone %d", i))
	}
	return zones
}

func BenchmarkGoroutines(b *testing.B) {
	for _, numConcurrentCalls := range []int{4, 8, 16, 32, 64} {
		b.Run(fmt.Sprintf("%d goroutines:", numConcurrentCalls), func(b *testing.B) {
			s := &BenchmarkService{}
			ListRecommendations(s, "", "", numConcurrentCalls)
		})
	}
}
