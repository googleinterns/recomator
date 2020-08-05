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

import "google.golang.org/api/compute/v1"

// CreateSnapshot calls the disks.createSnapshot method.
// Requires compute.disks.createSnapshot or compute.snapshots.create permission.
func (s *googleService) CreateSnapshot(project, zone, disk, name string) error {
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
