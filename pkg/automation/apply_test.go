package automation

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/recommender/v1"

	"github.com/stretchr/testify/assert"
)

type gcloudStateInfo = recommender.GoogleCloudRecommenderV1RecommendationStateInfo
type gcloudContent = recommender.GoogleCloudRecommenderV1RecommendationContent
type calledFunction struct {
	functionName string
	arguments    []interface{}
	results      []interface{}
}

type ApplyMockService struct {
	GoogleService
	calledFunctions   []calledFunction
	getInstanceResult *compute.Instance
}

func (s *ApplyMockService) GetInstance(project string, zone string, instance string) (*compute.Instance, error) {
	newCalledFunction := calledFunction{"GetInstance", []interface{}{project, zone, instance}, []interface{}{s.getInstanceResult, nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return s.getInstanceResult, nil
}

func (s *ApplyMockService) StopInstance(project string, zone string, instance string) error {
	newCalledFunction := calledFunction{"StopInstance", []interface{}{project, zone, instance}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *ApplyMockService) ChangeMachineType(project string, zone string, instance string, machineType string) error {
	newCalledFunction := calledFunction{"ChangeMachineType", []interface{}{project, zone, instance, machineType}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *ApplyMockService) StartInstance(project string, zone string, instance string) error {
	newCalledFunction := calledFunction{"StartInstance", []interface{}{project, zone, instance}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *ApplyMockService) CreateSnapshot(project string, zone string, disk string, name string) error {
	// it is not possible to say what the name should be equal to
	newCalledFunction := calledFunction{"CreateSnapshot", []interface{}{project, zone, disk, ""}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *ApplyMockService) DeleteDisk(project string, zone string, disk string) error {
	newCalledFunction := calledFunction{"DeleteDisk", []interface{}{project, zone, disk}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *ApplyMockService) MarkRecommendationClaimed(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationClaimed", []interface{}{name, etag}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil, nil
}

func (s *ApplyMockService) MarkRecommendationSucceeded(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationSucceeded", []interface{}{name, etag}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil, nil
}

func (s *ApplyMockService) MarkRecommendationFailed(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationFailed", []interface{}{name, etag}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil, nil
}

// Creates an array of type calledFunction, given arrays of functions,
// their arguments and results
func newCalledFunctions(functions []string, arguments [][]interface{}, results [][]interface{}) ([]calledFunction, error) {
	result := []calledFunction{}

	if len(functions) != len(arguments) || len(arguments) != len(results) {
		return nil, errors.New("lengths of the arguments must be equal")
	}

	for i := range functions {
		result = append(result, calledFunction{functions[i], arguments[i], results[i]})
	}

	return result, nil
}

// Checks if two arrays of calledFunction type are equal
func compareCalledFunctions(t *testing.T, expected, received []calledFunction) {
	assert.Equal(t, len(expected), len(received), "wrong number of functions were called")
	for i := range received {
		assert.Equal(t, expected[i].functionName, received[i].functionName, "a wrong function was called")

		assert.Equal(t, len(expected[i].arguments), len(received[i].arguments), fmt.Sprintf("function %s was called with a wrong number of arguments", expected[i].functionName))
		for j := range received[i].arguments {
			assert.Equal(t, expected[i].arguments[j], received[i].arguments[j], fmt.Sprintf("function %s was called with a wrong argument", expected[i].functionName))
		}

		assert.Equal(t, len(expected[i].results), len(received[i].results), fmt.Sprintf("function %s returned wrong number of values", expected[i].functionName))
		for j := range received[i].results {
			assert.Equal(t, expected[i].results[j], received[i].results[j], fmt.Sprintf("function %s returned a wrong result", expected[i].functionName))
		}
	}
}

// Checks if the test machine type operation works as expected.
func TestTestMachineTypeOperation(t *testing.T) {
	operation := gcloudOperation{
		Action:       "test",
		Path:         "/machineType",
		Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
		ResourceType: "compute.googleapis.com/Instance",
		ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/n1-standard-4"},
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{MachineType: "zones/us-east1-b/machineTypes/n1-standard-4"}}
	err := DoOperation(&service, &operation)
	assert.NoErrorf(t, err, "DoOperation shouldn't return an error")

	expectedFunctions := []string{"GetInstance"}
	expectedArguments := [][]interface{}{{"rightsizer-test", "us-east1-b", "alicja-test"}}
	expectedResults := [][]interface{}{{&compute.Instance{MachineType: "zones/us-east1-b/machineTypes/n1-standard-4"}, nil}}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	assert.Equal(t, expected, service.calledFunctions)
}

// Checks if the test status operation works as expected.
func TestTestStatusOperation(t *testing.T) {
	operation := gcloudOperation{
		Action:       "test",
		Path:         "/status",
		Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
		ResourceType: "compute.googleapis.com/Instance",
		Value:        "RUNNING",
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{Status: "RUNNING"}}
	err := DoOperation(&service, &operation)
	assert.NoErrorf(t, err, "DoOperation shouldn't return an error")

	expectedFunctions := []string{"GetInstance"}
	expectedArguments := [][]interface{}{{"rightsizer-test", "us-central1-a", "vkovalova-instance-memory-1"}}
	expectedResults := [][]interface{}{{&compute.Instance{Status: "RUNNING"}, nil}}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if the replace machine type operation works as expected.
func TestReplaceMachineTypeOperation(t *testing.T) {
	operation := gcloudOperation{
		Action:       "replace",
		Path:         "/machineType",
		Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
		ResourceType: "compute.googleapis.com/Instance",
		Value:        "zones/us-east1-b/machineTypes/custom-2-5120",
	}

	service := ApplyMockService{}
	err := DoOperation(&service, &operation)
	assert.NoErrorf(t, err, "DoOperation shouldn't return an error")

	expectedFunctions := []string{"StopInstance", "ChangeMachineType", "StartInstance"}
	expectedArguments := [][]interface{}{
		{"rightsizer-test", "us-east1-b", "alicja-test"},
		{"rightsizer-test", "us-east1-b", "alicja-test", "custom-2-5120"},
		{"rightsizer-test", "us-east1-b", "alicja-test"},
	}
	expectedResults := [][]interface{}{{nil}, {nil}, {nil}}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if the replace status operation works as expected.
func TestReplaceStatusOperation(t *testing.T) {
	operation := gcloudOperation{
		Action:       "replace",
		Path:         "/status",
		Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
		ResourceType: "compute.googleapis.com/Instance",
		Value:        "TERMINATED",
	}

	service := ApplyMockService{}
	err := DoOperation(&service, &operation)
	assert.NoErrorf(t, err, "DoOperation shouldn't return an error")

	expectedFunctions := []string{"StopInstance"}
	expectedArguments := [][]interface{}{{"rightsizer-test", "us-central1-a", "vkovalova-instance-memory-1"}}
	expectedResults := [][]interface{}{{nil}}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if the add snapshot operation works as expected.
func TestAddSnapshotOperation(t *testing.T) {
	operation := gcloudOperation{
		Action:       "add",
		Path:         "//compute.googleapis.com/projects/rightsizer-test/global/snapshots/$snapshot-name",
		ResourceType: "compute.googleapis.com/Snapshot",
		Value:        valueAddSnapshot{Name: "$snapshot-name", SourceDisk: "projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress", StorageLocations: []string{"europe-west1-d"}},
	}

	service := ApplyMockService{}
	err := DoOperation(&service, &operation)
	assert.NoErrorf(t, err, "DoOperation shouldn't return an error")

	expectedFunctions := []string{"CreateSnapshot"}
	expectedArguments := [][]interface{}{{"rightsizer-test", "europe-west1-d", "vertical-scaling-krzysztofk-wordpress", ""}}
	expectedResults := [][]interface{}{{nil}}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if the remove disk operation works as expected.
func TestRemoveDiskOperation(t *testing.T) {
	operation := gcloudOperation{
		Action:       "remove",
		Path:         "/",
		Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
		ResourceType: "compute.googleapis.com/Disk",
	}

	service := ApplyMockService{}
	err := DoOperation(&service, &operation)
	assert.NoErrorf(t, err, "DoOperation shouldn't return an error")

	expectedFunctions := []string{"DeleteDisk"}
	expectedArguments := [][]interface{}{{"rightsizer-test", "europe-west1-d", "vertical-scaling-krzysztofk-wordpress"}}
	expectedResults := [][]interface{}{{nil}}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if receiving an operation without necessary parameter
// returns the correct error.
func TestResourceWithoutNecessaryParams(t *testing.T) {
	operation := gcloudOperation{
		Action:       "remove",
		Path:         "/",
		Resource:     "//compute.googleapis.com/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
		ResourceType: "compute.googleapis.com/Disk",
	}

	service := ApplyMockService{}
	err := DoOperation(&service, &operation)
	assert.Error(t, err, fmt.Sprintf("url %s does not contain the parameter %s", operation.Resource, projectParam))

	expectedFunctions := []string{}
	expectedArguments := [][]interface{}{}
	expectedResults := [][]interface{}{}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if applying recommendation with stopping the machine
// works as expected.
func TestStopRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/status",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "RUNNING",
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/status",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "TERMINATED",
						},
					},
				},
			},
		},
		Etag:      "\"9f58395697934a1a\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.IdleResourceRecommender/recommendations/63378bdf-9ffe-4ea4-b8ee-04145f2a59c9",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{Status: "RUNNING"}}
	err := Apply(&service, &recommendation)
	assert.NoErrorf(t, err, "Apply shouldn't return an error")

	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"GetInstance",
		"StopInstance",
		"MarkRecommendationSucceeded",
	}
	expectedArguments := [][]interface{}{
		{"projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.IdleResourceRecommender/recommendations/63378bdf-9ffe-4ea4-b8ee-04145f2a59c9", "\"9f58395697934a1a\""},
		{"rightsizer-test", "us-central1-a", "vkovalova-instance-memory-1"},
		{"rightsizer-test", "us-central1-a", "vkovalova-instance-memory-1"},
		{"projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.IdleResourceRecommender/recommendations/63378bdf-9ffe-4ea4-b8ee-04145f2a59c9", "\"9f58395697934a1a\""},
	}
	expectedResults := [][]interface{}{
		{nil},
		{&compute.Instance{Status: "RUNNING"}, nil},
		{nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if applying recommmendation with adding snapshot of a machine
// and then deleting it works as expected.
func TestSnapshotAndDeleteRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "add",
							Path:         "/",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/global/snapshots/$snapshot-name",
							ResourceType: "compute.googleapis.com/Snapshot",
							Value:        valueAddSnapshot{Name: "$snapshot-name", SourceDisk: "projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress", StorageLocations: []string{"europe-west1-d"}},
						},
						&gcloudOperation{
							Action:       "remove",
							Path:         "/",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
							ResourceType: "compute.googleapis.com/Disk",
						},
					},
				},
			},
		},
		Etag:      "\"856260fc666866a3\"",
		Name:      "projects/323016592286/locations/europe-west1-d/recommenders/google.compute.disk.IdleResourceRecommender/recommendations/1e32196d-fc39-4358-9c9b-cec17a85f4ea",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{calledFunctions: []calledFunction{}}
	err := Apply(&service, &recommendation)
	assert.NoErrorf(t, err, "Apply shouldn't return an error")

	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"CreateSnapshot",
		"DeleteDisk",
		"MarkRecommendationSucceeded",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{"rightsizer-test", "europe-west1-d", "vertical-scaling-krzysztofk-wordpress", ""},
		{"rightsizer-test", "europe-west1-d", "vertical-scaling-krzysztofk-wordpress"},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{nil},
		{nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks if applying a recommendation with replacing machine type
// works as expected.
func TestReplaceRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{calledFunctions: []calledFunction{}, getInstanceResult: &compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}}
	err := Apply(&service, &recommendation)
	assert.NoErrorf(t, err, "Apply shouldn't return an error")

	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"GetInstance",
		"StopInstance",
		"ChangeMachineType",
		"StartInstance",
		"MarkRecommendationSucceeded",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{"rightsizer-test", "us-central1-a", "sidsharan-e2-with-stackdriver"},
		{"rightsizer-test", "us-central1-a", "sidharan-e2-with-stackdriver"},
		{"rightsizer-test", "us-central1-a", "sidharan-e2-with-stackdriver", "e2-medium"},
		{"rightsizer-test", "us-central1-a", "sidharan-e2-with-stackdriver"},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{&compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}, nil},
		{nil},
		{nil},
		{nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks, that the attempt to apply a not active recommendation
// results in the expected error.
func TestNotActiveRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Claimed"},
	}

	service := ApplyMockService{}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, "to apply a recommendation, its status must be active")
}

