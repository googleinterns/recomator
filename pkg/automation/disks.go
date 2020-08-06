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

	"google.golang.org/api/compute/v1"
)

const maxDisknameLen = 20
const maxZonenameLen = 20
const maxSnapshotnameLen = 63

const characters = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// Returns current time in the format YYYYMMDDHHMMSS
func getTimestamp() string {
	t := time.Now().UTC()
	return t.Format("20060102150405")
}

// Returns a random string of the given length
func randomString(sequenceLen int) string {
	result := make([]byte, sequenceLen)
	for i := range result {
		result[i] = characters[generator.Intn(len(characters))]
	}

	return string(result)
}

// Returns the name of the snapshot, following
// the convention described here:
// https://cloud.google.com/compute/docs/disks/scheduled-snapshots#names_for_scheduled_snapshots
func randomSnapshotName(zone, disk string) string {
	result := ""
	result += disk[:min(maxDisknameLen, len(disk))]
	result += "-"
	result += zone[:min(maxZonenameLen, len(zone))]
	result += "-"
	result += getTimestamp()
	result += "-"
	randomSequenceLen := maxSnapshotnameLen - len(result)
	result += randomString(randomSequenceLen)

	return result
}

// CreateSnapshot calls the disks.createSnapshot method.
// Requires compute.disks.createSnapshot or compute.snapshots.create permission.
// For a given name, there can only be one snapshot having it.
// The maximum name length is 63.
func (s *googleService) CreateSnapshot(project, zone, disk, name string) error {
	if len(name) > maxSnapshotnameLen {
		return fmt.Errorf("Length of the snapshot name must not exceed %d", maxSnapshotnameLen)
	}
	disksService := compute.NewDisksService(s.computeService)
	snapshot := &compute.Snapshot{Name: name}
	_, err := disksService.CreateSnapshot(project, zone, disk, snapshot).Do()
	return err
}

// DeleteDisk calls the disks.delete method.
// Requires compute.disks.delete permission.
func (s *googleService) DeleteDisk(project, zone, disk string) error {
	disksService := compute.NewDisksService(s.computeService)
	_, err := disksService.Delete(project, zone, disk).Do()
	return err
}
