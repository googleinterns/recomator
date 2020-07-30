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

	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

var requiredAPIs = []string{"compute.googleapis.com", "recommender.googleapis.com", "cloudresourcemanager.googleapis.com"}

// ListDisabledAPIs returns the list of not enabled apis for the specified project from apis list.
// Uses services.get method. If Service Usage API is not enabled returns it in the list of disabled APIs.
func (s *googleService) ListDisabledAPIs(project string, apis []string) ([]string, error) {
	servicesService := serviceusage.NewServicesService(s.serviceUsageService)
	var result []string
	for _, api := range apis {
		response, err := servicesService.Get("projects/" + project + "/services/" + api).Do()
		if err != nil {
			googleErr, ok := err.(*googleapi.Error)
			if ok && googleErr.Code == http.StatusForbidden {
				return []string{"serviceusage.googleapis.com"}, nil

			}
			return nil, err
		}
		if response.State != "ENABLED" {
			result = append(result, api)
		}
	}
	return result, nil
}

var requiredPermissions = [][]string{}

// TODO permissions check
