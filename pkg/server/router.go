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
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/recommender/v1"
)

type gcloudRecommendation recommender.GoogleCloudRecommenderV1Recommendation

// ErrorResponse is response with error containing error message.
type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

const defaultErrorCode = http.StatusInternalServerError

// sendError sends error message in ErrorResponse.
// If the code is not specified in err, errorCode will be used.
// If errorCode is not specified, defaultErrorCode will be used.
func sendError(c *gin.Context, err error, errorCode ...int) {
	log.Printf("Error happened while processing request: %s: %v", c.Request.URL, err)
	googleErr, ok := err.(*googleapi.Error)
	if ok {
		c.JSON(googleErr.Code, ErrorResponse{googleErr.Message})
		return
	}
	code := defaultErrorCode
	if len(errorCode) != 0 {
		code = errorCode[0]
	}
	c.JSON(code, ErrorResponse{err.Error()})
}

// SetUpRouter creates new router using SharedService.
func SetUpRouter(service *SharedService) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

	router.GET("/api/redirect", redirectHandler(service))

	router.GET("/api/auth", getAuthHandler(service))

	router.GET("/api/projects", getProjectsHandler(service))

	router.POST("/api/requirements", getStartCheckingHandler(service))

	router.GET("/api/requirements", getCheckRequirementsHandler(service))

	router.POST("/api/recommendations", getStartListingHandler(service))

	router.GET("/api/recommendations", getListHandler(service))

	router.POST("/api/recommendations/apply", getApplyHandler(service))

	router.GET("/api/recommendations/checkStatus", getCheckStatusHandler(service))
	return router
}

// SharedService is the struct that contains authorization service
// and information about currently processed requests.
type SharedService struct {
	auth     AuthorizationService
	requests RequestsMap
}

// NewSharedService creates new sharedService to access GoogleAPIs.
func NewSharedService(conf oauth2.Config) (*SharedService, error) {
	var service SharedService
	auth, err := NewAuthorizationService(conf)
	if err != nil {
		return nil, err
	}
	service.auth = auth
	service.requests = NewRequestsMap()
	return &service, nil
}
