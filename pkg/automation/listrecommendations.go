package automation

import (
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"
)

// ListRecommendations returns the list of recommendations for specified project, zone, recommender.
// projects.locations.recommenders.recommendations/list method from Recommender API is used
func (s *googleService) ListRecommendations(project string, location string, recommenderID string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	listCall := s.recommenderService.List(fmt.Sprintf("projects/%s/locations/%s/recommenders/%s", project, location, recommenderID))
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

func (s *googleService) ListZonesNames(project string) []string {
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

// ListRecommendations returns the list of recommendations for a Cloud project.
// Requires the recommender.*.list IAM permission for the specified recommender.
// numConcurrentCalls specifies the maximum number of concurrent calls to ListZoneRecommendations,
// non-positive values are ignored, instead the default value is used.
func ListRecommendations(service GoogleService, project string, recommenderID string, numConcurrentCalls int) []*recommender.GoogleCloudRecommenderV1Recommendation {
	zones := service.ListZonesNames(project)
	numberOfZones := len(zones)

	numWorkers := numConcurrentCalls
	const defaultNumWorkers = 16
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	results := make(chan []*recommender.GoogleCloudRecommenderV1Recommendation, numberOfZones)
	zonesJobs := make(chan string, numberOfZones)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for zone := range zonesJobs {
				results <- service.ListRecommendations(project, zone, recommenderID)
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
