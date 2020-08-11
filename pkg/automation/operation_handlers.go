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
	StorageLocations string
}

type gcloudOperation = recommender.GoogleCloudRecommenderV1Operation

func (s *googleService) testMachineType(operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, "projects")
	zone, errZone := extractFromURL(path, "zones")
	instance, errInstance := extractFromURL(path, "instance")
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	result, err := s.TestMachineType(project, zone, instance, operation.Value, operation.ValueMatcher)
	if err != nil {
		return err
	}

	if result == false {
		return errors.New("testing of the machine type failed")
	}

	return nil
}

func (s *googleService) testStatus(operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, "projects")
	zone, errZone := extractFromURL(path, "zones")
	instance, errInstance := extractFromURL(path, "instance")
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	result, err := s.TestStatus(project, zone, instance, operation.Value, operation.ValueMatcher)
	if err != nil {
		return err
	}

	if result == false {
		return errors.New("testing of the status failed")
	}

	return nil
}

func (s *googleService) replaceMachineType(operation *gcloudOperation) error {
	path1 := operation.Resource
	path2, ok := operation.Value.(string)
	if !ok {
		return errors.New("wrong value type for operation replace machine type")
	}

	project, errProject := extractFromURL(path1, "projects")
	instance, errInstance := extractFromURL(path1, "instances")
	machineType, errMachine := extractFromURL(path2, "machineTypes")
	zone, errZone := extractFromURL(path2, "zones")
	err := chooseNotNil(errProject, errInstance, errMachine, errZone)
	if err != nil {
		return err
	}

	return s.ChangeMachineType(project, zone, instance, machineType)
}

func (s *googleService) replaceStatus(operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, "projects")
	zone, errZone := extractFromURL(path, "zones")
	instance, errInstance := extractFromURL(path, "instance")
	err := chooseNotNil(errProject, errZone, errInstance)
	if err != nil {
		return err
	}

	return s.StopInstance(project, zone, instance)
}

func (s *googleService) addSnapshot(operation *gcloudOperation) error {
	value, ok := operation.Value.(valueAddSnapshot)
	if !ok {
		return errors.New("wrong value type for operation add snapshot")
	}
	path := value.SourceDisk

	project, errProject := extractFromURL(path, "projects")
	zone, errZone := extractFromURL(path, "zones")
	disk, errDisk := extractFromURL(path, "disks")
	err := chooseNotNil(errProject, errZone, errDisk)
	if err != nil {
		return err
	}

	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	name, err := randomSnapshotName(zone, disk, generator)
	if err != nil {
		return err
	}

	return s.CreateSnapshot(project, zone, disk, name)
}

func (s *googleService) removeDisk(operation *gcloudOperation) error {
	path := operation.Resource

	project, errProject := extractFromURL(path, "projects")
	zone, errZone := extractFromURL(path, "zones")
	disk, errDisk := extractFromURL(path, "disks")
	err := chooseNotNil(errProject, errZone, errDisk)
	if err != nil {
		return err
	}

	return s.DeleteDisk(project, zone, disk)
}
