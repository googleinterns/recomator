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
	"time"

	"github.com/xyproto/randomstring"
	"google.golang.org/api/compute/v1"
)

const maxDisknameLen = 20
const maxZonenameLen = 26
const maxSnapshotnameLen = 63

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// Returns current time in the fromat YYYYMMDDHHMMSS
func getTimestamp() string {
	t := time.Now().UTC()
	return t.Format("20060102150405")
}

// Returns the name of the snapshot, following
// the convention described here:
// https://cloud.google.com/compute/docs/disks/scheduled-snapshots#names_for_scheduled_snapshots
func getSnapshotName(zone, disk string) string {
	result := ""
	result += disk[:min(maxDisknameLen, len(disk))]
	result += "-"
	result += zone[:min(maxZonenameLen, len(zone))]
	result += "-"
	result += getTimestamp()
	result += "-"
	randomSequenceLen := maxSnapshotnameLen - len(result)
	result += randomstring.HumanFriendlyString(randomSequenceLen)

	return result
}

// CreateSnapshot calls the disks.createSnapshot method.
// Requires compute.disks.createSnapshot or compute.snapshots.create permission.
func (s *googleService) CreateSnapshot(project, zone, disk string) error {
	disksService := compute.NewDisksService(s.computeService)
	name := getSnapshotName(zone, disk)
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
