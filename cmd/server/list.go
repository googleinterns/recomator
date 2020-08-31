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
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"google.golang.org/api/recommender/v1"
)

const defaultNumConcurrentCalls = 100

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
	result             *automation.ListResult
	service            automation.GoogleService
	task               automation.Task
	numConcurrentCalls int
	err                error
}

func (h *listRequestHandler) ListRecommendations() {
	h.task.SetNumberOfSubtasks(1) // 1 call to ListAllProjectsRecommendations
	h.result, h.err = automation.ListAllProjectsRecommendations(h.service, h.numConcurrentCalls, h.task.GetNextSubtask())
	h.task.SetAllDone()
}

func (h *listRequestHandler) GetProgress() (int32, int32) {
	return h.task.GetProgress()
}

func (h *listRequestHandler) GetResult() (*automation.ListResult, error) {
	return h.result, h.err
}

type listRequestsMap struct {
	data  map[string]*listRequestHandler
	mutex sync.Mutex
}

func (m *listRequestsMap) LoadOrStore(email string, handler *listRequestHandler) (*listRequestHandler, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	h, ok := m.data[email]
	if !ok {
		m.data[email] = handler
		h = handler
	}
	return h, ok
}

func (m *listRequestsMap) Delete(email string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, email)
}

func getListHandler(authService AuthorizationService) func(c *gin.Context) {
	return func(c *gin.Context) {
		authCode := c.Request.Header["Authorization"]
		if len(authCode) != 0 {
			user, err := authService.Authorize(authCode[0])

			if err == nil {

				handler, loaded := listRequestsInProcess.LoadOrStore(user.email, &listRequestHandler{service: user.service, numConcurrentCalls: defaultNumConcurrentCalls})
				if !loaded {
					go handler.ListRecommendations()
				}
				done, all := handler.GetProgress()
				if done < all {
					c.JSON(http.StatusOK, ListRecommendationsProgressResponse{int(done), int(all)})
				} else {
					listRequestsInProcess.Delete(user.email)
					listResult, err := handler.GetResult()
					if err != nil {
						sendError(c, err)
					} else {
						c.JSON(http.StatusOK, ListRecommendationsResponse{
							Recommendations: listResult.Recommendations,
							FailedProjects:  listResult.FailedProjects})
					}

				}
				return
			}
			sendError(c, err, http.StatusUnauthorized)
			return
		}
		sendError(c, fmt.Errorf("Authorization code not specified"), http.StatusUnauthorized)
	}
}
