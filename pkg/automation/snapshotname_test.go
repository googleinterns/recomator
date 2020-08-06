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
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests if the generated names are different
func TestDifferentNames(t *testing.T) {
	testZone := strings.Repeat("a", 100)
	testDiskName := strings.Repeat("a", 100)

	numberOfTrials := 10000
	results := make([]string, numberOfTrials)

	for i := range results {
		results[i] = randomSnapshotName(testZone, testDiskName)
	}

	sort.Strings(results)

	for i := 0; i < numberOfTrials-1; i++ {
		assert.NotEqual(t, results[i], results[i+1], "Generated shapshot names should be different")
	}
}

// Tests that the returned name is not to long if the zone name is long
func TestLongZoneName(t *testing.T) {
	testZone := strings.Repeat("a", 2*maxSnapshotnameLen)
	testDiskName := ""

	result := randomSnapshotName(testZone, testDiskName)

	assert.LessOrEqual(t, len(result), maxSnapshotnameLen, fmt.Sprintf("The length of the returned snapshot name must be less than %d", maxSnapshotnameLen))
}

// Tests that the returned name is not to long if the disk name is long
func TestLongDiskName(t *testing.T) {
	testZone := strings.Repeat("a", 2*maxSnapshotnameLen)
	testDiskName := ""

	result := randomSnapshotName(testZone, testDiskName)

	assert.LessOrEqual(t, len(result), maxSnapshotnameLen, fmt.Sprintf("The length of the returned snapshot name must be less than %d", maxSnapshotnameLen))
}
