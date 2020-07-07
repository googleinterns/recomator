package automation

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/recommender/v1"
)

type MockService struct {
	zones                                 []string
	mutex                                 sync.Mutex
	numberOfTimesListRecommendationsCalls int
	zonesCalled                           []string
	numConcurrentCalls                    int
}

func (s *MockService) listZoneRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	s.mutex.Lock()
	s.numberOfTimesListRecommendationsCalls++
	s.zonesCalled = append(s.zonesCalled, location)
	s.mutex.Unlock()
	return []*recommender.GoogleCloudRecommenderV1Recommendation{}
}

func (s *MockService) listZonesNames(project string) []string {
	return s.zones
}

func (s *MockService) getNumConcurrentCalls() int {
	return s.numConcurrentCalls
}

func TestListRecommendations(t *testing.T) {
	for numConcurrentCalls := 1; numConcurrentCalls <= 5; numConcurrentCalls++ {
		zones := []string{"zone1", "zone2", "zone3"}
		mock := &MockService{zones: zones, numConcurrentCalls: numConcurrentCalls}
		result := ListRecommendations(mock, "", "")

		assert.Equal(t, 0, len(result), "No recommedations expected")
		assert.Equal(t, len(zones), mock.numberOfTimesListRecommendationsCalls, "Wrong number of listZoneRecommendaions calls")
		assert.ElementsMatch(t, mock.zones, mock.zonesCalled, "listZoneRecommendations was called for different zones")
	}
}

type SleepingService struct {
	zones                                 []string
	numberOfTimesListRecommendationsCalls int
	requestsDurations                     []time.Duration
	numConcurrentCalls                    int
}

func (s *SleepingService) listZoneRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	d, err := time.ParseDuration(location)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(d)
	s.requestsDurations = append(s.requestsDurations, d)
	return []*recommender.GoogleCloudRecommenderV1Recommendation{}
}

func (s *SleepingService) listZonesNames(project string) []string {
	return s.zones
}

func (s *SleepingService) getNumConcurrentCalls() int {
	return s.numConcurrentCalls
}

func TestConcurrency(t *testing.T) {
	durations := []time.Duration{time.Second, time.Millisecond * 50}
	zones := []string{}
	for _, duration := range durations {
		zones = append(zones, fmt.Sprint(duration))
	}
	mock := &SleepingService{zones: zones, numConcurrentCalls: 2}
	result := ListRecommendations(mock, "", "")

	assert.Equal(t, 0, len(result), "No recommedations expected")
	assert.Equal(t, len(zones), len(mock.requestsDurations), "Wrong number of calls to listZoneRecommendations")
	assert.True(t, sort.SliceIsSorted(mock.requestsDurations, func(i, j int) bool {
		return mock.requestsDurations[i] < mock.requestsDurations[j]
	}), "Faster goroutine should have finished first")
}

type BenchmarkService struct {
	numConcurrentCalls int
}

func (s *BenchmarkService) listZoneRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	time.Sleep(time.Second)
	return []*recommender.GoogleCloudRecommenderV1Recommendation{}
}

func (s *BenchmarkService) listZonesNames(project string) []string {
	zones := []string{}
	for i := 0; i < 100; i++ {
		zones = append(zones, fmt.Sprintf("zone %d", i))
	}
	return zones
}

func (s *BenchmarkService) getNumConcurrentCalls() int {
	return s.numConcurrentCalls
}

func BenchmarkGoroutines(b *testing.B) {
	for _, numConcurrentCalls := range []int{4, 8, 16, 32, 64} {
		b.Run(fmt.Sprintf("%d goroutines:", numConcurrentCalls), func(b *testing.B) {
			s := &BenchmarkService{numConcurrentCalls}
			ListRecommendations(s, "", "")
		})
	}
}
