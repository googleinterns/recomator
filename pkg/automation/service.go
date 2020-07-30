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

// GoogleService is the inferface that prodives methods required to list recommendations and apply them
type GoogleService interface {
	// changes the machine type of an instance
	ChangeMachineType(project, zone, instance, machineType string) error

	// creates a snapshot of a disk
	CreateSnapshot(project, location, disk, name string) error

	// delete persistent disk
	DeleteDisk(project, location, disk string) error

	// listing recommendations for specified project, zone and recommender
	ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error)

	// listing every zone available for the project methods
	ListZonesNames(project string) ([]string, error)

	// listing every region available for the project methods
	ListRegionsNames(project string) ([]string, error)

	// stops the specified instance
	StopInstance(project, zone, instance string) error
}

// googleService implements GoogleService interface for Recommender and Compute APIs,
// using projects.locations.recommenders.recommendations/list and zones/list methods.
type googleService struct {
	ctx                context.Context
	computeService     *compute.Service
	recommenderService *recommender.Service
}

// NewGoogleService creates new googleServices.
// If creation failed the error will be non-nil.
func NewGoogleService(ctx context.Context) (GoogleService, error) {
	recommenderService, err := recommender.NewService(ctx)
	if err != nil {
		return nil, err
	}

	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	return &googleService{
		ctx:                ctx,
		computeService:     computeService,
		recommenderService: recommenderService,
	}, nil
}
