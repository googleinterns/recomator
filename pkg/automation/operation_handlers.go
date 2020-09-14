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
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/api/recommender/v1"
)

type gcloudOperation = recommender.GoogleCloudRecommenderV1Operation

// Assumes that the operation action is test.
// According to Recommender API, in a test operation, either value or valueMatcher is specified.
// The value specified by the path field in the operation struct must match value or valueMatcher,
// depending on which one is defined. More can be read here:
// https://cloud.google.com/recommender/docs/reference/rest/v1/projects.locations.recommenders.recommendations#operation
func testInstanceField(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectParam)
	zone, errZone := extractFromURL(path, zoneParam)
	instance, errInstance := extractFromURL(path, instanceParam)
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	machineInstance, err := service.GetInstance(project, zone, instance)
	if err != nil {
		return err
	}

	var result bool
	var field string

	switch operation.Path {
	case "/machineType":
		field = "machine type"
		result, err = testMatching(machineInstance.MachineType, operation.Value, operation.ValueMatcher)
	case "/status":
		field = "status"
		result, err = testMatching(machineInstance.Status, operation.Value, operation.ValueMatcher)
	default:
		return errors.New(operationNotSupportedMessage)
	}

	if err != nil {
		return err
	}

	if result == false {
		return fmt.Errorf("%s is not as expected", field)
	}

	return nil
}

// Assumes, that the operation's action is replace and path is /machineType.
// Replaces the machine type with a new one.
func replaceMachineType(service GoogleService, operation *gcloudOperation) error {
	path1 := operation.Resource
	path2, ok := operation.Value.(string)
	if !ok {
		return errors.New("wrong value type for operation replace machine type")
	}

	project, errProject := extractFromURL(path1, projectParam)
	instance, errInstance := extractFromURL(path1, instanceParam)

	machineType, errMachine := extractFromURL(path2, machineTypeParam)
	zone, errZone := extractFromURL(path2, zoneParam)
	err := chooseNotNil(errProject, errInstance, errMachine, errZone)
	if err != nil {
		return err
	}

	err = service.StopInstance(project, zone, instance)
	if err != nil {
		return err
	}

	err = service.ChangeMachineType(project, zone, instance, machineType)
	if err != nil {
		return err
	}

	return service.StartInstance(project, zone, instance)
}

// Assumes that operation's action is replace, path is status and value
// is terminated. Stops the given machine.
func stopInstance(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectParam)
	zone, errZone := extractFromURL(path, zoneParam)
	instance, errInstance := extractFromURL(path, instanceParam)
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	return service.StopInstance(project, zone, instance)
}

// Assumes that operation's action is add, and ResourceType
// is compute.googleapis.com/Snapshot. Adds a snapshot of the given machine.
func addSnapshot(service GoogleService, operation *gcloudOperation) error {
	value, ok := operation.Value.(map[string]interface{})

	if !ok {
		return errors.New("wrong value type for operation add snapshot")
	}

	path, ok := value["source_disk"].(string)
	if !ok {
		return fmt.Errorf("wrong source disk type for operation add snapshot: %t", value["source_disk"])
	}

	project, errProject := extractFromURL(path, projectParam)
	zone, errZone := extractFromURL(path, zoneParam)
	disk, errDisk := extractFromURL(path, diskParam)
	err := chooseNotNil(errProject, errZone, errDisk)
	if err != nil {
		return err
	}

	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	name, err := randomSnapshotName(zone, disk, generator)

	if err != nil {
		return err
	}

	return service.CreateSnapshot(project, zone, disk, name)
}

// Assumes that the operation's action is remove and its resource type
// is compute.googleapis.com/Disk. Removes the given disk.
func removeDisk(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectParam)
	zone, errZone := extractFromURL(path, zoneParam)
	disk, errDisk := extractFromURL(path, diskParam)
	err := chooseNotNil(errProject, errZone, errDisk)
	if err != nil {
		return err
	}

	return service.DeleteDisk(project, zone, disk)
}
