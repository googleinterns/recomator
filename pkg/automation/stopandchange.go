package automation

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

func stopInstance(project string, zone string, instance string) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	instancesService := compute.NewInstancesService(computeService)
	_, err = instancesService.Stop(project, zone, instance).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func changeMachineType(project string, zone string, instance string, machineType string) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	instancesService := compute.NewInstancesService(computeService)
	machineType = fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
	request := &compute.InstancesSetMachineTypeRequest{MachineType: machineType}
	_, err = instancesService.SetMachineType(project, zone, instance, request).Do()
	if err != nil {
		log.Fatal(err)
	}
}
