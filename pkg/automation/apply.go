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
	"strings"
)

// DoOperation does the action specified in the operation.
func (s *googleService) DoOperation(operation *gcloudOperation) error {
	switch strings.ToLower(operation.Action) {
	case "test":
		switch operation.Path {
		case "/machineType":
			return s.testMachineType(operation)
		case "/status":
			return s.testStatus(operation)
		default:
			return errors.New("the opperation is not supported")
		}
	case "replace":
		switch operation.Path {
		case "/machineType":
			return s.replaceMachineType(operation)
		case "/status":
			return s.replaceStatus(operation)
		default:
			return errors.New("the opperation is not supported")
		}
	case "add":
		switch operation.ResourceType {
		case "compute.googleapis.com/Snapshot":
			return s.addSnapshot(operation)
		default:
			return errors.New("the opperation is not supported")
		}

	case "remove":
		switch operation.ResourceType {
		case "compute.googleapis.com/Disk":
			return s.removeDisk(operation)
		default:
			return errors.New("the opperation is not supported")
		}

	default:
		return errors.New("the opperation is not supported")
	}
}

// Apply is the method used to apply recommendations from Recommender API.
// Supports recommendations from the following recommenders:
// - google.compute.disk.IdleResourceRecommender
// - google.compute.instance.IdleResourceRecommender
// - google.compute.instance.MachineTypeRecommender
func (s *googleService) Apply(recommendation *gcloudRecommendation) error {
	// check that status is active
	s.MarkRecommendationSucceded(recommendation.Name, recommendation.Etag)
	for _, operationGroup := range recommendation.Content.OperationGroups {
		for _, operation := range operationGroup.Operations {
			err := s.DoOperation(operation)
			if err != nil {
				s.MarkRecommendationFailed(recommendation.Name, recommendation.Etag)
				return err
			}
		}
	}
	s.MarkRecommendationSucceded(recommendation.Name, recommendation.Etag)

	return nil
}
