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
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAllCompletedService struct {
	GoogleService
}

func (s *mockAllCompletedService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	var result []*Requirement
	for _, api := range apis {
		result = append(result, &Requirement{Name: api, Status: RequirementCompleted})
	}
	return result, nil
}

func (s *mockAllCompletedService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	var result []*Requirement
	for _, perm := range permissions {
		result = append(result, &Requirement{Name: perm[0], Status: RequirementCompleted})
	}
	return result, nil
}

func checkAllRequirementsCompleted(t *testing.T, reqs []*Requirement) {
	var actualNames, expectedNames []string
	for _, req := range reqs {
		assert.Equal(t, RequirementCompleted, req.Status, "Requirements should be compeleted")
		actualNames = append(actualNames, req.Name)
	}
	for _, api := range requiredAPIs {
		expectedNames = append(expectedNames, api)
	}
	for _, perm := range requiredPermissions {
		expectedNames = append(expectedNames, perm[0])
	}
	assert.ElementsMatch(t, expectedNames, actualNames, "Requirements should contain all APIs and permissions")
}

func TestAllCompleted(t *testing.T) {
	mock := &mockAllCompletedService{}
	reqs, err := ListProjectRequirements(mock, "")
	if assert.NoError(t, err, "No error from ListProjectRequirements expected") {
		checkAllRequirementsCompleted(t, reqs)
	}
}

type mockService struct {
	GoogleService
	apiReqs               []*Requirement
	permissionReqs        []*Requirement
	listPermissionsCalled bool
}

func (s *mockService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	return s.apiReqs, nil
}

func (s *mockService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	s.listPermissionsCalled = true
	return s.permissionReqs, nil
}

var failedRequirements = []*Requirement{&Requirement{Status: RequirementFailed}}

func TestFailAPIRequirements(t *testing.T) {
	mock := &mockService{apiReqs: failedRequirements}

	reqs, err := ListProjectRequirements(mock, "")
	if assert.NoError(t, err, "ListProjectRequirements should not return error") {
		assert.ElementsMatch(t, mock.apiReqs, reqs, "ListProjectRequirements should return api requirements")
		assert.False(t, mock.listPermissionsCalled, "ListPermissionRequirements should not be called if api reqs are failed")
	}
}

func TestFailPermissionRequirements(t *testing.T) {
	mock := &mockService{apiReqs: []*Requirement{&Requirement{Status: RequirementCompleted}},
		permissionReqs: failedRequirements}

	reqs, err := ListProjectRequirements(mock, "")
	if assert.NoError(t, err, "ListProjectRequirements should not return error") {
		assert.ElementsMatch(t, append(mock.apiReqs, mock.permissionReqs...), reqs,
			"ListProjectRequirements should return api & permission requirements, if api requirements are completed")
	}
}

type errorAPIService struct {
	GoogleService
	err error
}

func (s *errorAPIService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	return nil, s.err
}

type errorPermissionService struct {
	GoogleService
	err error
}

func (s *errorPermissionService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	return nil, nil
}

func (s *errorPermissionService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	return nil, s.err
}

func TestErrorAPIRequirements(t *testing.T) {
	for _, mock := range []GoogleService{&errorAPIService{err: errors.New("hi! i'm error")},
		&errorPermissionService{err: errors.New("another error")}} {
		reqs, err := ListProjectRequirements(mock, "")
		if assert.Error(t, err, "ListProjectRequirements should result in error") {
			assert.Nil(t, reqs, "Only one of returned values should be non-nil")
		}
	}
}

const (
	errorProject  = "err"
	failedProject = "fail"
)

type mockProjectService struct {
	GoogleService
}

func getService(project string) GoogleService {
	var service GoogleService
	switch project {
	case errorProject:
		service = &errorAPIService{err: errors.New("")}
	case failedProject:
		service = &mockService{apiReqs: failedRequirements}
	default:
		service = &mockAllCompletedService{}
	}
	return service
}

func (s *mockProjectService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	rec, err := getService(project).ListAPIRequirements(project, apis)
	return rec, err
}

func (s *mockProjectService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	rec, err := getService(project).ListPermissionRequirements(project, permissions)
	return rec, err
}

func TestListRequirements(t *testing.T) {
	rand.Seed(42)
	for numProjects := 0; numProjects < 5; numProjects++ {
		for i := 0; i < 100; i++ {
			var projects []string
			for ind := 0; ind < numProjects; ind++ {
				project := fmt.Sprintf("project %d", ind)
				if rand.Int()%2 == 0 {
					project = failedProject
				}
				projects = append(projects, project)
			}

			task := &Task{}
			reqs, err := ListRequirements(&mockProjectService{}, projects, task)
			if assert.NoError(t, err, "No error expected from ListRequirements") {
				var actualProjects []string
				for _, req := range reqs {
					actualProjects = append(actualProjects, req.Project)
					if req.Project == failedProject {
						assert.ElementsMatch(t, failedRequirements, req.Requirements, "Should be equal to failed requirements for this project")
					} else {
						checkAllRequirementsCompleted(t, req.Requirements)
					}
				}
				assert.ElementsMatch(t, projects, actualProjects, "Should contain the same projects")

				done, all := task.GetProgress()
				assert.True(t, done == all, "ListRequirements should be done now")
				task.mutex.Lock()
				assert.Equal(t, task.subtasksDone, len(task.subtasks), "All subtasks should be done")
				task.mutex.Unlock()
			}
		}
	}
}

func TestErrorListRequirements(t *testing.T) {
	projects := []string{"ok", errorProject, failedProject}
	task := &Task{}
	reqs, err := ListRequirements(&mockProjectService{}, projects, task)
	if assert.Error(t, err) {
		assert.Nil(t, reqs, "Only one value should be non-nil")
		done, all := task.GetProgress()
		assert.True(t, done < all, "ListRequirements should not be done because of error")
	}
}
