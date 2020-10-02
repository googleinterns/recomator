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

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	done    bool
	err     error
	mutex   sync.Mutex
	content interface{}
}

func (h *mockHandler) Start() {
}

const inProgressResponse = "request in progresss"

func (h *mockHandler) GetResponse() (Response, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.done {
		return Response{Content: h.content, Error: h.err}, true
	}
	return Response{Content: inProgressResponse}, false
}

func (h *mockHandler) SetDone() {
	h.mutex.Lock()
	h.done = true
	h.mutex.Unlock()
}

func TestAddRequest(t *testing.T) {
	content := "done"
	mock := &mockHandler{content: content}
	requestsMap := NewRequestsMap()
	info := RequestInfo{"me@google.com", "1"}
	err := requestsMap.StartProcessing(info, mock)
	assert.NoError(t, err, "No error expected")

	resp, found := requestsMap.GetResponse(info)
	assert.True(t, found, "Should be in map")
	str, inProgress := resp.Content.(string)
	assert.True(t, inProgress, "Content should be string, because request is in progress")
	assert.Equal(t, str, inProgressResponse, "Should return response in progress")
	mock.SetDone()
	resp, found = requestsMap.GetResponse(info)
	assert.True(t, found, "Should be in map")
	res, ok := resp.Content.(string)
	assert.True(t, ok, "Should be string")
	assert.Equal(t, content, res, "Should be equal to set content")
}

func TestCollision(t *testing.T) {
	content := "done"
	mock := &mockHandler{content: content}
	requestsMap := NewRequestsMap()
	info := RequestInfo{"email", "12"}
	err := requestsMap.StartProcessing(info, mock)
	assert.NoError(t, err, "No error expected")

	err = requestsMap.StartProcessing(info, mock)
	assert.Error(t, err, "Expected collision")

	mock.SetDone()
	resp, found := requestsMap.GetResponse(info)
	assert.True(t, found, "Should be in map")
	res, ok := resp.Content.(string)
	assert.True(t, ok, "Should be string")
	assert.Equal(t, content, res, "Should be equal to set content")
}

func TestError(t *testing.T) {
	errInResponse := fmt.Errorf("Something happened")
	mock := &mockHandler{done: true, err: errInResponse}
	requestsMap := NewRequestsMap()
	info := RequestInfo{"email", "1"}
	err := requestsMap.StartProcessing(info, mock)
	assert.NoError(t, err, "No error expected")
	resp, found := requestsMap.GetResponse(info)
	assert.True(t, found, "Should be in map")
	assert.EqualError(t, resp.Error, errInResponse.Error(), "Error should be returned")
}

func TestMultipleRequests(t *testing.T) {
	numRequests := 5
	requestsMap := NewRequestsMap()
	content := "done"
	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		info := RequestInfo{"email", fmt.Sprint(i)}
		mock := &mockHandler{done: true, content: content}
		go func(info RequestInfo, mock RequestHandler) {
			err := requestsMap.StartProcessing(info, mock)
			assert.NoError(t, err, "No error expected")
			resp, found := requestsMap.GetResponse(info)
			assert.True(t, found, "Should be in map")
			res, ok := resp.Content.(string)
			assert.True(t, ok, "Should be string")
			assert.Equal(t, content, res, "Should be equal to set content")
			wg.Done()
		}(info, mock)
	}
	wg.Wait()
}
