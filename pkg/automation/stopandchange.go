package automation

import (
	"fmt"

	"google.golang.org/api/compute/v1"
)

func (s *googleService) ChangeMachineType(project string, zone string, instance string, machineType string) error {
	machineType = fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
	request := &compute.InstancesSetMachineTypeRequest{MachineType: machineType}
	_, err := s.instancesService.SetMachineType(project, zone, instance, request).Do()
	return err
}

func (s *googleService) StopInstance(project string, zone string, instance string) error {
	_, err := s.instancesService.Stop(project, zone, instance).Do()
	return err
}
