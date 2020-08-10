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
	"math/rand"
	"strings"
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

type gcloudOperationGroup = recommender.GoogleCloudRecommenderV1OperationGroup

// DoOperation does the action specified in the operation.
func (s *googleService) DoOperation(operation *recommender.GoogleCloudRecommenderV1Operation) error {
	switch strings.ToLower(operation.Action) {
	case "test":
		switch operation.Path {
		case "/machineType":
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			instance := extractFromURL(path, "instance")

			result, err := s.TestMachineType(project, zone, instance, operation.Value, operation.ValueMatcher)
			if err != nil {
				return err
			}

			if result == false {
				return errors.New("testing of the machine type failed")
			}
		case "/status":
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			instance := extractFromURL(path, "instance")

			result, err := s.TestStatus(project, zone, instance, operation.Value, operation.ValueMatcher)
			if err != nil {
				return err
			}

			if result == false {
				return errors.New("testing of the status failed")
			}
		default:
			return errors.New("the opperation is not supported")
		}
	case "replace":
		switch operation.Path {
		case "/machineType":
			path1 := operation.Resource
			path2, ok := operation.Value.(string)
			if !ok {
				return errors.New("wrong value type for operation replace machine type")
			}

			project := extractFromURL(path1, "projects")
			instance := extractFromURL(path1, "instances")
			machineType := extractFromURL(path2, "machineTypes")
			zone := extractFromURL(path2, "zones")

			s.ChangeMachineType(project, zone, instance, machineType)
		case "/status":
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			instance := extractFromURL(path, "instance")

			err := s.StopInstance(project, zone, instance)
			if err != nil {
				return err
			}
		default:
			return errors.New("the opperation is not supported")
		}
	case "add":
		switch operation.ResourceType {
		case "compute.googleapis.com/Snapshot":
			value, ok := operation.Value.(valueAddSnapshot)
			if !ok {
				return errors.New("wrong value type for operation add snapshot")
			}
			path := value.SourceDisk

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			disk := extractFromURL(path, "disks")

			generator := rand.New(rand.NewSource(time.Now().UnixNano()))
			name, err := randomSnapshotName(zone, disk, generator)
			if err != nil {
				return err
			}

			err = s.CreateSnapshot(project, zone, disk, name)
			if err != nil {
				return err
			}
		default:
			return errors.New("the opperation is not supported")
		}

	case "remove":
		switch operation.ResourceType {
		case "compute.googleapis.com/Disk":
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			disk := extractFromURL(path, "disks")

			err := s.DeleteDisk(project, zone, disk)
			if err != nil {
				return err
			}

		default:
			return errors.New("the opperation is not supported")
		}

	default:
		return errors.New("the opperation is not supported")
	}

	return nil
}

// Apply is the method used to apply recommendations from Recommender API.
// Supports recommendations from the following recommenders:
// - google.compute.disk.IdleResourceRecommender
// - google.compute.instance.IdleResourceRecommender
// - google.compute.instance.MachineTypeRecommender
func (s *googleService) Apply(recommendation *gcloudRecommendation) error {
	// check the state is ACTIVE
	// claim the recommendation
	// this may somehow be concurrent
	// if test fails, just proceed to the next group? Or what?
	for _, operationGroup := range recommendation.Content.OperationGroups {
		for _, operation := range operationGroup.Operations {
			err := s.DoOperation(operation)
			if err != nil {
				// mark failed
				return err
			}
		}
	}
	// mark succedeed
	return nil
}