// Checks, that the attempt to apply a recommendation with unknown action
// results in the expected error.
func TestUnsupportedAction(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "copy",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, operationNotSupportedMessage)
	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks, that the attempt to apply a recommendation with unknown path
// results in the expected error.
func TestUnsupportedPath(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/coreCount",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, operationNotSupportedMessage)
	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"GetInstance",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{"rightsizer-test", "us-central1-a", "sidsharan-e2-with-stackdriver"},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{&compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}, nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks, that the attempt to apply a recommendation with
// unknown resource type results in the expected error.
func TestUnsupportedResourceType(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/CPU",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, operationNotSupportedMessage)
	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks, that the attempt to apply a recommendation with unknown
// unknown replace value results in the expected error.
func TestUnsupportedReplaceValue(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/status",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "RUNNING",
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/status",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "CLOSED",
						},
					},
				},
			},
		},
		Etag:      "\"da62b100443c341b\"",
		Name:      "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd692f-14b7-499a-be95-a09fe0893911",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, operationNotSupportedMessage)
	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"GetInstance",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{"rightsizer-test", "us-central1-a", "vkovalova-instance-memory-1"},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{&compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}, nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks, that the attempt to apply a recommendation that adds
// a resource of unknown type results in the expected error.
func TestUnsupportedAddResourceType(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "add",
							Path:         "/",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/global/snapshots/$snapshot-name",
							ResourceType: "compute.googleapis.com/Schnappschuss",
							Value: &valueAddSnapshot{
								Name:       "$snapshot-name",
								SourceDisk: "projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
								StorageLocations: []string{
									"europe-west1-d",
								},
							},
						},
						&gcloudOperation{
							Action:       "remove",
							Path:         "/",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
							ResourceType: "compute.googleapis.com/Disk",
						},
					},
				},
			},
		},
		Etag:      "\"da62b100443c341b\"",
		Name:      "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd692f-14b7-499a-be95-a09fe0893911",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, operationNotSupportedMessage)
	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

