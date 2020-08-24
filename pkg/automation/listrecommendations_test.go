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
	regions                               []string
	callsToList                           []query
}

type query struct {
	location      string
	recommenderID string
}

func (s *MockService) ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error) {
	s.mutex.Lock()
	s.numberOfTimesListRecommendationsCalls++
	s.callsToList = append(s.callsToList, query{location, recommenderID})
	s.mutex.Unlock()
	return []*gcloudRecommendation{nil}, nil
}

func (s *MockService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func (s *MockService) ListRegionsNames(project string) ([]string, error) {
	return s.regions, nil
}

func makeQueries(locations []string) []query {
	var queries []query
	for _, rec := range googleRecommenders {
		for _, loc := range locations {
			queries = append(queries, query{loc, rec})
		}
	}
	return queries
}

func TestListRecommendations(t *testing.T) {
	for numConcurrentCalls := 0; numConcurrentCalls <= 7; numConcurrentCalls++ {
		zones := []string{"zone1", "zone2", "zone3"}
		regions := []string{"region1", "region2", "region3"}
		mock := &MockService{zones: zones, regions: regions}
		task := &Task{}
		result, err := ListRecommendations(mock, "", numConcurrentCalls, task)

		if assert.NoError(t, err, "Unexpected error from ListRecommendations") {
			locations := append(mock.zones, mock.regions...)
			queries := makeQueries(locations)
			assert.Equal(t, len(queries), len(result), "One recommendation from each query was expected")
			assert.Equal(t, len(queries), mock.numberOfTimesListRecommendationsCalls, "Wrong number of ListRecommendations calls")
			assert.ElementsMatch(t, queries, mock.callsToList, "ListRecommendations was called for different locations and recommenders")

			done, all := task.GetProgress()
			assert.True(t, done == all, "List recommendations task should be done already")
			task.mutex.Lock()
			assert.Equal(t, len(task.subtasks), task.subtasksDone, "All subtasks should be done")
			task.mutex.Unlock()
		}
	}
}

type ErrorZonesService struct {
	GoogleService
	err     error
	regions []string
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

	task := &Task{}
	_, err := ListRecommendations(&ErrorZonesService{err: fmt.Errorf(errorMessage), regions: regions}, "", 2, task)
	assert.EqualError(t, err, errorMessage, "Expected error calling ListZones")

	done, all := task.GetProgress()
	assert.True(t, done < all, "List recommendations task should be not finished because of error")
}

type ErrorRegionsService struct {
	GoogleService
	err   error
	zones []string
}

func (s *ErrorRegionsService) ListZonesNames(project string) ([]string, error) {
	return s.zones, nil
}

func (s *ErrorRegionsService) ListRegionsNames(project string) ([]string, error) {
	return []string{}, s.err
}

func TestErrorInListRegions(t *testing.T) {
	errorMessage := "error listing regions"
	zones := []string{"zone1", "zone2", "zone3"}

	task := &Task{}
	_, err := ListRecommendations(&ErrorRegionsService{err: fmt.Errorf(errorMessage), zones: zones}, "", 2, task)
	assert.EqualError(t, err, errorMessage, "Expected error calling ListRegions")

	done, all := task.GetProgress()
	assert.True(t, done < all, "List recommendations task should be not finished because of error")
}

type ErrorRecommendationService struct {
	GoogleService
	err                 error
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
	for i := 1; i <= 3; i++ {
		regions = append(regions, fmt.Sprintf("region %d", i))
	}

	locations := append(zones, regions...)

	for _, location := range locations {
		for numConcurrentCalls := 1; numConcurrentCalls <= 10; numConcurrentCalls++ {
			service := &ErrorRecommendationService{
				err:           fmt.Errorf(errorMessage),
				zones:         zones,
				regions:       regions,
				errorLocation: location,
			}

			task := &Task{}
			_, err := ListRecommendations(service, "", numConcurrentCalls, task)
			assert.EqualError(t, err, errorMessage, "Expected error calling ListRecommendations")
			numQueries := len(locations) * len(googleRecommenders)
			assert.Equal(t, numQueries, service.numberOfTimesCalled, "ListRecommendations called wrong number of times")

			done, all := task.GetProgress()
			assert.True(t, done < all, "List recommendations task should be not finished because of error")
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
		regions = append(regions, fmt.Sprintf("region %d", i))
	}
	return regions, nil
}

func BenchmarkGoroutines(b *testing.B) {
	for _, numConcurrentCalls := range []int{4, 8, 16, 32, 64, 128} {
		b.Run(fmt.Sprintf("%d goroutines:", numConcurrentCalls), func(b *testing.B) {
			s := &BenchmarkService{}
			ListRecommendations(s, "", numConcurrentCalls, &Task{})
		})
	}
}

type projectRecommender struct {
	project       string
	recommenderID string
}

type MockProjectsService struct {
	GoogleService
	queries                          []projectRecommender
	numberOfListRecommendationsCalls int
	apiCalls                         []string
	permissionCalls                  []string
	mutex                            sync.Mutex
	projects                         []string
}

func (s *MockProjectsService) ListProjects() ([]string, error) {
	return s.projects, nil
}

func (s *MockProjectsService) ListZonesNames(project string) ([]string, error) {
	return []string{"one zone"}, nil
}

func (s *MockProjectsService) ListRegionsNames(project string) ([]string, error) {
	return nil, nil
}

func (s *MockProjectsService) ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error) {
	s.mutex.Lock()
	s.numberOfListRecommendationsCalls++
	s.queries = append(s.queries, projectRecommender{project, recommenderID})
	s.mutex.Unlock()
	return []*gcloudRecommendation{nil}, nil
}

func makeProjectsQueries(projects []string) []projectRecommender {
	var result []projectRecommender
	for _, pr := range projects {
		for _, rec := range googleRecommenders {
			result = append(result, projectRecommender{pr, rec})
		}
	}
	return result
}

var okRequirements = []*Requirement{&Requirement{Status: RequirementCompleted}}

func (s *MockProjectsService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.apiCalls = append(s.apiCalls, project)
	if project == failedProject {
		return failedRequirements, nil
	}
	return okRequirements, nil
}

func (s *MockProjectsService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.permissionCalls = append(s.permissionCalls, project)
	if project == failedProject {
		return failedRequirements, nil
	}
	return okRequirements, nil
}

func TestListAllProjectsRecommendations(t *testing.T) {
	for numConcurrentCalls := 0; numConcurrentCalls < 10; numConcurrentCalls++ {
		for numProjects := 0; numProjects < 5; numProjects++ {
			for numFailed := 0; numFailed <= numProjects; numFailed++ {
				var okProjects, failedProjects []string
				for i := 0; i < numFailed; i++ {
					failedProjects = append(failedProjects, failedProject)
				}
				numOk := numProjects - numFailed
				for i := 0; i < numOk; i++ {
					okProjects = append(okProjects, fmt.Sprintf("project %d", i))
				}
				projects := append(okProjects, failedProjects...)
				task := &Task{}
				mock := &MockProjectsService{projects: projects}
				res, err := ListAllProjectsRecommendations(mock, numConcurrentCalls, task)
				if assert.NoError(t, err) {
					done, all := task.GetProgress()
					assert.True(t, done == all, "Task List all recommendations should be finished already")
					task.mutex.Lock()
					assert.Equal(t, task.subtasksDone, len(task.subtasks), "All subtasks should be done")
					task.mutex.Unlock()

					queries := makeProjectsQueries(okProjects)
					assert.Equal(t, len(queries), mock.numberOfListRecommendationsCalls, "List recommendations called wrong number of times")
					assert.ElementsMatch(t, queries, mock.queries, "List Recommendations was called with wrong parameters")

					assert.ElementsMatch(t, projects, mock.apiCalls, "List api requirements was called for different projects")
					assert.ElementsMatch(t, okProjects, mock.permissionCalls, "List permission requirements was called for different projects")

					assert.Equal(t, len(queries), len(res.recommendations), "Wrong number of overall recommendations")
					var failedProjectsRequirements []*ProjectRequirements
					for i := 0; i < numFailed; i++ {
						failedProjectsRequirements = append(failedProjectsRequirements, &ProjectRequirements{Project: failedProject, Requirements: failedRequirements})
					}
					assert.ElementsMatch(t, failedProjectsRequirements, res.failedProjects, "Wrong failed projects requirements list")
				}
			}
		}
	}
}
