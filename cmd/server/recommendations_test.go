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
	"sync"
	"testing"

	"github.com/googleinterns/recomator/pkg/automation"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	automation.GoogleService
	mutex           sync.Mutex
	wait            bool // ListRegionsNames won't return until it's set to false
	projectsRegions []string
	projectsZones   []string
}

func (s *mockService) StopWaiting() {
	s.mutex.Lock()
	s.wait = false
	s.mutex.Unlock()
}

func (s *mockService) ListRegionsNames(project string) ([]string, error) {
	s.mutex.Lock()
	s.projectsRegions = append(s.projectsRegions, project)
	s.mutex.Unlock()
	for {
		// check if wait is false
		s.mutex.Lock()
		if !s.wait {
			s.mutex.Unlock()
			break
		}
		s.mutex.Unlock()
	}
	return []string{}, nil
}

func (s *mockService) ListZonesNames(project string) ([]string, error) {
	s.mutex.Lock()
	s.projectsZones = append(s.projectsZones, project)
	s.mutex.Unlock()
	return []string{}, nil
}

func (s *mockService) ListAPIRequirements(project string, apis []string) ([]*automation.Requirement, error) {
	return []*automation.Requirement{}, nil
}

func (s *mockService) ListPermissionRequirements(project string, permissions [][]string) ([]*automation.Requirement, error) {
	return []*automation.Requirement{}, nil
}

func TestListingRecommendations(t *testing.T) {
	projects = []string{"some", "projects"}
	mock := &mockService{wait: true} // we set this parameter for ListRegionsNames to wait until we check that progress is returned
	handler := NewListRequestHandler(mock, projects)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		handler.Start()
		wg.Done()
	}()

	resp, done := handler.GetResponse()

	assert.False(t, done, "Should not be done yet")

	progress, ok := resp.Content.(Progress)
	assert.NoError(t, resp.Error, "No error expected")
	assert.True(t, ok, "Progress should be returned")
	assert.True(t, progress.BatchesProcessed < progress.NumberOfBatches, "Should not be done yet")
	mock.StopWaiting() // now ListRegionsNames will return
	wg.Wait()          // wait until mock will return result

	resp, done = handler.GetResponse()
	assert.True(t, done, "Should be done now")
	assert.NoError(t, resp.Error, "No error expected")
	response, ok := resp.Content.(ListRecommendationsResponse)
	assert.True(t, ok, "Should be of type ListRecommendationsResponse")
	assert.Equal(t, 0, len(response.FailedProjects), "No failed projects")
	assert.Equal(t, 0, len(response.Recommendations), "No recommendations")
	assert.ElementsMatch(t, projects, mock.projectsRegions, "ListRegionsNames should be called for all projects")
	assert.ElementsMatch(t, projects, mock.projectsZones, "ListZonesNames should be called for all projects")
}
