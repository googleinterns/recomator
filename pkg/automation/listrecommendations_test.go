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

package automation

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockService struct {
	GoogleService
	mutex                                 sync.Mutex
	numberOfTimesListRecommendationsCalls int
	zones                                 []string
	zonesCalled                           []string
}

func (s *MockService) ListRecommendations(project string, location string, recommenderID string) ([]*gcloudRecommendation, error) {
	s.mutex.Lock()
	s.numberOfTimesListRecommendationsCalls++
	s.zonesCalled = append(s.zonesCalled, location)
	s.mutex.Unlock()
	return []*gcloudRecommendation{}, nil
}

func (s *MockService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func TestListRecommendations(t *testing.T) {
	for numConcurrentCalls := 0; numConcurrentCalls <= 5; numConcurrentCalls++ {
		zones := []string{"zone1", "zone2", "zone3"}
		mock := &MockService{zones: zones}
		result, err := ListRecommendations(mock, "", "", numConcurrentCalls)

		if assert.NoError(t, err, "Unexpected error from ListRecommendations") {
			assert.Equal(t, 0, len(result), "No recommendations expected")
			assert.Equal(t, len(zones), mock.numberOfTimesListRecommendationsCalls, "Wrong number of ListRecommendations calls")
			assert.ElementsMatch(t, mock.zones, mock.zonesCalled, "ListRecommendations was called for different zones")
		}
	}
}

type ErrorZonesService struct {
	GoogleService
	err error
}

func (s *ErrorZonesService) ListZonesNames(project string) ([]string, error) {
	return []string{}, s.err
}

func TestErrorInListZones(t *testing.T) {
	errorMessage := "error listing zones"
	_, err := ListRecommendations(&ErrorZonesService{err: fmt.Errorf(errorMessage)}, "", "", 2)
	assert.EqualError(t, err, errorMessage, "Expected error calling ListZones")
}

type ErrorRecommendationService struct {
	GoogleService
	err       error
	errorZone string
	mutex     sync.Mutex
	zones     []string
}

func (s *ErrorRecommendationService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func (s *ErrorRecommendationService) ListRecommendations(project string, location string, recommenderID string) ([]*gcloudRecommendation, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if location == s.errorZone {
		return []*gcloudRecommendation{}, s.err
	}
	return []*gcloudRecommendation{}, nil
}

func TestErrorInRecommendations(t *testing.T) {
	errorMessage := "error listing recommendations"
	zones := []string{}
	for i := 1; i < 5; i++ {
		zones = append(zones, fmt.Sprintf("zone %d", i))
	}

	for _, zone := range zones {
		for numConcurrentCalls := 1; numConcurrentCalls <= 10; numConcurrentCalls++ {
			_, err := ListRecommendations(
				&ErrorRecommendationService{
					err:       fmt.Errorf(errorMessage),
					zones:     zones,
					errorZone: zone}, "", "", numConcurrentCalls)
			assert.EqualError(t, err, errorMessage, "Expected error calling ListRecommendations")
		}
	}
}

type BenchmarkService struct {
	GoogleService
}

func (s *BenchmarkService) ListRecommendations(project string, location string, recommenderID string) ([]*gcloudRecommendation, error) {
	time.Sleep(time.Millisecond * 100)
	return []*gcloudRecommendation{}, nil
}

func (s *BenchmarkService) ListZonesNames(project string) ([]string, error) {
	zones := []string{}
	for i := 0; i < 100; i++ {
		zones = append(zones, fmt.Sprintf("zone %d", i))
	}
	return zones, nil
}

func BenchmarkGoroutines(b *testing.B) {
	for _, numConcurrentCalls := range []int{4, 8, 16, 32, 64} {
		b.Run(fmt.Sprintf("%d goroutines:", numConcurrentCalls), func(b *testing.B) {
			s := &BenchmarkService{}
			ListRecommendations(s, "", "", numConcurrentCalls)
		})
	}
}
