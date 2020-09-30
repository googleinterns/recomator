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
	"io/ioutil"

	"github.com/coreos/go-oidc"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// NewConfig creates the configuration of oauth2 using given clientID, clientSecret and redirectURL.
// Scopes and provider are set by default.
func NewConfig(file string) (*oauth2.Config, error) {
	config := &oauth2.Config{
		Scopes: []string{oidc.ScopeOpenID, "https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/userinfo.email"},
	}
	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var data map[string]string
	if err := json.Unmarshal(byt, &data); err != nil {
		return nil, err
	}
	provider, err := oidc.NewProvider(oauth2.NoContext, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	config.Endpoint = provider.Endpoint()
	config.ClientID = data["clientID"]
	config.ClientSecret = data["clientSecret"]
	config.RedirectURL = data["redirectURL"]
	return config, nil
}
