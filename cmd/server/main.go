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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/recommender/v1"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
)

type gcloudRecommendation recommender.GoogleCloudRecommenderV1Recommendation

// ErrorResponse is response with error containing error message.
type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func setAuthorizationHeader(c *gin.Context, token *oauth2.Token) bool {
	tokenByte, err := json.Marshal(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{ErrorMessage: err.Error()})
		return false
	}
	c.Header("Authorization", string(tokenByte))
	return true
}

func sendError(c *gin.Context, err error) {
	googleErr, ok := err.(*googleapi.Error)
	if ok {
		c.JSON(googleErr.Code, ErrorResponse{googleErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
}

func main() {
	if err := setConfig(&conf, "config.json"); err != nil {
		log.Fatal(err)
	}
	listRequestsInProcess := listRequestsMap{data: make(map[oauth2.Token]*listRequestHandler)} // the key is the token

	// applyRequestsInProcess := applyRequestsMap{data: make(map[string]*applyRequestHandler)} // the key is recommendation name

	router := gin.Default()
	router.Use(corsMiddleware())

	router.GET("/auth", func(c *gin.Context) {
		code := c.Query("code")
		path := c.Query("state")
		token, err := conf.Exchange(context.Background(), code, authOptions...)
		if err != nil {
			sendError(c, err)
			return
		}
		if !setAuthorizationHeader(c, token) {
			return
		}
		c.JSON(http.StatusUnauthorized, ErrorResponse{ErrorMessage: fmt.Sprintf("Resend %s request with this authorization header", path)})
	})

	router.GET("/recommendations", func(c *gin.Context) {
		tokenStrings := c.Request.Header["Authorization"]
		log.Printf("Authorization header: %v\n", tokenStrings)
		if len(tokenStrings) != 0 {
			var token oauth2.Token
			log.Println(tokenStrings[0])
			tokenByte := []byte(tokenStrings[0])

			err := json.Unmarshal(tokenByte, &token)
			if err == nil {
				newService, err := automation.NewGoogleService(context.Background(), &conf, &token)
				if err != nil {
					sendError(c, err)
					return
				}

				handler, loaded := listRequestsInProcess.LoadOrStore(token, &listRequestHandler{service: newService, token: &token})
				if !loaded {
					go handler.ListRecommendations()
				}
				done, all := handler.GetProgress()
				if done < all {
					if setAuthorizationHeader(c, &token) {
						c.JSON(http.StatusOK, ListRecommendationsProgressResponse{int(done), int(all)})
					}
				} else {
					listRequestsInProcess.Delete(token)
					listResult, err := handler.GetResult()
					if err != nil {
						sendError(c, err)
					} else {
						if setAuthorizationHeader(c, &token) {
							c.JSON(http.StatusOK, ListRecommendationsResponse{
								Recommendations: listResult.Recommendations,
								FailedProjects:  listResult.FailedProjects})
						}
					}
				}
				return
			}
			log.Println(err)
		}
		// TODO send state field
		url := conf.AuthCodeURL(c.Request.URL.Path, authOptions...)
		c.Redirect(http.StatusSeeOther, url)
		return
	})

	/*
		router.POST("/recommendations/apply", func(c *gin.Context) {

		})

		router.GET("/recommendations/checkStatus", func(c *gin.Context) {

		})
	*/

	router.Run(":8000")

}
