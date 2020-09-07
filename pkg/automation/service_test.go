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
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
)

func TestAwaitCompletionFastFailure(t *testing.T) {
	err := AwaitCompletion(fastFailingOperationGen, time.Millisecond)
	assert.EqualError(t, err, "Oh no")
}

func TestAwaitingCompletionEndsWithSuccess(t *testing.T) {
	calledTimes := 0
	err := AwaitCompletion(longOperationGen(true, &calledTimes), time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, calledTimes, 3)
}
func TestAwaitingCompletionEndsWithFailure(t *testing.T) {
	calledTimes := 0
	err := AwaitCompletion(longOperationGen(false, &calledTimes), time.Millisecond)
	assert.EqualError(t, err, "Oh no")
	assert.Equal(t, calledTimes, 3)
}

func fastFailingOperationGen() (*compute.Operation, error) {
	return nil, errors.New("Oh no")
}

func longOperationGen(succeeds bool, genCounter *int) operationGenerator {
	*genCounter = 0
	return func() (*compute.Operation, error) {
		(*genCounter)++
		switch *genCounter {
		case 1, 2:
			return &compute.Operation{Status: "PROCESSING"}, nil
		case 3:
			if succeeds {
				return &compute.Operation{Status: "DONE"}, nil
			}
			return &compute.Operation{Status: "Done"}, errors.New("Oh no")
		default:
			panic("too many calls")
		}
	}
}
