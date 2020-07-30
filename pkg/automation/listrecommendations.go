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
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"
)

// gcloudRecommendation is a type alias for Google Cloud Recommendation recommender.GoogleCloudRecommenderV1Recommendation
type gcloudRecommendation = recommender.GoogleCloudRecommenderV1Recommendation

// ListRecommendations returns the list of recommendations for specified project, zone, recommender.
// projects.locations.recommenders.recommendations/list method from Recommender API is used.
// If the error occurred the returned error is not nil.
func (s *googleService) ListRecommendations(project, location, recommenderID string) ([]*gcloudRecommendation, error) {
	recommendationsService := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	listCall := recommendationsService.List(fmt.Sprintf("projects/%s/locations/%s/recommenders/%s", project, location, recommenderID))
	var recommendations []*gcloudRecommendation
	addRecommendations := func(response *recommender.GoogleCloudRecommenderV1ListRecommendationsResponse) error {
		recommendations = append(recommendations, response.Recommendations...)
		return nil
	}

	err := listCall.Pages(s.ctx, addRecommendations)
	if err != nil {
		return []*gcloudRecommendation{}, err
	}
	return recommendations, nil
}

// ListZonesNames returns list of zone names for the specified project.
// Uses zones/list method from Compute API.
// If the error occurred the returned error is not nil.
func (s *googleService) ListZonesNames(project string) ([]string, error) {
	zonesService := compute.NewZonesService(s.computeService)
	listCall := zonesService.List(project)

	var zones []string
	addZones := func(zoneList *compute.ZoneList) error {
		for _, zone := range zoneList.Items {
			zones = append(zones, zone.Name)
		}
		return nil
	}
	err := listCall.Pages(s.ctx, addZones)
	if err != nil {
		return []string{}, err
	}
	return zones, nil
}

// ListRegionsNames returns list of region names for the specified project.
// Uses region/list method from Compute API.
// If the error occurred the returned error is not nil.
func (s *googleService) ListRegionNames(project string) ([]string, error) {
	regionsService := compute.NewRegionsService(s.computeService)
	listCall := regionsService.List(project)

	var regions []string
	addRegions := func(regionList *compute.RegionList) error {
		for _, region := range regionList.Items {
			regions = append(regions, region.Name)
		}
		return nil
	}
	err := listCall.Pages(s.ctx, addRegions)
	if err != nil {
		return []string{}, err
	}
	return regions, nil
}

type result struct {
	recommendations []*gcloudRecommendation
	err             error
}

// ListRecommendations returns the list of recommendations for a Cloud project.
// Requires the recommender.*.list IAM permission for the specified recommender.
// numConcurrentCalls specifies the maximum number of concurrent calls to ListRecommendations method,
// non-positive values are ignored, instead the default value is used.
func ListRecommendations(service GoogleService, project, recommenderID string, numConcurrentCalls int) ([]*gcloudRecommendation, error) {
	zones, err := service.ListZonesNames(project)
	if err != nil {
		return []*gcloudRecommendation{}, err
	}
	numberOfZones := len(zones)

	regions, err := service.ListRegionsNames(project)
	if err != nil {
		return []*gcloudRecommendation{}, err
	}
	numberOfRegions := len(regions)

	locations := append(zones, regions...)
	numberOfLocations := len(locations);

	numWorkers := numConcurrentCalls
	const defaultNumWorkers = 16
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	results := make(chan result, numberOfLocations)
	locationsJobs := make(chan string, numberOfLocations)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for location := range locationsJobs {
				recs, err := service.ListRecommendations(project, location, recommenderID)
				results <- result{recs, err}
			}
		}()
	}

	for _, location := range locations {
		locationsJobs <- location
	}

	close(locationsJobs)

	var recommendations []*gcloudRecommendation
	err = nil
	for range locations {
		locationResult := <-results
		if locationResult.err != nil {
			err = locationResult.err
		} else {
			recommendations = append(recommendations, locationResult.recommendations...)
		}
	}
	if err != nil {
		return []*gcloudRecommendation{}, err
	}
	return recommendations, nil
}
