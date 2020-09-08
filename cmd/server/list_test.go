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
	"sync/atomic"
	"testing"

	"github.com/googleinterns/recomator/pkg/automation"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	automation.GoogleService
	listProjectsCalls int32 // these int32 fields are used with sync/atomic from different threads
	wait              int32 // ListProjects won't return until it's set to 0
}

func (s *mockService) ListProjects() ([]string, error) {
	atomic.AddInt32(&s.listProjectsCalls, 1)
	for {
		// check if wait is false
		if wait := atomic.LoadInt32(&s.wait); wait == 0 {
			break
		}
	}
	return []string{}, nil
}

func TestListingRecommendations(t *testing.T) {
	listRequests := listRequestsMap{data: make(map[string]*listRequestHandler)}
	mock := &mockService{wait: 1} // we set this parameter for ListProjects to wait until we check that progress is returned
	handler := listRequestHandler{service: mock, numConcurrentCalls: defaultNumConcurrentCalls}

	resp, err := listRequests.GetListRequestResponse("email", &handler)
	assert.NoError(t, err, "No error expected")

	progress, ok := resp.(ListRecommendationsProgressResponse)
	if assert.True(t, ok, "Progress should be returned") {
		atomic.StoreInt32(&mock.wait, 0) // now ListProjects won't wait
		assert.True(t, progress.BatchesProcessed < progress.NumberOfBatches, "Should not be done yet")
	}

	for {
		resp, err := listRequests.GetListRequestResponse("email", nil)
		assert.NoError(t, err, "No error expected")
		response, ok := resp.(ListRecommendationsResponse)
		if ok {
			assert.Equal(t, 0, len(response.FailedProjects), "No failed projects")
			assert.Equal(t, 0, len(response.Recommendations), "No recommendations")
			assert.Equal(t, int32(1), atomic.LoadInt32(&mock.listProjectsCalls), "ListProjects should be called once")
			break
		}
	}
}

func TestMultipleListingRecommendations(t *testing.T) {
	var services []*mockService
	listRequests := listRequestsMap{data: make(map[string]*listRequestHandler)}
	var emails []string

	numUsers := 10
	requestsStarted := make(chan string, numUsers)

	for i := 0; i < numUsers; i++ {
		mock := &mockService{wait: 1}
		services = append(services, mock)
		email := fmt.Sprintf("email %d", i)

		go func(s *mockService, email string) {
			_, err := listRequests.GetListRequestResponse(email, &listRequestHandler{service: s, numConcurrentCalls: defaultNumConcurrentCalls})
			requestsStarted <- email
			assert.NoError(t, err, "Unexpected error while listing")
			atomic.StoreInt32(&s.wait, 0) // now ListProjects will return
		}(mock, email)
		emails = append(emails, email)
	}

	var actualEmails []string
	for range emails {
		email := <-requestsStarted
		actualEmails = append(actualEmails, email)
		for {
			resp, err := listRequests.GetListRequestResponse(email, nil)

			if !assert.NoError(t, err, "Unexpected error while listing") {
				return
			}

			if _, ok := resp.(ListRecommendationsProgressResponse); ok {
				continue
			}
			_, ok := resp.(ListRecommendationsResponse)
			assert.True(t, ok, "Should be of type ListRecommendationsResponse")
			break
		}
	}
	assert.ElementsMatch(t, emails, actualEmails, "Not all requests ")

	for _, mock := range services {
		assert.Equal(t, int32(1), atomic.LoadInt32(&mock.listProjectsCalls), "Should be called exactly once")
	}
}
