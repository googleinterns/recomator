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
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"google.golang.org/api/googleapi"
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

func (h *applyRequestHandler) Apply() {
	h.task.SetNumberOfSubtasks(1) // 1 call to ApplyByName
	h.err = automation.ApplyByName(h.service, h.name, h.task.GetNextSubtask())
	h.task.SetAllDone()
}

func (h *applyRequestHandler) GetProgress() (int32, int32) {
	return h.task.GetProgress()
}

func (h *applyRequestHandler) GetStatus() string {
	if h.err != nil {
		return failedStatus
	}
	return succeededStatus
}

func (h *applyRequestHandler) GetError() error {
	return h.err
}

type applyInfo struct {
	recommendationName string
	userEmail          string
}

type applyRequestsMap struct {
	data  map[applyInfo]*applyRequestHandler
	mutex sync.Mutex
}

func (m *applyRequestsMap) DeleteRequest(info applyInfo) {
	m.mutex.Lock()
	delete(m.data, info)
	m.mutex.Unlock()
}

// Returns checkStatusResponse if request is in process.
// If there's no such request returns false in second value.
func (m *applyRequestsMap) CheckStatus(info applyInfo) (CheckStatusResponse, bool) {
	m.mutex.Lock()
	handler, ok := m.data[info]
	m.mutex.Unlock()
	if ok {
		done, all := handler.GetProgress()
		status := inProgressStatus
		errMessage := ""
		if done == all {
			m.DeleteRequest(info)
			status = handler.GetStatus()
			if status == failedStatus {
				errMessage = handler.GetError().Error()
			}
		}
		return CheckStatusResponse{status, errMessage}, true
	}
	return CheckStatusResponse{}, false
}

// returns error if failed to start applying (e.g. recommendation is already being applied)
func (m *applyRequestsMap) StartApplying(info applyInfo, handler *applyRequestHandler) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.data[info]
	if !ok {
		m.data[info] = handler
		go handler.Apply()
		return nil
	}
	return &googleapi.Error{Message: "Recommendation is already being applied", Code: http.StatusMethodNotAllowed}

}

func getApplyHandler(service *sharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		err = service.applyRequestsInProcess.StartApplying(applyInfo{name, user.email}, &applyRequestHandler{service: user.service, name: name})
		if err != nil {
			sendError(c, err)
		}
	}
}

func getCheckStatusHandler(service *sharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		user, err := authorizeRequest(service.auth, c.Request)

		if err != nil {
			sendError(c, err)
			return
		}

		response, loaded := service.applyRequestsInProcess.CheckStatus(applyInfo{name, user.email})

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
