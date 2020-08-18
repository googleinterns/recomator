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

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"
	"google.golang.org/api/serviceusage/v1"
)

// GoogleService is the inferface that prodives methods required to list recommendations and apply them
type GoogleService interface {
	// changes the machine type of an instance
	ChangeMachineType(project, zone, instance, machineType string) error

	// creates a snapshot of a disk
	CreateSnapshot(project, zone, disk, name string) error

	// deletes persistent disk
	DeleteDisk(project, zone, disk string) error

	// gets the specified instance resource
	GetInstance(project string, zone string, instance string) (*compute.Instance, error)

	// lists whether the requirements have been met for all APIs (APIs enabled).
	ListAPIRequirements(project string, apis []string) ([]Requirement, error)

	// lists whether the requirements have been met for all required permissions.
	ListPermissionRequirements(project string, permissions [][]string) ([]Requirement, error)

	// lists projects
	ListProjects() ([]string, error)

	// listing recommendations for specified project, zone and recommender
	ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error)

	// listing every zone available for the project methods
	ListZonesNames(project string) ([]string, error)

	// listing every region available for the project methods
	ListRegionsNames(project string) ([]string, error)

	// marks recommendation for the project with given etag and name claimed
	MarkRecommendationClaimed(name, etag string) (*gcloudRecommendation, error)

	// marks recommendation for the project with given etag and name succeeded
	MarkRecommendationSucceeded(name, etag string) (*gcloudRecommendation, error)

	// marks recommendation for the project with given etag and name failed
	MarkRecommendationFailed(name, etag string) (*gcloudRecommendation, error)

	// stops the specified instance
	StopInstance(project, zone, instance string) error

	// starts the specified instance
	StartInstance(project, zone, instance string) error
}

// googleService implements GoogleService interface for Recommender and Compute APIs.
type googleService struct {
	ctx                    context.Context
	computeService         *compute.Service
	recommenderService     *recommender.Service
	resourceManagerService *cloudresourcemanager.Service
	serviceUsageService    *serviceusage.Service
}

// NewGoogleService creates new googleServices.
// If creation failed the error will be non-nil.
func NewGoogleService(ctx context.Context) (GoogleService, error) {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	recommenderService, err := recommender.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resourceManagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return nil, err
	}

	serviceUsageService, err := serviceusage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	return &googleService{
		ctx:                    ctx,
		computeService:         computeService,
		recommenderService:     recommenderService,
		resourceManagerService: resourceManagerService,
		serviceUsageService:    serviceUsageService,
	}, nil
}
