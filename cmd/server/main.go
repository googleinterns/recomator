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
	"os"

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

func main() {
	if err := setConfig(&config, "config.json"); err != nil {
		log.Fatal(err)
	}

	service, err := NewSharedService()
	if err != nil {
		log.Fatal(err)
	}
	router := setUpRouter(service)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)

}
