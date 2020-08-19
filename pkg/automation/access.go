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
	"net/http"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

// requiredAPIs are APIs required for googleService
var requiredAPIs = []string{"compute.googleapis.com", "recommender.googleapis.com", "cloudresourcemanager.googleapis.com"}

const (
	// RequirementFailed corresponds to the failed status of requirement
	RequirementFailed = "FAILED"
	// RequirementCompleted corresponds to the completed status of requirement
	RequirementCompleted = "COMPLETED"
)

// Requirement contains information about the required permission or api.
// Status states whether permission is given and api is enabled.
// If Status if FAILED, ErrorMessage will contain information about the problem.
type Requirement struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

// ListAPIRequirements returns the list of APIs and their statuses for the specified project.
// services.get permission is required for this method.
// First requirement in returned list will be related to Service Usage API.
// If it is not enabled or user doesn't have services.get permission other APIs won't be checked.
func (s *googleService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	servicesService := serviceusage.NewServicesService(s.serviceUsageService)
	serviceUsageAPI := "serviceusage.googleapis.com"
	serviceUsageName := "Service Usage API and services.get permission"
	_, err := servicesService.Get("projects/" + project + "/services/" + serviceUsageAPI).Do()
	if err != nil {
		googleErr := err.(*googleapi.Error)
		if googleErr.Code == http.StatusForbidden {
			return []*Requirement{&Requirement{
				Name:         serviceUsageName,
				Status:       RequirementFailed,
				ErrorMessage: googleErr.Message,
			}}, nil
		}
		return nil, err
	}
	result := []*Requirement{&Requirement{
		Name:   serviceUsageName,
		Status: RequirementCompleted,
	}}
	for _, api := range apis {
		response, err := servicesService.Get("projects/" + project + "/services/" + api).Do()
		if err != nil {
			return nil, err
		}
		status := RequirementCompleted
		errorMessage := ""
		if response.State != "ENABLED" {
			status = RequirementFailed
			errorMessage = response.Config.Name + " has not been used in project " + project +
				" before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/" + api +
				"/overview?project=" + project + " then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry."
		}

		requirement := &Requirement{
			Name:         response.Config.Title,
			Status:       status,
			ErrorMessage: errorMessage,
		}
		result = append(result, requirement)
	}
	return result, nil
}

// requiredPermissions are permissions required for googleService
var requiredPermissions = [][]string{
	[]string{"compute.instances.setMachineType"},                              // ChangeMachineType
	[]string{"compute.disks.createSnapshot", "compute.snapshots.create"},      // CreateSnapshot
	[]string{"compute.disks.delete"},                                          // DeleteDisk
	[]string{"compute.instances.get"},                                         // GetInstance
	[]string{"recommender.computeDiskIdleResourceRecommendations.list"},       // ListRecommendations for google.compute.disk.IdleResourceRecommender
	[]string{"recommender.computeInstanceIdleResourceRecommendations.list"},   // ListRecommendations for google.compute.instance.IdleResourceRecommender
	[]string{"recommender.computeInstanceMachineTypeRecommendations.list"},    // ListRecommendations for google.compute.instance.MachineTypeRecommender
	[]string{"recommender.computeDiskIdleResourceRecommendations.get"},        // GetRecommendation for google.compute.disk.IdleResourceRecommender
	[]string{"recommender.computeInstanceIdleResourceRecommendations.get"},    // GetRecommendation for google.compute.instance.IdleResourceRecommender
	[]string{"recommender.computeInstanceMachineTypeRecommendations.get"},     // GetRecommendation for google.compute.instance.MachineTypeRecommender
	[]string{"recommender.computeDiskIdleResourceRecommendations.update"},     // MarkClaimed/Failed/Suceeded for google.compute.disk.IdleResourceRecommender
	[]string{"recommender.computeInstanceIdleResourceRecommendations.update"}, // MarkClaimed/Failed/Suceeded for google.compute.instance.IdleResourceRecommender
	[]string{"recommender.computeInstanceMachineTypeRecommendations.update"},  // ListRecommendations for google.compute.instance.MachineTypeRecommender
	[]string{"compute.regions.list"},                                          // ListRegionsNames
	[]string{"compute.zones.list"},                                            // ListZonesNames
	[]string{"compute.instances.stop"},                                        // StopInstance
}

// ListPermissionRequirements returns the list of permissions and their statuses for the project.
// No permissions required for this method.
func (s *googleService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	var result []*Requirement
	for _, permissionsGroup := range permissions {
		status := RequirementCompleted
		errorMessage := ""
		request := cloudresourcemanager.TestIamPermissionsRequest{Permissions: permissionsGroup}
		projectsService := cloudresourcemanager.NewProjectsService(s.resourceManagerService)
		response, err := projectsService.TestIamPermissions(project, &request).Do()
		if err != nil {
			return nil, err
		}

		if len(response.Permissions) == 0 {
			status = RequirementFailed
			errorMessage = "At least one of these permissions is needed. None found."
		}
		name := strings.Join(permissionsGroup, ", ")
		result = append(result, &Requirement{Name: name, Status: status, ErrorMessage: errorMessage})
	}
	return result, nil
}

// ProjectRequirements contains information about permissions for the user for the project.
type ProjectRequirements struct {
	Project      string         `json:"project"`
	Requirements []*Requirement `json:"requirements"`
}

// ListProjectRequirements is a function that lists all permissions and APIs and their statuses for a project.
// If all statuses are equal to RequirementCompleted, user has all required permissions.
func ListProjectRequirements(s GoogleService, project string) ([]*Requirement, error) {
	requirements, err := s.ListAPIRequirements(project, requiredAPIs)
	if err != nil {
		return nil, err
	}
	for _, req := range requirements {
		if req.Status == RequirementFailed {
			return requirements, nil
		}
	}

	permissions, err := s.ListPermissionRequirements(project, requiredPermissions)
	if err != nil {
		return nil, err
	}
	requirements = append(requirements, permissions...)
	return requirements, nil
}

// ListRequirements lists the requirements and their statuses for every project.
// task structure tracks how many projects have been processed already.
func ListRequirements(s GoogleService, projects []string, task *Task) ([]*ProjectRequirements, error) {
	task.SetNumberOfSubtasks(len(projects))
	var result []*ProjectRequirements
	for _, project := range projects {
		requirements, err := ListProjectRequirements(s, project)
		if err != nil {
			return nil, err
		}
		result = append(result, &ProjectRequirements{project, requirements})
		task.IncrementDone()
	}
	task.SetAllDone()
	return result, nil
}
