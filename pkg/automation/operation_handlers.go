package automation

import (
	"errors"
	"math/rand"
	"time"

	"google.golang.org/api/recommender/v1"
)

// The type, that the value field of operation should be
// interpretable as in add snapshot operation
type valueAddSnapshot struct {
	Name             string
	SourceDisk       string
	StorageLocations []string
}

type gcloudOperation = recommender.GoogleCloudRecommenderV1Operation

func testMachineType(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectPath)
	zone, errZone := extractFromURL(path, zonePath)
	instance, errInstance := extractFromURL(path, instancePath)
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	result, err := service.TestMachineType(project, zone, instance, operation.Value, operation.ValueMatcher)
	if err != nil {
		return err
	}

	if result == false {
		return errors.New("machine type is not as expected")
	}

	return nil
}

func testStatus(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectPath)
	zone, errZone := extractFromURL(path, zonePath)
	instance, errInstance := extractFromURL(path, instancePath)
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	result, err := service.TestStatus(project, zone, instance, operation.Value, operation.ValueMatcher)
	if err != nil {
		return err
	}

	if result == false {
		return errors.New("status of instance is not as expected")
	}

	return nil
}

func replaceMachineType(service GoogleService, operation *gcloudOperation) error {
	path1 := operation.Resource
	path2, ok := operation.Value.(string)
	if !ok {
		return errors.New("wrong value type for operation replace machine type")
	}

	project, errProject := extractFromURL(path1, projectPath)
	instance, errInstance := extractFromURL(path1, instancePath)

	machineType, errMachine := extractFromURL(path2, machineTypePath)
	zone, errZone := extractFromURL(path2, zonePath)
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

func stopInstance(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectPath)
	zone, errZone := extractFromURL(path, zonePath)
	instance, errInstance := extractFromURL(path, instancePath)
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	return service.StopInstance(project, zone, instance)
}

func addSnapshot(service GoogleService, operation *gcloudOperation) error {
	value, ok := operation.Value.(valueAddSnapshot)
	if !ok {
		return errors.New("wrong value type for operation add snapshot")
	}
	path := value.SourceDisk

	project, errProject := extractFromURL(path, projectPath)
	zone, errZone := extractFromURL(path, zonePath)
	disk, errDisk := extractFromURL(path, diskPath)
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

func removeDisk(service GoogleService, operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, projectPath)
	zone, errZone := extractFromURL(path, zonePath)
	disk, errDisk := extractFromURL(path, diskPath)
	err := chooseNotNil(errProject, errZone, errDisk)
	if err != nil {
		return err
	}

	return service.DeleteDisk(project, zone, disk)
}
