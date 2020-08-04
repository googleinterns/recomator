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
	"strings"

	"google.golang.org/api/recommender/v1"
)

type valueAddSnapshot struct {
	Name              string
	Source_disk       string
	Storage_locations string
}

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

			}

			if result == false {

			}
		case "/status":
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			instance := extractFromURL(path, "instance")

			result, err := s.TestStatus(project, zone, instance, operation.Value, operation.ValueMatcher)
			if err != nil {

			}

			if result == false {

			}
		default:
			// return error operation not supported
		}
	case "replace":
		switch operation.Path {
		case "/machineType":
			// TODO do both zones need to be the same?
			path1 := operation.Resource
			path2, ok := operation.Value.(string)
			if !ok {
				// handle error (if nil, it's fine)
			}

			project := extractFromURL(path1, "projects")
			instance := extractFromURL(path1, "instances")
			machineType := extractFromURL(path2, "machineTypes")
			zone := extractFromURL(path2, "zones")

			s.ChangeMachineType(project, zone, instance, machineType)
		case "/status":
			// stop (or start?) the machine
			// TODO start?
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			instance := extractFromURL(path, "instance")

			s.StopInstance(project, zone, instance)

		default:
			// return error operation not supported
		}
	case "add":
		switch operation.ResourceType {
		case "compute.googleapis.com/Snapshot":
			value, ok := operation.Value.(valueAddSnapshot)
			if !ok {
				// handle error (if nil, it's fine)
			}
			path := value.Source_disk

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			disk := extractFromURL(path, "disks")

			s.CreateSnapshot(project, zone, disk)
		default:
			//
		}

	case "remove":
		switch operation.ResourceType {
		case "compute.googleapis.com/Disk":
			path := operation.Resource

			project := extractFromURL(path, "projects")
			zone := extractFromURL(path, "zones")
			disk := extractFromURL(path, "disks")

			s.DeleteDisk(project, zone, disk)

		default:
			//
		}

	default:
		// return error operation not supported
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
