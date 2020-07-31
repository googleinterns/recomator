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
	regions								  []string
	locationsCalled                       []string
}

func (s *MockService) ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error) {
	s.mutex.Lock()
	s.numberOfTimesListRecommendationsCalls++
	s.locationsCalled = append(s.locationsCalled, location)
	s.mutex.Unlock()
	return []*gcloudRecommendation{}, nil
}

func (s *MockService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func (s *MockService) ListRegionsNames(project string) ([]string, error) {
	return s.regions, nil
}

func TestListRecommendations(t *testing.T) {
	for numConcurrentCalls := 0; numConcurrentCalls <= 5; numConcurrentCalls++ {
		zones := []string{"zone1", "zone2", "zone3"}
		regions := []string{"region1", "region2", "region3"}
		mock := &MockService{zones: zones, regions: regions}
		result, err := ListRecommendations(mock, "", "", numConcurrentCalls)

		if assert.NoError(t, err, "Unexpected error from ListRecommendations") {
			assert.Equal(t, 0, len(result), "No recommendations expected")
			assert.Equal(t, len(zones) + len(regions), mock.numberOfTimesListRecommendationsCalls, "Wrong number of ListRecommendations calls")
			assert.ElementsMatch(t, append(mock.zones, mock.regions...), mock.locationsCalled, "ListRecommendations was called for different locations")
		}
	}
}

type ErrorZonesService struct {
	GoogleService
	err 			error
	regions 		[]string
}

func (s *ErrorZonesService) ListZonesNames(project string) ([]string, error) {
	return []string{}, s.err
}

func (s *ErrorZonesService) ListRegionsNames(project string) ([]string, error) {
	return s.regions, nil
}

func TestErrorInListZones(t *testing.T) {
	errorMessage := "error listing zones"
	regions := []string{"region1", "region2", "region3"}

	_, err := ListRecommendations(&ErrorZonesService{err: fmt.Errorf(errorMessage), regions: regions}, "", "", 2)
	assert.EqualError(t, err, errorMessage, "Expected error calling ListZones")
}

type ErrorRegionsService struct {
	GoogleService
	err 			error
	zones 			[]string
}

func (s *ErrorRegionsService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func (s *ErrorRegionsService) ListRegionsNames(project string) ([]string, error) {
	return []string{}, s.err
}

func TestErrorInListRegions(t *testing.T) {
	errorMessage := "error listing zones"
	zones := []string{"region1", "region2", "region3"}

	_, err := ListRecommendations(&ErrorRegionsService{err: fmt.Errorf(errorMessage), zones: zones}, "", "", 2)
	assert.EqualError(t, err, errorMessage, "Expected error calling ListZones")
}

type ErrorRecommendationService struct {
	GoogleService
	err                	error
	errorLocation       string
	mutex               sync.Mutex
	numberOfTimesCalled int
	zones               []string
	regions             []string
}

func (s *ErrorRecommendationService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func (s *ErrorRecommendationService) ListRegionsNames(project string) ([]string, error) {
	return s.regions, nil
}

func (s *ErrorRecommendationService) ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error) {
	s.mutex.Lock()
	s.numberOfTimesCalled++
	s.mutex.Unlock()

	if location == s.errorLocation {
		return []*gcloudRecommendation{}, s.err
	}
	return []*gcloudRecommendation{}, nil
}

func TestErrorInRecommendations(t *testing.T) {
	errorMessage := "error listing recommendations"
	zones := []string{}
	for i := 1; i <= 5; i++ {
		zones = append(zones, fmt.Sprintf("zone %d", i))
	}

	regions := []string{}
	for i := 1; i <=7; i++ {
		regions = append(regions, fmt.Sprintf("region %d", i))
	}

	locations := append(zones, regions...);

	for _, location := range locations {
		for numConcurrentCalls := 1; numConcurrentCalls <= 10; numConcurrentCalls++ {
			service := &ErrorRecommendationService{
				err:       		fmt.Errorf(errorMessage),
				zones:    		zones,
				regions:   		regions,
				errorLocation: 	location,
			}

			_, err := ListRecommendations(service, "", "", numConcurrentCalls)
			assert.EqualError(t, err, errorMessage, "Expected error calling ListRecommendations")
			assert.Equal(t, len(locations), service.numberOfTimesCalled, "ListRecommendations called wrong number of times")
		}
	}
}

type BenchmarkService struct {
	GoogleService
}

func (s *BenchmarkService) ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error) {
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

func (s *BenchmarkService) ListRegionsNames(project string) ([]string, error) {
	regions := []string{}
	for i := 0; i < 25; i++ {
		regions = append(regions, fmt.Sprintf("zone %d", i))
	}
	return regions, nil
}

func BenchmarkGoroutines(b *testing.B) {
	for _, numConcurrentCalls := range []int{4, 8, 16, 32, 64} {
		b.Run(fmt.Sprintf("%d goroutines:", numConcurrentCalls), func(b *testing.B) {
			s := &BenchmarkService{}
			ListRecommendations(s, "", "", numConcurrentCalls)
		})
	}
}
