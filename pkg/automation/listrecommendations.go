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
		return nil, err
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
		return nil, err
	}
	return zones, nil
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
	err := listCall.Pages(s.ctx, addRegions)
	if err != nil {
		return []string{}, err
	}
	return regions, nil
}

// ListRecommendations returns the list of recommendations for a Cloud project.
// Requires the recommender.*.list IAM permission for the specified recommender.
// numConcurrentCalls specifies the maximum number of concurrent calls to ListRecommendations method,
// non-positive values are ignored, instead the default value is used.
func ListRecommendations(service GoogleService, project, recommenderID string, numConcurrentCalls int, tasks ...*Task) ([]*gcloudRecommendation, error) {
	var task *Task
	if len(tasks) != 0 {
		task = tasks[0]
	}

	zones, err := service.ListZonesNames(project)
	if err != nil {
		return nil, err
	}

	regions, err := service.ListRegionsNames(project)
	if err != nil {
		return []*gcloudRecommendation{}, err
	}

	locations := append(zones, regions...)
	numberOfLocations := len(locations)

	task.AddSubtasks(numberOfLocations)

	numWorkers := numConcurrentCalls
	const defaultNumWorkers = 16
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	type result struct {
		recommendations []*gcloudRecommendation
		err             error
	}

	results := make(chan result, numberOfLocations)
  	locationsJobs := make(chan string, numberOfLocations)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for location := range locationsJobs {
				recs, err := service.ListRecommendations(project, location, recommenderID)
				results <- result{recs, err}
				task.IncrementDone()
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
		return nil, err
	}
	defer task.SetAllDone()
	return recommendations, nil
}

var googleRecommenders = []string{
	"google.compute.disk.IdleResourceRecommender",
	"google.compute.instance.IdleResourceRecommender",
	"google.compute.instance.MachineTypeRecommender",
}

// ListAllRecommendersRecommendations lists recommendations for all googleRecommenders
func ListAllRecommendersRecommendations(service GoogleService, project string, numConcurrentCalls int, tasks ...*Task) ([]*gcloudRecommendation, error) {
	var task *Task
	if len(tasks) != 0 {
		task = tasks[0]
	}
	task.AddSubtasks(len(googleRecommenders))
	var recommendations []*gcloudRecommendation
	for _, recommender := range googleRecommenders {

		newRecommendations, err := ListRecommendations(service, project, recommender, numConcurrentCalls, task.GetNextSubtask())
		if err != nil {
			return nil, err
		}
		recommendations = append(recommendations, newRecommendations...)

		task.IncrementDone()
	}
	defer task.SetAllDone()
	return recommendations, nil
}

// ListResult contains information about listing recommendations for all projects.
// If user doesn't have enough permissions for the project, the requirements, including failed ones, are listed in failedProjects.
// Otherwise, recommendations for the project are appended to recommendations.
type ListResult struct {
	recommendations []*gcloudRecommendation
	failedProjects  []*ProjectRequirements
}

func listRecommendationsIfRequirementsCompleted(service GoogleService, projectsRequirements []*ProjectRequirements, numConcurrentCalls int, tasks ...*Task) (*ListResult, error) {
	var task *Task
	if len(tasks) != 0 {
		task = tasks[0]
	}
	task.AddSubtasks(len(projectsRequirements))

	var listResult ListResult
	for _, projectRequirements := range projectsRequirements {
		ok := true
		for _, req := range projectRequirements.Requirements {
			if req.Status == RequirementFailed {
				ok = false
				break
			}
		}
		if ok {
			newRecs, err := ListAllRecommendersRecommendations(service, projectRequirements.Project, numConcurrentCalls, task.GetNextSubtask())
			if err != nil {
				return nil, err
			}
			listResult.recommendations = append(listResult.recommendations, newRecs...)
		} else {
			listResult.failedProjects = append(listResult.failedProjects, projectRequirements)
		}

		task.IncrementDone()
	}

	defer task.SetAllDone()
	return &listResult, nil
}

// ListAllProjectsRecommendations gets all projects for which user has projects.get permission.
// If the user has enough permissions to apply and list recommendations, recommendations for projects are listed.
// Otherwise, projects requirements, including failed ones, are added to `failedProjects` to help show warnings to the user.
func ListAllProjectsRecommendations(service GoogleService, numConcurrentCalls int, tasks ...*Task) (*ListResult, error) {
	var task *Task
	if len(tasks) != 0 {
		task = tasks[0]
	}
	projects, err := service.ListProjects()
	if err != nil {
		return nil, err
	}

	task.AddSubtasks(2)

	projectsRequirements, err := ListRequirements(service, projects, task.GetNextSubtask())
	if err != nil {
		return nil, err
	}

	task.IncrementDone()

	listResult, err := listRecommendationsIfRequirementsCompleted(service, projectsRequirements, numConcurrentCalls, task.GetNextSubtask())

	if err != nil {
		return nil, err
	}
	task.IncrementDone()

	defer task.SetAllDone()
	return listResult, nil
}
