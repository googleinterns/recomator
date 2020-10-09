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

	"github.com/google/uuid"
	"google.golang.org/api/compute/v1"
)

// ChangeMachineType changes machine type using instances.setMachineType method.
// Requires compute.instances.setMachineType permission.
func (s *googleService) ChangeMachineType(project string, zone string, instance string, machineType string) error {
	machineType = fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
	request := &compute.InstancesSetMachineTypeRequest{MachineType: machineType}
	instancesService := compute.NewInstancesService(s.computeService)

	requestID := uuid.New().String()
	return DoRequestWithRetries(func() error {
		return AwaitCompletion(func() (*compute.Operation, error) {
			return instancesService.SetMachineType(project, zone, instance, request).RequestId(requestID).Do()
		}, sleepTimeChangingMachineType)
	})
}

// GetInstance gets instance using instances.get method.
// Requires compute.instances.get permission.
func (s *googleService) GetInstance(project string, zone string, instance string) (*compute.Instance, error) {
	instancesService := compute.NewInstancesService(s.computeService)
	var instanceVal *compute.Instance
	err := DoRequestWithRetries(func() error {
		inst, err := instancesService.Get(project, zone, instance).Do()
		instanceVal = inst
		return err
	})
	return instanceVal, err
}

// StopInstance stops instance using instances.stop method.
// Requires compute.instances.stop permission.
func (s *googleService) StopInstance(project string, zone string, instance string) error {
	instancesService := compute.NewInstancesService(s.computeService)
	requestID := uuid.New().String()
	return DoRequestWithRetries(func() error {
		return AwaitCompletion(func() (*compute.Operation, error) {
			return instancesService.Stop(project, zone, instance).RequestId(requestID).Do()
		}, sleepTimeStoppingInstance)
	})
}

// StartInstance starts instance using instances.start method.
// Requires compute.instances.start permission.
func (s *googleService) StartInstance(project string, zone string, instance string) error {
	instancesService := compute.NewInstancesService(s.computeService)
	requestID := uuid.New().String()
	return DoRequestWithRetries(func() error {
		return AwaitCompletion(func() (*compute.Operation, error) {
			return instancesService.Start(project, zone, instance).RequestId(requestID).Do()
		}, sleepTimeStartingInstance)
	})
}
