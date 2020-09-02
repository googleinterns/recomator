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
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func setUpRouter(auth AuthorizationService) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

	router.GET("/redirect", redirectHandler)

	router.GET("/auth", getAuthHandler(auth))

	router.GET("/recommendations", getListHandler(auth))

	router.POST("/recommendations/apply", getApplyHandler(auth))

	router.GET("/recommendations/checkStatus", getCheckStatusHandler(auth))
	return router
}

var listRequestsInProcess listRequestsMap
var applyRequestsInProcess applyRequestsMap

func main() {
	if err := setConfig(&config, "config.json"); err != nil {
		log.Fatal(err)
	}

	listRequestsInProcess = listRequestsMap{data: make(map[string]*listRequestHandler)} // the key is email address of the user

	applyRequestsInProcess = applyRequestsMap{data: make(map[applyInfo]*applyRequestHandler)} // the key is recommendation name & user email

	auth, err := NewAuthorizationService()
	if err != nil {
		log.Fatal(err)
	}
	router := setUpRouter(auth)

	router.Run(":8000")

}
