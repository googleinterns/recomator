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
	"io/ioutil"
	"log"
	"os"

	"github.com/googleinterns/recomator/pkg/server"
)

func main() {
	byt, err := ioutil.ReadFile(file)
	var clientID, clientSecret, redirectURL string
	if err == nil {
		var data map[string]string
		if err := json.Unmarshal(byt, &data); err != nil {
			return nil, err
		}
		clientID = data["clientID"]
		clientSecret = data["clientSecret"]
		redirectURL = data["redirectURL"]
	} else {
		clientID := os.Getenv("CLIENT_ID")
		clientSecret := os.Getenv("CLIENT_SECRET")
		redirectURL := os.Getenv("REDIRECT_URL")
	}

	if conf, err := server.NewConfig(&config, clientID, clientSecret, redirectURL); err != nil {
		log.Fatal(err)
	}

	service, err := server.NewSharedService(*conf)
	if err != nil {
		log.Fatal(err)
	}
	router := server.SetUpRouter(service)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}
