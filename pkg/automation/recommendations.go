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

// GetRecommendation implements projects.locations.recommenders.recommendations/get method
func (s *googleService) GetRecommendation(name string) (*gcloudRecommendation, error) {
	service := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	var recommendation *gcloudRecommendation
	err := DoRequestWithRetries(func() error {
		rec, err := service.Get(name).Do()
		recommendation = rec
		return err
	})
	return recommendation, err
}

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
	err := DoRequestWithRetries(func() error {
		recommendations = nil
		return listCall.Pages(s.ctx, addRecommendations)
	})
	return recommendations, err
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
	err := DoRequestWithRetries(func() error {
		zones = nil
		return listCall.Pages(s.ctx, addZones)
	})
	return zones, err
}

// ListRegionsNames returns list of region names for the specified project.
// Uses regions/list method from Compute API.
// If the error occurred the returned error is not nil.
func (s *googleService) ListRegionsNames(project string) ([]string, error) {
	regionsService := compute.NewRegionsService(s.computeService)
	listCall := regionsService.List(project)

	var regions []string
	addRegions := func(regionList *compute.RegionList) error {
		for _, region := range regionList.Items {
			regions = append(regions, region.Name)
		}
		return nil
	}
	err := DoRequestWithRetries(func() error {
		regions = nil
		return listCall.Pages(s.ctx, addRegions)
	})
	return regions, err
}

// ListLocations return the list of all locations per project(zones and regions).
// Exactly one of returned values will be non-nil.
func ListLocations(service GoogleService, project string) ([]string, error) {
	zones, err := service.ListZonesNames(project)
	if err != nil {
		return nil, err
	}

	regions, err := service.ListRegionsNames(project)
	if err != nil {
		return nil, err
	}

	locations := append(zones, regions...)
	return locations, nil
}

var googleRecommenders = []string{
	"google.compute.disk.IdleResourceRecommender",
	"google.compute.instance.IdleResourceRecommender",
	"google.compute.instance.MachineTypeRecommender",
}

type recommendationsResult struct {
	recommendations []*gcloudRecommendation
	err             error
}

// concatResults receives numberOfResults values from results channel.
// Returns concatenated slice of all recommendations. If one of results contains error, returns error.
// At most one of returned values will be non-nil.
func concatResults(results <-chan recommendationsResult, numberOfResults int) ([]*gcloudRecommendation, error) {
	var err error
	var recommendations []*gcloudRecommendation
	for i := 0; i < numberOfResults; i++ {
		result := <-results
		if result.err != nil {
			err = result.err
		} else {
			recommendations = append(recommendations, result.recommendations...)
		}
	}

	if err != nil {
		return nil, err
	}
	return recommendations, nil
}

// ListRecommendations returns the list of recommendations for a Cloud project from googleRecommenders.
// Requires the recommender.*.list IAM permissions for the recommenders.
// numConcurrentCalls specifies the maximum number of concurrent calls to ListRecommendations method,
// non-positive values are ignored, instead the default value is used.
// task structure tracks the progress of the function.
func ListRecommendations(service GoogleService, project string, numConcurrentCalls int, task *Task) ([]*gcloudRecommendation, error) {
	locations, err := ListLocations(service, project)
	if err != nil {
		return nil, err
	}

	numWorkers := numConcurrentCalls
	const defaultNumWorkers = 16
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	type query struct {
		location      string
		recommenderID string
	}

	numberOfQueries := len(locations) * len(googleRecommenders)
	task.SetNumberOfSubtasks(numberOfQueries)

	results := make(chan recommendationsResult, numberOfQueries)
	queries := make(chan query, numberOfQueries)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for query := range queries {
				recs, err := service.ListRecommendations(project, query.location, query.recommenderID)
				results <- recommendationsResult{recs, err}
				task.IncrementDone()
			}
		}()
	}

	for _, recommenderID := range googleRecommenders {
		for _, location := range locations {
			queries <- query{location: location, recommenderID: recommenderID}
		}
	}

	close(queries)

	recommendations, err := concatResults(results, numberOfQueries)
	if err == nil {
		task.SetAllDone()
	}
	return recommendations, err
}

// ListResult contains information about listing recommendations for all projects.
// If user doesn't have enough permissions for the project, the requirements, including failed ones, are listed in failedProjects.
// Otherwise, recommendations for the project are appended to recommendations.
type ListResult struct {
	Recommendations []*gcloudRecommendation
	FailedProjects  []*ProjectRequirements
}

// Lists requirements for the project, if all satisfied - lists recommendations.
// Adds results to listResult.
// Otherwise, adds project's requirements in FailedProjects field.
func listRecommendationsIfRequirementsSatisfied(service GoogleService, project string, numConcurrentCalls int, listResult *ListResult, task *Task) error {
	task.SetNumberOfSubtasks(2) // CheckRequirements and ListRecommendations

	task.GetNextSubtask()
	projectRequirements, err := ListProjectRequirements(service, project)

	if err != nil {
		return err
	}

	task.IncrementDone()

	for _, req := range projectRequirements {
		if !req.Satisfied {
			listResult.FailedProjects = append(listResult.FailedProjects,
				&ProjectRequirements{Project: project, Requirements: projectRequirements})
			task.SetAllDone()
			return nil
		}
	}
	newRecs, err := ListRecommendations(service, project, numConcurrentCalls, task.GetNextSubtask())
	if err != nil {
		return err
	}
	task.IncrementDone()
	listResult.Recommendations = append(listResult.Recommendations, newRecs...)
	return nil
}

// ListProjectsRecommendations gets recommendations for the specified projects.
// If the user has enough permissions to apply and list recommendations, recommendations for project are listed.
// Otherwise, projects requirements, including failed ones, are added to `failedProjects` to help show warnings to the user.
// task structure tracks how many subtasks have been done already.
func ListProjectsRecommendations(service GoogleService, projects []string, numConcurrentCalls int, task *Task) (*ListResult, error) {
	task.SetNumberOfSubtasks(len(projects)) // subtasks are calls to listRecommendationsIfRequirementsSatisfied for each project

	var listResult ListResult
	for _, project := range projects {
		err := listRecommendationsIfRequirementsSatisfied(service, project, numConcurrentCalls, &listResult, task.GetNextSubtask())
		if err != nil {
			return nil, err
		}
		task.IncrementDone()
	}
	task.SetAllDone()
	return &listResult, nil
}
