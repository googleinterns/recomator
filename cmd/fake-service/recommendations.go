package main

var recommendationsJSON = []string{`{
	"content": {
		"operationGroups": [
			{
				"operations": [
					{
						"action": "test",
						"path": "/machineType",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
						"resourceType": "compute.googleapis.com/Instance",
						"valueMatcher": {
							"matchesPattern": ".*zones/us-east1-b/machineTypes/n1-standard-4"
						}
					},
					{
						"action": "replace",
						"path": "/machineType",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
						"resourceType": "compute.googleapis.com/Instance",
						"value": "zones/us-east1-b/machineTypes/custom-2-5120"
					}
				]
			}
		]
	},
	"description": "Save cost by changing machine type from n1-standard-4 to custom-2-5120.",
	"etag": "\"da62b100443c341b\"",
	"lastRefreshTime": "2020-07-13T06:41:17Z",
	"name": "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd692f-14b7-499a-be95-a09fe0893911",
	"primaryImpact": {
		"category": "COST",
		"costProjection": {
			"cost": {
				"currencyCode": "USD",
				"nanos": -268972762,
				"units": "-73"
			},
			"duration": "2592000s"
		}
	},
	"recommenderSubtype": "CHANGE_MACHINE_TYPE",
	"stateInfo": {
		"state": "CLAIMED"
	}
}`,
	`{
	"content": {
		"operationGroups": [
			{
				"operations": [
					{
						"action": "test",
						"path": "/machineType",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
						"resourceType": "compute.googleapis.com/Instance",
						"valueMatcher": {
							"matchesPattern": ".*zones/us-central1-a/machineTypes/e2-standard-2"
						}
					},
					{
						"action": "replace",
						"path": "/machineType",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/sidsharan-e2-with-stackdriver",
						"resourceType": "compute.googleapis.com/Instance",
						"value": "zones/us-central1-a/machineTypes/e2-medium"
					}
				]
			}
		]
	},
	"description": "Save cost by changing machine type from e2-standard-2 to e2-medium.",
	"etag": "\"40204a1000e5befe\"",
	"lastRefreshTime": "2020-06-24T06:20:37Z",
	"name": "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/5df355d9-2f50-4567-a909-bcfcebcf7d66",
	"primaryImpact": {
		"category": "COST",
		"costProjection": {
			"cost": {
				"currencyCode": "USD",
				"nanos": -79835798,
				"units": "-24"
			},
			"duration": "2592000s"
		}
	},
	"recommenderSubtype": "CHANGE_MACHINE_TYPE",
	"stateInfo": {
		"state": "FAILED"
	}
}`,

	`{
	"content": {
		"operationGroups": [
			{
				"operations": [
					{
						"action": "test",
						"path": "/status",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
						"resourceType": "compute.googleapis.com/Instance",
						"value": "RUNNING"
					},
					{
						"action": "replace",
						"path": "/status",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/instances/vkovalova-instance-memory-1",
						"resourceType": "compute.googleapis.com/Instance",
						"value": "TERMINATED"
					}
				]
			}
		]
	},
	"description": "Save cost by stopping Idle VM 'vkovalova-instance-memory-1'.",
	"etag": "\"9f58395697934a1a\"",
	"lastRefreshTime": "2020-07-17T06:48:28Z",
	"name": "projects/323016592286/locations/us-central1-a/recommenders/google.compute.instance.IdleResourceRecommender/recommendations/63378bdf-9ffe-4ea4-b8ee-04145f2a59c9",
	"primaryImpact": {
		"category": "COST",
		"costProjection": {
			"cost": {
				"currencyCode": "USD",
				"nanos": -497242857,
				"units": "-5"
			},
			"duration": "2592000s"
		}
	},
	"recommenderSubtype": "STOP_VM",
	"stateInfo": {
		"state": "ACTIVE"
	}
}`,
	`{
	"content": {
		"operationGroups": [
			{
				"operations": [
					{
						"action": "test",
						"path": "/status",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-c/instances/timus-test-for-probers-n2-std-4-idling",
						"resourceType": "compute.googleapis.com/Instance",
						"value": "RUNNING"
					},
					{
						"action": "replace",
						"path": "/status",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-c/instances/timus-test-for-probers-n2-std-4-idling",
						"resourceType": "compute.googleapis.com/Instance",
						"value": "TERMINATED"
					}
				]
			}
		]
	},
	"description": "Save cost by stopping Idle VM 'timus-test-for-probers-n2-std-4-idling'.",
	"etag": "\"2e11293786a101ea\"",
	"lastRefreshTime": "2020-07-17T06:36:44Z",
	"name": "projects/323016592286/locations/us-central1-c/recommenders/google.compute.instance.IdleResourceRecommender/recommendations/6df88342-8116-441c-beb7-ab66d18a3078",
	"primaryImpact": {
		"category": "COST",
		"costProjection": {
			"cost": {
				"currencyCode": "USD",
				"nanos": -182895999,
				"units": "-140"
			},
			"duration": "2592000s"
		}
	},
	"recommenderSubtype": "STOP_VM",
	"stateInfo": {
		"state": "ACTIVE"
	}
}`,
	`{
	"associatedInsights": [
		{
			"insight": "projects/323016592286/locations/europe-west1-d/insightTypes/google.compute.disk.IdleResourceInsight/insights/620de196-ff06-4f97-ae52-648636c98c49"
		}
	],
	"content": {
		"operationGroups": [
			{
				"operations": [
					{
						"action": "add",
						"path": "/",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/global/snapshots/$snapshot-name",
						"resourceType": "compute.googleapis.com/Snapshot",
						"value": {
							"name": "$snapshot-name",
							"source_disk": "projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
							"storage_locations": [
								"europe-west1-d"
							]
						}
					},
					{
						"action": "remove",
						"path": "/",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
						"resourceType": "compute.googleapis.com/Disk"
					}
				]
			}
		]
	},
	"description": "Save cost by snapshotting and then deleting idle persistent disk 'vertical-scaling-krzysztofk-wordpress'.",
	"etag": "\"856260fc666866a3\"",
	"lastRefreshTime": "2020-07-17T07:00:00Z",
	"name": "projects/323016592286/locations/europe-west1-d/recommenders/google.compute.disk.IdleResourceRecommender/recommendations/1e32196d-fc39-4358-9c9b-cec17a85f4ea",
	"primaryImpact": {
		"category": "COST",
		"costProjection": {
			"cost": {
				"currencyCode": "USD",
				"nanos": -135483871
			},
			"duration": "2592000s"
		}
	},
	"recommenderSubtype": "SNAPSHOT_AND_DELETE_DISK",
	"stateInfo": {
		"state": "ACTIVE"
	}
}`,
	`{
	"associatedInsights": [
		{
			"insight": "projects/323016592286/locations/europe-west1-d/insightTypes/google.compute.disk.IdleResourceInsight/insights/2afa4c16-7812-4469-a56e-03099d83447b"
		}
	],
	"content": {
		"operationGroups": [
			{
				"operations": [
					{
						"action": "add",
						"path": "/",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/global/snapshots/$snapshot-name",
						"resourceType": "compute.googleapis.com/Snapshot",
						"value": {
							"name": "$snapshot-name",
							"source_disk": "projects/rightsizer-test/zones/europe-west1-d/disks/krzysztofk2",
							"storage_locations": [
								"europe-west1-d"
							]
						}
					},
					{
						"action": "remove",
						"path": "/",
						"resource": "//compute.googleapis.com/projects/rightsizer-test/zones/europe-west1-d/disks/krzysztofk2",
						"resourceType": "compute.googleapis.com/Disk"
					}
				]
			}
		]
	},
	"description": "Save cost by snapshotting and then deleting idle persistent disk 'krzysztofk2'.",
	"etag": "\"4159bea4e7c90c00\"",
	"lastRefreshTime": "2020-07-17T07:00:00Z",
	"name": "projects/323016592286/locations/europe-west1-d/recommenders/google.compute.disk.IdleResourceRecommender/recommendations/8962f57e-10c6-47cc-a48e-52e9f0f800c5",
	"primaryImpact": {
		"category": "COST",
		"costProjection": {
			"cost": {
				"currencyCode": "USD",
				"nanos": -354838709,
				"units": "-1"
			},
			"duration": "2592000s"
		}
	},
	"recommenderSubtype": "SNAPSHOT_AND_DELETE_DISK",
	"stateInfo": {
		"state": "ACTIVE"
	}
}`}
