package automation

import (
	"context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"
	"log"
)

func listZoneRecommendations(project string, location string, recommenderId string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	ctx := context.Background()
	service, err := recommender.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	recService := recommender.NewProjectsLocationsRecommendersRecommendationsService(service)

	listCall := recService.List(fmt.Sprintf("projects/%s/locations/%s/recommenders/%s", project, location, recommenderId))
	var recommendations []*recommender.GoogleCloudRecommenderV1Recommendation
	addRecommendations := func(response *recommender.GoogleCloudRecommenderV1ListRecommendationsResponse) error {
		recommendations = append(recommendations, response.Recommendations...)
		return nil
	}

	ctx = context.Background()
	err = listCall.Pages(ctx, addRecommendations)
	if err != nil {
		log.Fatal(err)
	}
	return recommendations
}

func listZonesNames(project string) []string {
	ctx := context.Background()
	service, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	zonesService := compute.NewZonesService(service)
	listCall := zonesService.List(project)

	var zones []string
	addZones := func(zoneList *compute.ZoneList) error {
		for _, zone := range zoneList.Items {
			zones = append(zones, zone.Name)
		}
		return nil
	}
	ctx = context.Background()
	err = listCall.Pages(ctx, addZones)
	if err != nil {
		log.Fatal(err)
	}
	return zones
}

func ListRecommendations(project string, recommenderId string) []*recommender.GoogleCloudRecommenderV1Recommendation {
	zones := listZonesNames(project)
	ch := make(chan []*recommender.GoogleCloudRecommenderV1Recommendation, len(zones))
	for _, zone := range zones {
		go func(zoneName string) {
			ch <- listZoneRecommendations(project, zoneName, recommenderId)
		}(zone)
	}
	var recommendations []*recommender.GoogleCloudRecommenderV1Recommendation
	for _ = range zones {
		recommendations = append(recommendations, <-ch...)
	}
	return recommendations
}
