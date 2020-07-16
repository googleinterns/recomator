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
	"context"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"
)

// GoogleService is the inferface that prodives the following methods:
// ChangeMachineType - changes the machine type of an instance,
// ListRecommendations - listing recommendations for specified project, zone and recommender,
// ListZonesNames - listing every zone available for the project methods,
// StopInstance - stops the specified instance.
type GoogleService interface {
	ChangeMachineType(project string, zone string, instance string, machineType string) error

	ListRecommendations(project string, location string, recommenderID string) ([]*gcloudRecommendation, error)

	ListZonesNames(project string) ([]string, error)

	StopInstance(project string, zone string, instance string) error
}

// googleService implements GoogleService interface for Recommender and Compute APIs,
// using projects.locations.recommenders.recommendations/list and zones/list methods.
type googleService struct {
	ctx                context.Context
	instancesService   *compute.InstancesService
	zonesService       *compute.ZonesService
	recommenderService *recommender.ProjectsLocationsRecommendersRecommendationsService
}

// NewGoogleService creates new googleServices.
// If creation failed the error will be non-nil.
func NewGoogleService(ctx context.Context) (GoogleService, error) {
	recService, err := recommender.NewService(ctx)
	if err != nil {
		return nil, err
	}

	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	return &googleService{
		ctx:                ctx,
		instancesService:   compute.NewInstancesService(computeService),
		recommenderService: recommender.NewProjectsLocationsRecommendersRecommendationsService(recService),
		zonesService:       compute.NewZonesService(computeService),
	}, nil
}
