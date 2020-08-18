package automation

import (
	"errors"
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

func TestAllCompleted(t *testing.T) {
	mock := &mockAllCompletedService{}
	reqs, err := ListProjectRequirements(mock, "")
	if assert.NoError(t, err, "No error from ListProjectrequierements expected") {
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
		assert.EqualValues(t, expectedNames, actualNames, "Requirements should contain all APIs and permissions")
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
		assert.EqualValues(t, mock.apiReqs, reqs, "ListProjectRequirements should return api requirements")
		assert.False(t, mock.listPermissionsCalled, "ListPermissionRequirements should not be called if api reqs are failed")
	}
}

func TestFailPermissionRequirements(t *testing.T) {
	mock := &mockService{apiReqs: []*Requirement{&Requirement{Status: RequirementCompleted}},
		permissionReqs: failedRequirements}

	reqs, err := ListProjectRequirements(mock, "")
	if assert.NoError(t, err, "ListProjectRequirements should not return error") {
		assert.EqualValues(t, append(mock.apiReqs, mock.permissionReqs...), reqs,
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
