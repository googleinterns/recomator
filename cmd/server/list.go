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

	"github.com/googleinterns/recomator/pkg/automation"
	"golang.org/x/oauth2"
	"google.golang.org/api/recommender/v1"
)

const numConcurrentCalls = 100

// ListRecommendationsResponse is response to list/recommendations method
type ListRecommendationsResponse struct {
	Recommendations []*recommender.GoogleCloudRecommenderV1Recommendation `json:"recommendations"`
	FailedProjects  []*automation.ProjectRequirements                     `json:"failedProjects"`
}

// ListRecommendationsProgressResponse is response to list/recommendations method
// if not all recommendations have been processed yet.
type ListRecommendationsProgressResponse struct {
	BatchesProcessed int `json:"batchesProcessed"`
	NumberOfBatches  int `json:"numberOfBatches"`
}

type listRequestHandler struct {
	result  *automation.ListResult
	service automation.GoogleService
	token   *oauth2.Token
	task    automation.Task
	err     error
}

func (h *listRequestHandler) ListRecommendations() {
	h.task.SetNumberOfSubtasks(1) // 1 call to ListAllProjectsRecommendations
	h.result, h.err = automation.ListAllProjectsRecommendations(h.service, numConcurrentCalls, h.task.GetNextSubtask())
	h.task.SetAllDone()
}

func (h *listRequestHandler) GetProgress() (int32, int32) {
	return h.task.GetProgress()
}

func (h *listRequestHandler) GetResult() (*automation.ListResult, error) {
	return h.result, h.err
}

type listRequestsMap struct {
	data  map[oauth2.Token]*listRequestHandler
	mutex sync.Mutex
}

func (m *listRequestsMap) LoadOrStore(token oauth2.Token, handler *listRequestHandler) (*listRequestHandler, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	h, ok := m.data[token]
	if !ok {
		m.data[token] = handler
		h = handler
	}
	return h, ok
}

func (m *listRequestsMap) Delete(token oauth2.Token) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, token)
}
