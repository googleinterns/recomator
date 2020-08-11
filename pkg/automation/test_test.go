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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing using correct value results in a match
func TestCorrectValue(t *testing.T) {
	toTest := "aaaaa"
	value := "aaaaa"
	var valueMatcher *gcloudValueMatcher = nil
	result, err := testMatching(toTest, value, valueMatcher)
	assert.True(t, result, "When value is equal to the tested value then they match")
	assert.Nil(t, err)
}

// Testing using correct value doesn't result in a match
func TestIncorrectValue(t *testing.T) {
	toTest := "aaaa"
	value := "aaaaa"
	var valueMatcher *gcloudValueMatcher = nil
	result, err := testMatching(toTest, value, valueMatcher)
	assert.False(t, result, "When value is not equal to the tested value then they don't match")
	assert.Nil(t, err)
}

// Testing using correct value matcher results in a match
func TestCorrectValueMatcher(t *testing.T) {
	toTest := "aaaa"
	var value interface{} = nil
	valueMatcher := gcloudValueMatcher{
		MatchesPattern: "a*",
	}
	result, err := testMatching(toTest, value, &valueMatcher)
	assert.True(t, result, "When value matcher matches the tested value then they match")
	assert.Nil(t, err)
}

// Testing using correct value matcher doesn't result in a match
func TestIncorrectValueMatcher(t *testing.T) {
	toTest := "aaaa"
	var value interface{} = nil
	valueMatcher := gcloudValueMatcher{
		MatchesPattern: "ba*",
	}
	result, err := testMatching(toTest, value, &valueMatcher)
	assert.False(t, result, "When value matcher doesn't match the tested value then they don't match")
	assert.Nil(t, err)
}

// Testing using correct value and matching value matcher results in a match
func TestCorrectValueAndValueMatcher(t *testing.T) {
	toTest := "aaaa"
	value := "aaaa"
	valueMatcher := gcloudValueMatcher{
		MatchesPattern: "a*",
	}
	result, err := testMatching(toTest, value, &valueMatcher)
	assert.True(t, result, "When value matcher and value matches the tested value then they match")
	assert.Nil(t, err)
}

// Testing using correct value and not matching value matcher doesn't result in a match
func TestCorrectValueIncorrectValueMatcher(t *testing.T) {
	toTest := "aaaa"
	value := "aaaa"
	valueMatcher := gcloudValueMatcher{
		MatchesPattern: "ba*",
	}
	result, err := testMatching(toTest, value, &valueMatcher)
	assert.False(t, result, "When value matcher doesn't match the tested value and value matches the tested value then they don't match")
	assert.Nil(t, err)
}

// Testing using incorrect value and matching value matcher doesn't result in a match
func TestIncorrectValueCorrectValueMatcher(t *testing.T) {
	toTest := "aaaa"
	value := "aaa"
	valueMatcher := gcloudValueMatcher{
		MatchesPattern: "a*",
	}
	result, _ := testMatching(toTest, value, &valueMatcher)
	assert.False(t, result, "When value matcher matches the tested value and value doesn't match the tested value then they don't match")
}

// Testing using nil as value and value matcher results in a match
func TestNilValueAndValueMatcher(t *testing.T) {
	toTest := "aaaa"
	var value interface{} = nil
	var valueMatcher *gcloudValueMatcher = nil
	result, err := testMatching(toTest, value, valueMatcher)
	assert.True(t, result, "When value matcher and value are nil, then everything matches")
	assert.Nil(t, err)
}

// Testing using non-string value results in an error
func TestNonstringValue(t *testing.T) {
	toTest := "aaaa"
	value := 123
	var valueMatcher *gcloudValueMatcher = nil
	_, err := testMatching(toTest, value, valueMatcher)
	assert.Errorf(t, err, "If the given value is not a string, an error should be returned")
}