// Checks that when a test operation fails, the other one is not performed,
// and the expected error is returned.
func TestFailedTest(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := ApplyMockService{getInstanceResult: &compute.Instance{MachineType: "@#$%!E"}}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, "machine type is not as expected")
	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"GetInstance",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{"rightsizer-test", "us-central1-a", "sidsharan-e2-with-stackdriver"},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{&compute.Instance{MachineType: "@#$%!E"}, nil},
		{nil},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

type FailedClaimService struct {
	GoogleService
}

func (s *FailedClaimService) MarkRecommendationClaimed(name string, etag string) (*gcloudRecommendation, error) {
	return nil, errors.New("recommendation couldn't be marked claimed")
}

// Checks, that a failure to mark a recommendation as claimed
// leads to the expected behaviour.
func TestFailedClaimRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := FailedClaimService{}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, "recommendation couldn't be marked claimed")
}

type FailedSucceedService struct {
	GoogleService
	calledFunctions   []calledFunction
	getInstanceResult *compute.Instance
}

func (s *FailedSucceedService) GetInstance(project string, zone string, instance string) (*compute.Instance, error) {
	newCalledFunction := calledFunction{"GetInstance", []interface{}{project, zone, instance}, []interface{}{s.getInstanceResult, nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return s.getInstanceResult, nil
}

func (s *FailedSucceedService) StopInstance(project string, zone string, instance string) error {
	newCalledFunction := calledFunction{"StopInstance", []interface{}{project, zone, instance}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *FailedSucceedService) ChangeMachineType(project string, zone string, instance string, machineType string) error {
	newCalledFunction := calledFunction{"ChangeMachineType", []interface{}{project, zone, instance, machineType}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *FailedSucceedService) StartInstance(project string, zone string, instance string) error {
	newCalledFunction := calledFunction{"StartInstance", []interface{}{project, zone, instance}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *FailedSucceedService) MarkRecommendationClaimed(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationClaimed", []interface{}{name, etag}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *FailedSucceedService) MarkRecommendationSucceeded(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationSucceeded", []interface{}{name, etag}, []interface{}{errors.New("recommendation couldn't be marked succeeded")}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return errors.New("recommendation couldn't be marked succeeded")
}

// Checks that failing to mark a recommendation as succeeded leads
// to the expected behaviour.
func TestFailedSucceedRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "test",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := FailedSucceedService{getInstanceResult: &compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, "recommendation couldn't be marked succeeded")

	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"GetInstance",
		"StopInstance",
		"ChangeMachineType",
		"StartInstance",
		"MarkRecommendationSucceeded",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{"rightsizer-test", "us-central1-a", "sidsharan-e2-with-stackdriver"},
		{"rightsizer-test", "us-central1-a", "sidharan-e2-with-stackdriver"},
		{"rightsizer-test", "us-central1-a", "sidharan-e2-with-stackdriver", "e2-medium"},
		{"rightsizer-test", "us-central1-a", "sidharan-e2-with-stackdriver"},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{&compute.Instance{MachineType: "zones/us-east1-b/machineTypes/e2-standard-2"}, nil},
		{nil},
		{nil},
		{nil},
		{errors.New("recommendation couldn't be marked succeeded")},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}

type FailedFailedService struct {
	GoogleService
	calledFunctions   []calledFunction
	getInstanceResult *compute.Instance
}

func (s *FailedFailedService) MarkRecommendationClaimed(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationClaimed", []interface{}{name, etag}, []interface{}{nil}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return nil
}

func (s *FailedFailedService) MarkRecommendationFailed(name string, etag string) (*gcloudRecommendation, error) {
	newCalledFunction := calledFunction{"MarkRecommendationFailed", []interface{}{name, etag}, []interface{}{errors.New("recommendation couldn't be marked failed")}}
	s.calledFunctions = append(s.calledFunctions, newCalledFunction)
	return errors.New("recommendation couldn't be marked failed")
}

// Checks that failing to mark a recommendation as failed leads
// to the expected behaviour.
func TestFailedFailedRecommendation(t *testing.T) {
	recommendation := gcloudRecommendation{
		Content: &gcloudContent{
			OperationGroups: []*gcloudOperationGroup{
				&gcloudOperationGroup{
					Operations: []*gcloudOperation{
						&gcloudOperation{
							Action:       "copy",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							ValueMatcher: &gcloudValueMatcher{MatchesPattern: ".*zones/us-east1-b/machineTypes/e2-standard-2"},
						},
						&gcloudOperation{
							Action:       "replace",
							Path:         "/machineType",
							Resource:     "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidharan-e2-with-stackdriver",
							ResourceType: "compute.googleapis.com/Instance",
							Value:        "zones/us-central1-a/machineTypes/e2-medium",
						},
					},
				},
			},
		},
		Etag:      "\"40204a1000e5befe\"",
		Name:      "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
		StateInfo: &gcloudStateInfo{State: "Active"},
	}

	service := FailedFailedService{}
	err := Apply(&service, &recommendation)
	assert.Error(t, err, "recommendation couldn't be marked failed")

	expectedFunctions := []string{
		"MarkRecommendationClaimed",
		"MarkRecommendationFailed",
	}
	expectedArguments := [][]interface{}{
		{recommendation.Name, recommendation.Etag},
		{recommendation.Name, recommendation.Etag},
	}
	expectedResults := [][]interface{}{
		{nil},
		{errors.New("recommendation couldn't be marked failed")},
	}

	expected, _ := newCalledFunctions(expectedFunctions, expectedArguments, expectedResults)
	compareCalledFunctions(t, expected, service.calledFunctions)
}
