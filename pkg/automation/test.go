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
	"regexp"

	"google.golang.org/api/recommender/v1"
)

type gcloudValueMatcher = recommender.GoogleCloudRecommenderV1ValueMatcher

// Checks if the string toTest matches the member MatchesPattern of valueMatcher
// If valueMatcher is not nil. Otherwise, if value is not nil it is interpreted as string
// And equality of value.(string) and toTest is checked. If both value and valueMatcher are nil,
// then it results in an error
func test(toTest string, value interface{}, valueMatcher *gcloudValueMatcher) (bool, error) {
	if value == nil {
		if valueMatcher == nil {
			return false, operationUnsupportedError()
		}

		r, err := regexp.Compile("^" + valueMatcher.MatchesPattern + "$")
		if err != nil {
			return false, err
		}

		return r.MatchString(toTest), nil
	}

	valueString, ok := value.(string)
	if !ok {
		return false, typeError()
	}

	return valueString == toTest, nil
}

// Checks if the machine type of the instance specified by given project, zone and instance
// matches the given value or valueMatcher
func (s *googleService) TestMachineType(project string, zone string, instance string, value interface{}, valueMatcher *gcloudValueMatcher) (bool, error) {
	machineInstance, err := s.GetInstance(project, zone, instance)
	if err != nil {
		return false, err
	}

	machineType := machineInstance.MachineType

	return test(machineType, value, valueMatcher)
}

// Checks if the status of the instance specified by given project, zone and instance
// matches the given value or valueMatcher
func (s *googleService) TestStatus(project string, zone string, instance string, value interface{}, valueMatcher *gcloudValueMatcher) (bool, error) {
	machineInstance, err := s.GetInstance(project, zone, instance)
	if err != nil {
		return false, err
	}

	status := machineInstance.Status

	return test(status, value, valueMatcher)

}
