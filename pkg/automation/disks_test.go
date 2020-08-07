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
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests if the generated names are different
func TestDifferentNames(t *testing.T) {
	generator := rand.New(rand.NewSource(0))
	testZone := "europe-west1-d"
	testDisk := "vertical-scaling-krzysztofk-wordpress"

	numberOfTrials := 100
	results := make([]string, numberOfTrials)

	for i := range results {
		result, err := randomSnapshotName(testZone, testDisk, generator)
		assert.Equal(t, err, nil, "Error should not happen when zone name is not to long")
		results[i] = result
	}

	sort.Strings(results)

	for i := 0; i < numberOfTrials-1; i++ {
		assert.NotEqual(t, results[i], results[i+1], "Generated shapshot names should be different")
	}
}

// Tests if the generated name is of a correct format
func TestCorrectName(t *testing.T) {
	generator := rand.New(rand.NewSource(0))
	testZone := "europe-west1-d"
	testDisk := "vertical-scaling-krzysztofk-wordpress"
	expectedDisk := testDisk[:min(maxDisknameLen, len(testDisk))]

	result, err := randomSnapshotName(testZone, testDisk, generator)
	assert.Equal(t, err, nil, "Error should not happen when zone name is not to long")

	index := 0
	assert.Equal(t, result[:min(maxDisknameLen, len(testDisk))], expectedDisk, "The first part of the snapshot name should  be a prefix of the disk name")

	index = min(maxDisknameLen, len(testDisk))
	assert.Equal(t, string(result[index]), "-", "The parts of snapshot name should be separeted with -")
	index = index + 1

	assert.Equal(t, result[index:index+len(testZone)], testZone, "The second part of snapshot name should be zone name.")

	index = index + len(testZone)
	assert.Equal(t, string(result[index]), "-", "The parts of snapshot name should be separeted with -")
	index = index + 1

	_, err = time.Parse(timestampFormat, result[index:index+len(timestampFormat)])
	assert.Equal(t, err, nil, "The format of the timestamp should be YYYYMMDDHHMMSS")

	index = index + len(timestampFormat)
	assert.Equal(t, string(result[index]), "-", "The parts of snapshot name should be separeted with -")

	assert.LessOrEqual(t, len(result), maxSnapshotnameLen, fmt.Sprintf("The length of the returned snapshot name must be less than %d", maxSnapshotnameLen))
}

// Tests that an error is returned if the zone name is long
func TestLongZoneName(t *testing.T) {
	generator := rand.New(rand.NewSource(0))
	testZone := strings.Repeat("a", 2*maxSnapshotnameLen)
	testDisk := ""

	_, err := randomSnapshotName(testZone, testDisk, generator)
	assert.NotEqual(t, err, nil, "When zone name is to long an error should be returned")
}

// Tests that the returned name is not to long if the disk name is long
func TestLongDiskName(t *testing.T) {
	generator := rand.New(rand.NewSource(0))
	testZone := ""
	testDisk := strings.Repeat("a", 2*maxSnapshotnameLen)

	result, err := randomSnapshotName(testZone, testDisk, generator)
	assert.Equal(t, err, nil, "Error should not happen when zone name is not to long")

	assert.LessOrEqual(t, len(result), maxSnapshotnameLen, fmt.Sprintf("The length of the returned snapshot name must be less than %d", maxSnapshotnameLen))
}
