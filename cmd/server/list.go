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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/recommender/v1"
)

const defaultNumConcurrentCalls = 100

// ListRecommendationsResponse is response to list/recommendations method
type ListRecommendationsResponse struct {
	Recommendations []*recommender.GoogleCloudRecommenderV1Recommendation `json:"recommendations"`
	FailedProjects  []*automation.ProjectRequirements                     `json:"failedProjects"`
}

// Progress is the structure that contains information about request that is not finished yet.
type Progress struct {
	BatchesProcessed int `json:"batchesProcessed"`
	NumberOfBatches  int `json:"numberOfBatches"`
}

type listRequestHandler struct {
	result             *automation.ListResult
	service            automation.GoogleService
	task               automation.Task
	numConcurrentCalls int
	request            ListRequest
	err                error
}

func (h *listRequestHandler) ListRecommendations() {
	h.task.SetNumberOfSubtasks(1) // 1 call to ListAllProjectsRecommendations
	h.result, h.err = automation.ListProjectsRecommendations(h.service, h.request.Projects, h.numConcurrentCalls, h.task.GetNextSubtask())
	h.task.SetAllDone()
}

func (h *listRequestHandler) GetProgress() (int32, int32) {
	return h.task.GetProgress()
}

func (h *listRequestHandler) GetResult() (*automation.ListResult, error) {
	return h.result, h.err
}

type listInfo struct {
	userEmail string
}

type listRequestsMap struct {
	data  map[listInfo]*listRequestHandler
	mutex sync.Mutex
}

func (m *listRequestsMap) StartListing(info listInfo, handler *listRequestHandler) {
	m.mutex.Lock()
	_, ok := m.data[info]
	if !ok {
		m.data[info] = handler
		go handler.ListRecommendations()
	}
	m.mutex.Unlock()
}

// Return either ListRecommendationsResponse of Progress
func (m *listRequestsMap) GetListRequestResponse(info listInfo) (interface{}, error) {
	m.mutex.Lock()
	handler, ok := m.data[info]
	m.mutex.Unlock()

	if !ok {
		return nil, &googleapi.Error{Message: fmt.Sprintf("Specified request %v is not found", info),
			Code: http.StatusNotFound}
	}

	done, all := handler.GetProgress()
	if done < all {
		return Progress{int(done), int(all)}, nil
	}
	m.DeleteRequest(info)
	listResult, err := handler.GetResult()
	if err != nil {
		return nil, err
	}
	return ListRecommendationsResponse{
		Recommendations: listResult.Recommendations,
		FailedProjects:  listResult.FailedProjects}, nil
}

func (m *listRequestsMap) DeleteRequest(info listInfo) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, info)
}

// ListRequest is a struct containing fields from /recommendations request body.
type ListRequest struct {
	Projects []string `json:"projects"`
}

func getStartListingHandler(service *sharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var listRequest ListRequest

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			sendError(c, fmt.Errorf("Error reading body: %s", err.Error()), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &listRequest)

		if err != nil {
			sendError(c, fmt.Errorf("Error parsing body: %s", err.Error()), http.StatusBadRequest)
			return
		}

		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		service.listRequestsInProcess.StartListing(
			listInfo{user.email},
			&listRequestHandler{service: user.service, request: listRequest, numConcurrentCalls: defaultNumConcurrentCalls})
		c.String(http.StatusOK, "")
	}
}

func getListHandler(service *sharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		response, err := service.listRequestsInProcess.GetListRequestResponse(
			listInfo{user.email})

		if err != nil {
			sendError(c, err)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}
