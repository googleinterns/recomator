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

func (m *applyRequestsMap) Load(info applyInfo) (*applyRequestHandler, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	h, ok := m.data[info]
	return h, ok
}

func (m *applyRequestsMap) LoadOrStore(info applyInfo, handler *applyRequestHandler) (*applyRequestHandler, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	h, ok := m.data[info]
	if !ok {
		m.data[info] = handler
		h = handler
	}
	return h, ok
}

func getApplyHandler(authService AuthorizationService) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		user, err := authorizeRequest(authService, c.Request)

		if err != nil {
			sendError(c, err, http.StatusUnauthorized)
			return
		}

		handler, loaded := applyRequestsInProcess.LoadOrStore(applyInfo{name, user.email}, &applyRequestHandler{service: user.service})
		if !loaded {
			go handler.Apply()
		} else {
			sendError(c, fmt.Errorf("Recommendation is already being applied"), http.StatusMethodNotAllowed)
		}
	}
}

func getCheckStatusHandler(authService AuthorizationService) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Query("name")
		user, err := authorizeRequest(authService, c.Request)

		if err != nil {
			sendError(c, err, http.StatusUnauthorized)
			return
		}

		handler, loaded := applyRequestsInProcess.Load(applyInfo{name, user.email})
		if loaded {
			done, all := handler.GetProgress()
			status := inProgressStatus
			errMessage := ""
			if done == all {
				status = handler.GetStatus()
				if status == failedStatus {
					errMessage = handler.GetError().Error()
				}
			}
			c.JSON(http.StatusOK, CheckStatusResponse{status, errMessage})
		} else {
			rec, err := user.service.GetRecommendation(name)
			if err != nil {
				sendError(c, err)
			} else {
				c.JSON(http.StatusOK, CheckStatusResponse{Status: rec.StateInfo.State})
			}
		}

	}
}
