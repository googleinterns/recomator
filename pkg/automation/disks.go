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
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/compute/v1"
)

const maxDisknameLen = 20
const maxZonenameLen = 26
const maxSnapshotnameLen = 63

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

// Returns the name of the snapshot, generated using the given generator
// following the convention described here:
// https://cloud.google.com/compute/docs/disks/scheduled-snapshots#names_for_scheduled_snapshots
// This function will only be thread safe, if the given
// generator is thread safe.
func randomSnapshotName(zone string, disk string, generator *rand.Rand) (string, error) {
	if len(zone) > maxZonenameLen {
		return "", fmt.Errorf("length of the zone name must not exceed %d", maxZonenameLen)
	}

	result := ""
	result += disk[:min(maxDisknameLen, len(disk))]
	result += "-"
	result += zone
	result += "-"
	result += getTimestamp()
	result += "-"
	randomSequenceLen := maxSnapshotnameLen - len(result)
	result += randomString(randomSequenceLen, generator)

	return result, nil
}

// CreateSnapshot calls the disks.createSnapshot method.
// Requires compute.disks.createSnapshot or compute.snapshots.create permission.
// For a given name, there can only be one snapshot having it.
// The maximum name length is 63.
func (s *googleService) CreateSnapshot(project, zone, disk, name string) error {
	if len(name) > maxSnapshotnameLen {
		return fmt.Errorf("length of the snapshot name must not exceed %d", maxSnapshotnameLen)
	}
	disksService := compute.NewDisksService(s.computeService)
	snapshot := &compute.Snapshot{Name: name}
	requestID := uuid.New().String()
	err := AwaitUntilCompletion(func() (*compute.Operation, error) {
		return disksService.CreateSnapshot(project, zone, disk, snapshot).RequestId(requestID).Do()
	})
	return err
}

// DeleteDisk calls the disks.delete method.
// Requires compute.disks.delete permission.
func (s *googleService) DeleteDisk(project, zone, disk string) error {
	disksService := compute.NewDisksService(s.computeService)
	requestID := uuid.New().String()
	err := AwaitUntilCompletion(func() (*compute.Operation, error) {
		return disksService.Delete(project, zone, disk).RequestId(requestID).Do()
	})
	return err
}
