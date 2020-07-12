package automation

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)


// RecommendationApplyService is used to apply recommendations for VMs
type RecommendationApplyService struct {
	compute.InstancesService instansesService;
}

// NewRecommendationApplyService creates new RecommendationApplyService
func NewRecommendationApplyService(ctx context.Context) (*RecommendationApplyService, error) {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return RecommendationApplyService {
		compute.NewInstancesService(computeService)
	}, nil
}

func (s *RecommendationApplyService) stopInstance(project string, zone string, instance string) {
	_, err = s.instancesService.Stop(project, zone, instance).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *RecommendationApplyService) changeMachineType(project string, zone string, instance string, machineType string) {
	machineType = fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
	request := &compute.InstancesSetMachineTypeRequest{MachineType: machineType}
	_, err = s.instancesService.SetMachineType(project, zone, instance, request).Do()
	if err != nil {
		log.Fatal(err)
	}
}
