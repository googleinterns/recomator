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
	"fmt"
	"sync"
	"testing"

	"github.com/googleinterns/recomator/pkg/automation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/recommender/v1"
)

type mockApply struct {
	automation.GoogleService
	names []string
	mutex sync.Mutex
	wait  bool // doesn't finish applying until wait is true
}

func (s *mockApply) GetRecommendation(name string) (*recommender.GoogleCloudRecommenderV1Recommendation, error) {
	s.mutex.Lock()
	s.names = append(s.names, name)
	s.mutex.Unlock()
	return nil, fmt.Errorf("something happened")
}

func (s *mockApply) StopWaiting() {
	s.mutex.Lock()
	s.wait = false
	s.mutex.Unlock()
}

func CheckApplyingFails(t *testing.T, info applyInfo, applyRequests *applyRequestsMap) {
	for {
		status, ok := applyRequests.CheckStatus(info)
		assert.True(t, ok, "Recommendation should be in progress/applied now")
		if status.Status == failedStatus {
			assert.NotEmpty(t, status.ErrorMessage, "There should be some error message")
			break
		}
		assert.Equal(t, inProgressStatus, status.Status, "Recommendation should be in progress")
	}
}
func TestApplySimple(t *testing.T) {
	applyRequests := applyRequestsMap{data: make(map[applyInfo]*applyRequestHandler)}
	mock := &mockApply{wait: true}
	handler := &applyRequestHandler{name: "", service: mock}
	info := applyInfo{}
	err := applyRequests.StartApplying(info, handler)
	assert.NoError(t, err, "No error is expected for StartApplying")

	status, ok := applyRequests.CheckStatus(info)
	assert.True(t, ok, "Recommendation should be applied now")
	assert.Equal(t, inProgressStatus, status.Status, "Recommendation should be in progress")
	mock.StopWaiting()

	CheckApplyingFails(t, info, &applyRequests)

	assert.Equal(t, mock.names, []string{""})
}

func TestApplyTwoTimes(t *testing.T) {
	applyRequests := applyRequestsMap{data: make(map[applyInfo]*applyRequestHandler)}
	mock := &mockApply{wait: true}
	handler := &applyRequestHandler{name: "", service: mock}
	info := applyInfo{}
	err := applyRequests.StartApplying(info, handler)
	assert.NoError(t, err, "No error is expected for StartApplying")

	err = applyRequests.StartApplying(info, handler)
	assert.Error(t, err, "There should be an error because recommendation if already being applied")
	mock.StopWaiting()

	// check that request will still be processed
	CheckApplyingFails(t, info, &applyRequests)
}

func TestCheckStatusNotExisting(t *testing.T) {
	applyRequests := applyRequestsMap{data: make(map[applyInfo]*applyRequestHandler)}
	info := applyInfo{}
	_, ok := applyRequests.CheckStatus(info)
	assert.False(t, ok, "There shouldn't be such recommendation")
}

func TestCheckRequestIsDeleted(t *testing.T) {
	applyRequests := applyRequestsMap{data: make(map[applyInfo]*applyRequestHandler)}
	mock := &mockApply{}
	handler := &applyRequestHandler{name: "", service: mock}
	info := applyInfo{}
	applyRequests.StartApplying(info, handler)
	CheckApplyingFails(t, info, &applyRequests)
	_, ok := applyRequests.CheckStatus(info)
	assert.False(t, ok, "This recommendation should be already deleted")
}

func TestMultipleUsersAndRecommendations(t *testing.T) {
	applyRequests := applyRequestsMap{data: make(map[applyInfo]*applyRequestHandler)}
	var mocks []*mockApply
	numUsers := 5
	for i := 0; i < numUsers; i++ {
		mocks = append(mocks, &mockApply{})
	}
	numRequests := 5
	var recommendations []string
	for i := 0; i < numRequests; i++ {
		recommendations = append(recommendations, fmt.Sprintf("rec %d", i))
	}

	var wg sync.WaitGroup
	for i, mock := range mocks {
		for _, recommendation := range recommendations {
			info := applyInfo{userEmail: fmt.Sprintf("email %d", i), recommendationName: recommendation}
			wg.Add(1)
			go func(info applyInfo, mock automation.GoogleService) {
				defer wg.Done()
				handler := &applyRequestHandler{name: info.recommendationName, service: mock}
				err := applyRequests.StartApplying(info, handler)
				assert.NoError(t, err, "No error is expected for StartApplying")
				CheckApplyingFails(t, info, &applyRequests)
			}(info, mock)
		}
	}

	wg.Wait()
}
