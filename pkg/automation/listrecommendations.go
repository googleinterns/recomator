package automation

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"
)

// Service is the inferface that prodives the following methods:
// listZoneRecommendations - listing recommendations for specified project, zone and reccommender,
// listZonesNames - listing every zone available for the project methods,
// getNumConcurrentCalls returns the recommended maximum number of concurrent calls to listZonesNames.
type Service interface {
	listZoneRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation

	listZonesNames(project string) []string

	getNumConcurrentCalls() int
}

// ZonesRecommendationsService implements Service interface for Recommender and Compute APIs,
// using projects.locations.recommenders.recommendations/list and zones/list methods.
type ZonesRecommendationsService struct {
	ctx                   context.Context
	recommendationService *recommender.ProjectsLocationsRecommendersRecommendationsService
	zonesService          *compute.ZonesService
	numConcurrentCalls    int
}

// NewZonesRecommendationsService creates new ZonesRecommendationsServic.
// If creation failed the error will be non-nil.
// numConcurrentCalls should be positive, otherwise would return error.
func NewZonesRecommendationsService(ctx context.Context, numConcurrentCalls int) (*ZonesRecommendationsService, error) {
	compService, err := recommender.NewService(ctx)
	if err != nil {
		return nil, err
	}

	recommendationService := recommender.NewProjectsLocationsRecommendersRecommendationsService(compService)
	recService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	zonesService := compute.NewZonesService(recService)

	if numConcurrentCalls <= 0 {
		return nil, fmt.Errorf("Number of concurrent calls %d should be positive", numConcurrentCalls)
	}

	return &ZonesRecommendationsService{
		ctx:                   ctx,
		recommendationService: recommendationService,
		zonesService:          zonesService,
		numConcurrentCalls:    numConcurrentCalls,
	}, nil
}

func (s *ZonesRecommendationsService) listZoneRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	listCall := s.recommendationService.List(fmt.Sprintf("projects/%s/locations/%s/recommenders/%s", project, location, recommenderID))
	var recommendations []*recommender.GoogleCloudRecommenderV1Recommendation
	addRecommendations := func(response *recommender.GoogleCloudRecommenderV1ListRecommendationsResponse) error {
		recommendations = append(recommendations, response.Recommendations...)
		return nil
	}

	err := listCall.Pages(s.ctx, addRecommendations)
	if err != nil {
		log.Fatal(err)
	}
	return recommendations
}

func (s *ZonesRecommendationsService) listZonesNames(project string) []string {
	listCall := s.zonesService.List(project)

	var zones []string
	addZones := func(zoneList *compute.ZoneList) error {
		for _, zone := range zoneList.Items {
			zones = append(zones, zone.Name)
		}
		return nil
	}
	err := listCall.Pages(s.ctx, addZones)
	if err != nil {
		log.Fatal(err)
	}
	return zones
}

func (s *ZonesRecommendationsService) getNumConcurrentCalls() int {
	return s.numConcurrentCalls
}

// ListRecommendations returns the list of recommendations for a Cloud project.
// Requires the recommender.*.list IAM permission for the specified recommender.
func ListRecommendations(service Service, project string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	zones := service.listZonesNames(project)
	numberOfZones := len(zones)

	numWorkers := service.getNumConcurrentCalls()
	if numberOfZones < numWorkers {
		numWorkers = numberOfZones
	}
	results := make(chan []*recommender.GoogleCloudRecommenderV1Recommendation, numberOfZones)
	zonesJobs := make(chan string, numberOfZones)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for zone := range zonesJobs {
				results <- service.listZoneRecommendations(project, zone, recommenderID)
			}
		}()
	}

	for _, zone := range zones {
		zonesJobs <- zone
	}
	close(zonesJobs)

	var recommendations []*recommender.GoogleCloudRecommenderV1Recommendation
	for range zones {
		recommendations = append(recommendations, <-results...)
	}
	return recommendations
}
