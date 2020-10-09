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
	"sort"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

// requiredAPIs are APIs required for googleService
var requiredAPIs = []string{"compute.googleapis.com", "recommender.googleapis.com", "cloudresourcemanager.googleapis.com"}

// Requirement contains information about the required permission or api.
// Satisfied states whether permission is given and API is enabled.
// If Satisfied if false, ErrorMessage will contain information about the problem.
type Requirement struct {
	Name         string `json:"name"`
	Satisfied    bool   `json:"satisfied"`
	ErrorMessage string `json:"errorMessage"`
}

// ListAPIRequirements returns the list of APIs and their statuses for the specified project.
// services.get permission is required for this method.
// First requirement in returned list will be related to Service Usage API.
// If it is not enabled other APIs won't be checked.
func (s *googleService) ListAPIRequirements(project string, apis []string) ([]*Requirement, error) {
	servicesService := serviceusage.NewServicesService(s.serviceUsageService)
	serviceUsageName := "Service Usage API"
	result := []*Requirement{}
	for _, api := range apis {
		var response *serviceusage.GoogleApiServiceusageV1Service
		err := DoRequestWithRetries(func() error {
			resp, err := servicesService.Get("projects/" + project + "/services/" + api).Do()
			response = resp
			return err
		})
		if len(result) == 0 { // If it's the first request we also check implicitly if Service Usage API is enabled.
			if err != nil { // If there is an error because Service Usage API is not enabled, we return failed requirement.
				googleErr, ok := err.(*googleapi.Error)
				if ok && googleErr.Code == http.StatusForbidden {
					return []*Requirement{&Requirement{
						Name:         serviceUsageName,
						Satisfied:    false,
						ErrorMessage: googleErr.Message,
					}}, nil
				}
			} else {
				result = []*Requirement{&Requirement{
					Name:      serviceUsageName,
					Satisfied: true,
				}}
			}
		}
		if err != nil {
			return nil, err
		}
		satisfied := true
		errorMessage := ""
		if response.State != "ENABLED" {
			satisfied = false
			errorMessage = response.Config.Name + " has not been used in project " + project +
				" before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/" + api +
				"/overview?project=" + project + " then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry."
		}

		requirement := &Requirement{
			Name:         response.Config.Title,
			Satisfied:    satisfied,
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
	[]string{"compute.instances.start"},                                       // StartInstance
	[]string{"compute.instances.stop"},                                        // StopInstance
	[]string{"serviceusage.services.get"},                                     // ListAPIRequirements
}

// ListPermissionRequirements returns the list of permissions and their statuses for the project.
// No permissions required for this method.
// If cloud resource manager api is not enabled, will return not satisfied requirement for this API.
func (s *googleService) ListPermissionRequirements(project string, permissions [][]string) ([]*Requirement, error) {
	var result []*Requirement
	var allPermissions []string

	for _, permissionsGroup := range permissions {
		allPermissions = append(allPermissions, permissionsGroup...)
	}

	request := cloudresourcemanager.TestIamPermissionsRequest{Permissions: allPermissions}
	projectsService := cloudresourcemanager.NewProjectsService(s.resourceManagerService)
	var response *cloudresourcemanager.TestIamPermissionsResponse
	err := DoRequestWithRetries(func() error {
		resp, err := projectsService.TestIamPermissions(project, &request).Do()
		response = resp
		return err
	})
	if err != nil {
		googleErr, ok := err.(*googleapi.Error)
		if ok && googleErr.Code == http.StatusForbidden {
			return []*Requirement{&Requirement{
				Name:         "Cloud Resource Manager API",
				Satisfied:    false,
				ErrorMessage: googleErr.Message,
			}}, nil
		}
		return nil, err
	}

	satisfiedPermissions := response.Permissions
	sort.Strings(satisfiedPermissions)

	for _, permissionsGroup := range permissions {
		satisfied := false
		errorMessage := ""

		for _, permission := range permissionsGroup {
			ind := sort.SearchStrings(satisfiedPermissions, permission)
			// if one permission in group is satisfied then we are done
			if ind < len(satisfiedPermissions) && satisfiedPermissions[ind] == permission {
				satisfied = true
				break
			}
		}

		if !satisfied {
			errorMessage = "At least one of these permissions is needed. None found."
		}
		name := strings.Join(permissionsGroup, " or ")
		result = append(result, &Requirement{Name: name, Satisfied: satisfied, ErrorMessage: errorMessage})
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
	requirements, err := s.ListPermissionRequirements(project, requiredPermissions)
	if err != nil {
		return nil, err
	}

	for _, req := range requirements {
		if !req.Satisfied {
			return requirements, nil
		}
	}

	apiRequirements, err := s.ListAPIRequirements(project, requiredAPIs)
	if err != nil {
		return nil, err
	}
	requirements = append(requirements, apiRequirements...)
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
