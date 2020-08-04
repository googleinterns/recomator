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

import ("regexp")

// TODO change to returning error
func extractFromURL (url, parameterName string) string {
	r, err := regexp.Compile("/" + parameterName + "/[a-zA-Z0-9]+/")
	
	if err != nil {
		return ""
	}

	partialResult := r.FindString(url)

	if len(partialResult) == 0 {
		return ""
	}
	
	stringResult := string(partialResult)
	result := stringResult[len(parameterName) + 2: len(stringResult) - 1]
	
	return result
}