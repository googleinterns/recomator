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

package server

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
	for {
		// check if wait is false
		s.mutex.Lock()
		if !s.wait {
			s.mutex.Unlock()
			break
		}
		s.mutex.Unlock()
	}
	return nil, fmt.Errorf("something happened")
}

func (s *mockApply) StopWaiting() {
	s.mutex.Lock()
	s.wait = false
	s.mutex.Unlock()
}

func TestApplySimple(t *testing.T) {
	mock := &mockApply{wait: true}
	handler := NewApplyRequestHandler(mock, "name")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		handler.Start()
		wg.Done()
	}()

	resp, done := handler.GetResponse()
	assert.False(t, done, "Recommendation should be in progress now")
	status, ok := resp.Content.(CheckStatusResponse)
	assert.True(t, ok, "Response shouldbe of type CheckStatusResponse")
	assert.Equal(t, inProgressStatus, status.Status, "Recommendation should be in progress")
	mock.StopWaiting() // mock will return here
	wg.Wait()          // wait until returns

	resp, done = handler.GetResponse()
	status, ok = resp.Content.(CheckStatusResponse)
	assert.True(t, ok, "Response should be of type CheckStatusResponse")
	assert.NoError(t, resp.Error, "No error expected")
	assert.True(t, done, "Should be done already")
	assert.Equal(t, failedStatus, status.Status, "Status should be failed")
	assert.NotEmpty(t, status.ErrorMessage, "There should be some error message")

	assert.Equal(t, mock.names, []string{"name"})
}
