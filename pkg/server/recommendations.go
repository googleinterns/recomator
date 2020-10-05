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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"google.golang.org/api/recommender/v1"
)

const defaultNumConcurrentCalls = 200

// ListRecommendationsResponse is response to list/recommendations method
type ListRecommendationsResponse struct {
	Recommendations []*recommender.GoogleCloudRecommenderV1Recommendation `json:"recommendations"`
	FailedProjects  []*automation.ProjectRequirements                     `json:"failedProjects"`
}

type listRequestHandler struct {
	result             *automation.ListResult
	service            automation.GoogleService
	task               automation.Task
	projects           []string
	numConcurrentCalls int
	err                error
}

// NewListRequestHandler creates new listRequestHandler
func NewListRequestHandler(service automation.GoogleService, projects []string) RequestHandler {
	return &listRequestHandler{service: service, projects: projects, numConcurrentCalls: defaultNumConcurrentCalls}
}

func (h *listRequestHandler) Start() {
	h.task.SetNumberOfSubtasks(1) // 1 call to ListProjectsRecommendations
	h.result, h.err = automation.ListProjectsRecommendations(h.service, h.projects, h.numConcurrentCalls, h.task.GetNextSubtask())
	h.task.SetAllDone()
}

func (h *listRequestHandler) GetResponse() (Response, bool) {
	done, all := h.task.GetProgress()
	if done < all {
		return Response{Content: Progress{int(done), int(all)}}, false
	}
	if h.err != nil {
		return Response{Error: h.err}, true
	}
	return Response{Content: ListRecommendationsResponse{
		Recommendations: h.result.Recommendations,
		FailedProjects:  h.result.FailedProjects}}, true
}

// ListRequest contains the body of POST /recommendations request
type ListRequest struct {
	Projects []string `json:"projects"`
}

func getStartListingHandler(service *SharedService) func(c *gin.Context) {
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

		handler := NewListRequestHandler(user.service, listRequest.Projects)
		requestID := StartProcessingWithNewRequestID(&service.requests, user.email, handler)
		c.String(http.StatusCreated, requestID)
	}
}

func getListHandler(service *SharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Query("request_id")
		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		response, ok := service.requests.GetResponse(RequestInfo{user.email, id})

		if !ok {
			sendError(c, fmt.Errorf("No request for %s with id %s", user.email, id), http.StatusNotFound)
			return
		}

		if response.Error != nil {
			sendError(c, response.Error)
			return
		}

		c.JSON(http.StatusOK, response.Content)
	}
}
