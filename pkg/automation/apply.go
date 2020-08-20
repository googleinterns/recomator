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

	"google.golang.org/api/recommender/v1"
)

type gcloudOperationGroup = recommender.GoogleCloudRecommenderV1OperationGroup

const (
	operationNotSupportedMessage = "the operation is not supported"
)

const (
	projectParam     = "projects"
	zoneParam        = "zones"
	instanceParam    = "instances"
	diskParam        = "disks"
	machineTypeParam = "machineTypes"
)

// DoOperation does the action specified in the operation.
func DoOperation(service GoogleService, operation *gcloudOperation) error {
	switch strings.ToLower(operation.Action) {
	case "test":
		if operation.ResourceType != "compute.googleapis.com/Instance" {
			return errors.New(operationNotSupportedMessage)
		}
		return testInstanceField(service, operation)
	case "replace":
		if operation.ResourceType != "compute.googleapis.com/Instance" {
			return errors.New(operationNotSupportedMessage)
		}
		switch operation.Path {
		case "/machineType":
			return replaceMachineType(service, operation)
		case "/status":
			if operation.Value != "TERMINATED" {
				return errors.New(operationNotSupportedMessage)
			}

			return stopInstance(service, operation)
		default:
			return errors.New(operationNotSupportedMessage)
		}
	case "add":
		switch operation.ResourceType {
		case "compute.googleapis.com/Snapshot":
			return addSnapshot(service, operation)
		default:
			return errors.New(operationNotSupportedMessage)
		}

	case "remove":
		switch operation.ResourceType {
		case "compute.googleapis.com/Disk":
			return removeDisk(service, operation)
		}
	}

	return errors.New(operationNotSupportedMessage)
}

// DoOperations calls DoOperation for each operation specified in the recommendation
func DoOperations(service GoogleService, recommendation *gcloudRecommendation) error {
	for _, operationGroup := range recommendation.Content.OperationGroups {
		for _, operation := range operationGroup.Operations {
			err := DoOperation(service, operation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Apply is the method used to apply recommendations from Recommender API.
// Supports recommendations from the following recommenders:
// - google.compute.disk.IdleResourceRecommender
// - google.compute.instance.IdleResourceRecommender
// - google.compute.instance.MachineTypeRecommender
func Apply(service GoogleService, recommendation *gcloudRecommendation) error {
	if strings.ToLower(recommendation.StateInfo.State) != "active" {
		return errors.New("to apply a recommendation, its status must be active")
	}

	newRecommendation, err := service.MarkRecommendationClaimed(recommendation.Name, recommendation.Etag)
	if err != nil {
		return err
	}
	*recommendation = *newRecommendation

	err = DoOperations(service, recommendation)
	if err != nil {
		newRecommendation, errMark := service.MarkRecommendationFailed(recommendation.Name, recommendation.Etag)
		if errMark != nil {
			return errMark
		}
		*recommendation = *newRecommendation

		return err
	}

	newRecommendation, err = service.MarkRecommendationSucceeded(recommendation.Name, recommendation.Etag)
	if err != nil {
		return err
	}
	*recommendation = *newRecommendation

	return nil
}
