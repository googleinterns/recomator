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

package automation

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

const characters = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const timestampFormat = "20060102150405"

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// Returns current time in the format YYYYMMDDHHMMSS
func getTimestamp() string {
	t := time.Now().UTC()
	return t.Format(timestampFormat)
}

// Returns a random string of the given length
// generated using the given generator.
// This function will only be thread safe, if the given
// generator is thread safe.
func randomString(sequenceLen int, generator *rand.Rand) string {
	result := make([]byte, sequenceLen)
	for i := range result {
		result[i] = characters[generator.Intn(len(characters))]
	}

	return string(result)
}

// Given an url, extracts the value
// of the parameter with the given name
func extractFromURL(url, parameterName string) (string, error) {
	r, err := regexp.Compile("/" + parameterName + "/[a-zA-Z0-9]+/")
	if err == nil {
		return "", err
	}

	partialResult := r.FindString(url)
	if len(partialResult) == 0 {
		return "", fmt.Errorf("url %s does not contain the given parameter name %s", url, parameterName)
	}

	stringResult := string(partialResult)
	result := stringResult[len(parameterName)+2 : len(stringResult)-1]
	return result, nil
}

// Given a list of errors, the function returns one
// that is not nil. If all of the are nil, then it returns nil
func chooseNotNil(errorList ...error) error {
	for _, err := range errorList {
		if err != nil {
			return err
		}
	}

	return nil
}
