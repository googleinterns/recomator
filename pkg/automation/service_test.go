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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
)

func TestAwaitCompletion(t *testing.T) {
	// fast failure
	err := AwaitCompletion(fastFailingOperationGen)
	assert.EqualError(t, err, "Oh no")

	//ending in success
	longOperationGenCounter = 0
	longOperationGenShouldSucceed = true
	err = AwaitCompletion(longOperationGen)
	assert.Nil(t, err)
	assert.Equal(t, longOperationGenCounter, 3)

	//ending in failure
	longOperationGenCounter = 0
	longOperationGenShouldSucceed = false
	err = AwaitCompletion(longOperationGen)
	assert.EqualError(t, err, "Oh no")
	assert.Equal(t, longOperationGenCounter, 3)
}

func fastFailingOperationGen() (*compute.Operation, error) {
	return nil, errors.New("Oh no")
}

var longOperationGenCounter int
var longOperationGenShouldSucceed bool

func longOperationGen() (*compute.Operation, error) {
	longOperationGenCounter++
	switch longOperationGenCounter {
	case 1, 2:
		return &compute.Operation{Status: "PROCESSING"}, nil
	case 3:
		if longOperationGenShouldSucceed {
			return &compute.Operation{Status: "DONE"}, nil
		}
		return &compute.Operation{Status: "Done"}, errors.New("Oh no")
	default:
		panic("too many calls")
	}
}
