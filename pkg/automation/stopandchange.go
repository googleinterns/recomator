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
)

func (s *googleService) ChangeMachineType(project string, zone string, instance string, machineType string) error {
	machineType = fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
	request := &compute.InstancesSetMachineTypeRequest{MachineType: machineType}
	instancesService := compute.NewInstancesService(s.computeService)
	_, err := instancesService.SetMachineType(project, zone, instance, request).Do()
	return err
}

func (s *googleService) StopInstance(project string, zone string, instance string) error {
	instancesService := compute.NewInstancesService(s.computeService)
	_, err := instancesService.Stop(project, zone, instance).Do()
	return err
}
