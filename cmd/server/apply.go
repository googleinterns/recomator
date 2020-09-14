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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
)

const (
	inProgressStatus = "IN PROGRESS"
	failedStatus     = "FAILED"
	succeededStatus  = "SUCCEEDED"
)

// CheckStatusResponse is the response to recommendations/name/checkStatus method
type CheckStatusResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type applyRequestHandler struct {
	service automation.GoogleService
	name    string
	err     error
	task    automation.Task
}

// NewApplyRequestHandler creates new applyRequestHandler
func NewApplyRequestHandler(service automation.GoogleService, name string) RequestHandler {
	return &applyRequestHandler{service: service, name: name}
}

func (h *applyRequestHandler) Start() {
	h.task.SetNumberOfSubtasks(1) // 1 call to ApplyByName
	h.err = automation.ApplyByName(h.service, h.name, h.task.GetNextSubtask())
	h.task.SetAllDone()
}

func (h *applyRequestHandler) GetResponse() (Response, bool) {
	done, all := h.task.GetProgress()
	var response interface{}
	finished := false
	if done < all {
		response = CheckStatusResponse{Status: inProgressStatus}
	} else {
		finished = true
		if h.err != nil {
			response = CheckStatusResponse{Status: failedStatus, ErrorMessage: h.err.Error()}
		} else {
			response = CheckStatusResponse{Status: succeededStatus}
		}
	}
	return Response{Content: response}, finished
}

func getApplyHandler(service *SharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		err = service.requests.StartProcessing(RequestInfo{user.email, name},
			NewApplyRequestHandler(user.service, name))
		if err != nil {
			sendError(c, err)
			return
		}
		c.String(http.StatusCreated, "")
	}
}

func getCheckStatusHandler(service *SharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		response, loaded := service.requests.GetResponse(RequestInfo{user.email, name})

		if loaded {
			c.JSON(http.StatusOK, response)
			return
		}

		rec, err := user.service.GetRecommendation(name)
		if err != nil {
			sendError(c, err)
		} else {
			c.JSON(http.StatusOK, CheckStatusResponse{Status: rec.StateInfo.State})
		}
	}
}
