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

	ListRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation

	ListZonesNames(project string) []string

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
