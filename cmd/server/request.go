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
	"math/rand"
	"net/http"
	"sync"
	"time"

	"google.golang.org/api/googleapi"
)

// Response contains response that should be sent in response to the request.
// If Error is nil, response should be sent with http.StatusOK.
// Otherwise, something unexpected happened and error is specified in the Error field.
type Response struct {
	Content interface{}
	Error   error
}

// RequestHandler is the type that can start request processing and get results
type RequestHandler interface {
	// Start method starts doing the request
	Start()
	// GetResponse returns the current response and whether the request is already done.
	// For example, while request is still in proccess Response will contain some progress info,
	// and after request is done, Response will contain result, second value will be `true`.
	GetResponse() (Response, bool)
}

// RequestInfo contains information needed to indentify the request
type RequestInfo struct {
	email     string
	requestID string
}

// Progress is the structure that stores the fraction of work done for some request
type Progress struct {
	BatchesProcessed int `json:"batchesProcessed"`
	NumberOfBatches  int `json:"numberOfBatches"`
}

// RequestsMap contains current requests and handles getting response for them,
// deleting, adding new requests.
type RequestsMap struct {
	data  map[RequestInfo]RequestHandler
	mutex sync.Mutex
}

// NewRequestsMap creates new RequestsMap
func NewRequestsMap() RequestsMap {
	return RequestsMap{data: make(map[RequestInfo]RequestHandler)}
}

func (m *RequestsMap) deleteRequest(info RequestInfo) {
	m.mutex.Lock()
	delete(m.data, info)
	m.mutex.Unlock()
}

// StartProcessing starts the processing of the request.
// If such request is already in the map, returns error.
func (m *RequestsMap) StartProcessing(info RequestInfo, handler RequestHandler) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.data[info]
	if !ok {
		m.data[info] = handler
		go handler.Start()
		return nil
	}
	return &googleapi.Error{Message: "The request was already added", Code: http.StatusMethodNotAllowed}
}

// GetResponse returns response if request is in process or finished.
// If request is finished, deletes it from the map.
// If there's no such request returns false in second value.
func (m *RequestsMap) GetResponse(info RequestInfo) (Response, bool) {
	m.mutex.Lock()
	handler, ok := m.data[info]
	m.mutex.Unlock()
	if ok {
		response, done := handler.GetResponse()
		if done {
			m.deleteRequest(info)
		}
		return response, true
	}
	return Response{}, false
}

// StartProcessingWithNewRequestID finds unused request ID and calls StartProcessing method of RequestsMap
func StartProcessingWithNewRequestID(requests *RequestsMap, email string, handler RequestHandler) string {
	lengthOfID := 20
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	requestID := randomString(lengthOfID, generator)

	for {
		err := requests.StartProcessing(RequestInfo{email, requestID}, handler)
		// StartProcessing will only fail if we had a collision
		if err != nil {
			// create new requestID to avoid collision (most likely this code will never be reached as collisions are almost impossible)
			requestID = randomString(lengthOfID, generator)
		} else {
			break
		}
	}
	return requestID
}
