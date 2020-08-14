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
	"errors"
	"regexp"

	"google.golang.org/api/recommender/v1"
)

type gcloudValueMatcher = recommender.GoogleCloudRecommenderV1ValueMatcher

// Checks if the string toTest is equal to the string represented by value
// It is an error if value can't be interpreted as a string, unless
// value is nil. In that case true is returned.
func testValue(toTest string, value interface{}) (bool, error) {
	if value == nil {
		return true, nil
	}
	valueString, ok := value.(string)
	if !ok {
		return false, errors.New("if value is specified it must be of type string")
	}

	return valueString == toTest, nil
}

// Checks if the string toTest matches regex given by valueMatcher.MatchesPattern.
// If valueMatcher is nil then true is returned.
func testValueMatcher(toTest string, valueMatcher *gcloudValueMatcher) (bool, error) {
	if valueMatcher == nil {
		return true, nil
	}
	r, err := regexp.Compile("^" + valueMatcher.MatchesPattern + "$")
	if err != nil {
		return false, err
	}

	return r.MatchString(toTest), nil
}

// Checks if the string toTest matches the member MatchesPattern of valueMatcher,
// if valueMatcher is not nil. Otherwise, if value is not nil it is interpreted as string
// and equality of value.(string) and toTest is checked.
func testMatching(toTest string, value interface{}, valueMatcher *gcloudValueMatcher) (bool, error) {
	resultValue, err := testValue(toTest, value)
	if err != nil {
		return false, err
	}

	resultValueMatcher, err := testValueMatcher(toTest, valueMatcher)
	if err != nil {
		return false, err
	}

	return resultValue && resultValueMatcher, nil
}

// TestMachineType checks if the machine type of the instance specified by given project, zone and instance
// matches the given value or valueMatcher.
// According to Recommender API, in a test operation, either value or valueMatcher is specified.
// The value specified by the path field in the operation struct must match value or valueMatcher,
// depending on which one is defined. More can be read here:
// https://cloud.google.com/recommender/docs/reference/rest/v1/projects.locations.recommenders.recommendations#operation
func TestMachineType(service GoogleService, project string, zone string, instance string, value interface{}, valueMatcher *gcloudValueMatcher) (bool, error) {
	machineInstance, err := service.GetInstance(project, zone, instance)
	if err != nil {
		return false, err
	}

	return testMatching(machineInstance.MachineType, value, valueMatcher)
}

// TestStatus checks if the status of the instance specified by given project, zone and instance
// matches the given value or valueMatcher.
// According to Recommender API, in a test operation, either value or valueMatcher is specified.
// The value specified by the path field in the operation struct must match value or valueMatcher,
// depending on which one is defined. More can be read here:
// https://cloud.google.com/recommender/docs/reference/rest/v1/projects.locations.recommenders.recommendations#operation
func TestStatus(service GoogleService, project string, zone string, instance string, value interface{}, valueMatcher *gcloudValueMatcher) (bool, error) {
	machineInstance, err := service.GetInstance(project, zone, instance)
	if err != nil {
		return false, err
	}

	return testMatching(machineInstance.Status, value, valueMatcher)

}
