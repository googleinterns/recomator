/* Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

import { RecommendationRaw } from "@/store/data_model/recommendation_raw"; // --> OFF

// We don't want to enforce camelCase here
/* eslint @typescript-eslint/camelcase: 0 */

// only works for simple objects, maps and functions will be lost
function deepCopy(obj: object): object {
  return JSON.parse(JSON.stringify(obj));
}

const sampleRawRecommendation: RecommendationRaw = {
  content: {
    operationGroups: [
      {
        operations: [
          {
            action: "test",
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
            resourceType: "compute.googleapis.com/Instance",
            valueMatcher: {
              matchesPattern: ".*zones/us-east1-b/machineTypes/n1-standard-4"
            }
          },
          {
            action: "replace",
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
            resourceType: "compute.googleapis.com/Instance",
            value: "zones/us-east1-b/machineTypes/custom-2-5120"
          }
        ]
      }
    ]
  },
  description:
    "Save cost by changing machine type from n1-standard-4 to custom-2-5120.",
  etag: '"da62b100443c341b"',
  lastRefreshTime: "2020-07-13T06:41:17Z",
  name:
    "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd692f-14b7-499a-be95-a09fe0893911",
  primaryImpact: {
    category: "COST",
    costProjection: {
      cost: {
        currencyCode: "USD",
        nanos: -268972762,
        units: "-73"
      },
      duration: "2592000s"
    }
  },
  recommenderSubtype: "CHANGE_MACHINE_TYPE",
  stateInfo: {
    state: "CLAIMED"
  }
} as RecommendationRaw;

export function freshSampleRawRecommendation(): RecommendationRaw {
  return deepCopy(sampleRawRecommendation) as RecommendationRaw;
}

const savingRawRecommendation: RecommendationRaw = {
  content: {
    operationGroups: [
      {
        operations: [
          {
            action: "test",
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
            resourceType: "compute.googleapis.com/Instance"
          },
          {
            action: "replace",
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
            resourceType: "compute.googleapis.com/Instance"
          }
        ]
      }
    ]
  },
  description:
    "Save cost by changing machine type from n1-standard-4 to custom-2-5120.",
  name:
    "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd692f-14b7-499a-be95-a21370893911",
  primaryImpact: {
    category: "COST",
    costProjection: {
      cost: {
        currencyCode: "USD",
        units: "-73"
      },
      duration: "2592000s"
    }
  },
  recommenderSubtype: "CHANGE_MACHINE_TYPE",
  stateInfo: {
    state: "CLAIMED"
  }
} as RecommendationRaw;

export function freshSavingRawRecommendation(): RecommendationRaw {
  return deepCopy(savingRawRecommendation) as RecommendationRaw;
}

// increase performance type (CHANGE_MACHINE_TYPE)
const performanceRawRecommendation: RecommendationRaw = {
  additionalImpact: [
    {
      category: "COST",
      costProjection: {
        cost: { currencyCode: "USD", nanos: 417998195, units: "72" },
        duration: "2592000s"
      }
    }
  ],
  content: {
    operationGroups: [
      {
        operations: [
          {
            action: "test",
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/asia-east1-c/instances/timus-test-e2-24-cores-3",
            resourceType: "compute.googleapis.com/Instance",
            valueMatcher: {
              matchesPattern:
                ".*zones/asia-east1-c/machineTypes/e2-custom-24-98304"
            }
          },
          {
            action: "replace",
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/asia-east1-c/instances/timus-test-e2-24-cores-3",
            resourceType: "compute.googleapis.com/Instance",
            value: "zones/asia-east1-c/machineTypes/e2-custom-28-98304"
          }
        ]
      }
    ]
  },
  description:
    "Improve performance by changing machine type from e2-custom-24-98304 to e2-custom-28-98304.",
  etag: '"b64ee9c5f53fa731"',
  lastRefreshTime: "2020-09-11T06:34:11Z",
  name:
    "projects/323016592286/locations/asia-east1-c/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/1ec7145b-3f6b-44b0-89f3-778aa3c3cc46",
  primaryImpact: { category: "PERFORMANCE" },
  recommenderSubtype: "CHANGE_MACHINE_TYPE",
  stateInfo: { state: "ACTIVE" }
} as RecommendationRaw;

export function freshPerformanceRawRecommendation(): RecommendationRaw {
  return deepCopy(performanceRawRecommendation) as RecommendationRaw;
}

const sampleSnapshotRawRecommendation: RecommendationRaw = {
  associatedInsights: [
    {
      insight:
        "projects/323016592286/locations/europe-west1-d/insightTypes/google.compute.disk.IdleResourceInsight/insights/620de196-ff06-4f97-ae52-648636c98c49"
    }
  ],
  content: {
    operationGroups: [
      {
        operations: [
          {
            action: "add",
            path: "/",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/global/snapshots/$snapshot-name",
            resourceType: "compute.googleapis.com/Snapshot",
            value: {
              name: "$snapshot-name",
              source_disk:
                "projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
              storage_locations: ["europe-west1-d"]
            }
          },
          {
            action: "remove",
            path: "/",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress",
            resourceType: "compute.googleapis.com/Disk"
          }
        ]
      }
    ]
  },
  description:
    "Save cost by snapshotting and then deleting idle persistent disk 'vertical-scaling-krzysztofk-wordpress'.",
  etag: '"856260fc666866a3"',
  lastRefreshTime: "2020-07-17T07:00:00Z",
  name:
    "projects/323016592286/locations/europe-west1-d/recommenders/google.compute.disk.IdleResourceRecommender/recommendations/1e32196d-fc39-4358-9c9b-cec17a85f4ea",
  primaryImpact: {
    category: "COST",
    costProjection: {
      cost: {
        currencyCode: "USD",
        nanos: -135483871
      },
      duration: "2592000s"
    }
  },
  recommenderSubtype: "SNAPSHOT_AND_DELETE_DISK",
  stateInfo: {
    state: "ACTIVE"
  }
} as RecommendationRaw;

export function freshSampleSnapshotRawRecommendation(): RecommendationRaw {
  return deepCopy(sampleSnapshotRawRecommendation) as RecommendationRaw;
}

const sampleStopVMRawRecommendation: RecommendationRaw = {
  content: {
    operationGroups: [
      {
        operations: [
          {
            action: "test",
            path: "/status",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-c/instances/timus-test-for-probers-n2-std-4-idling",
            resourceType: "compute.googleapis.com/Instance",
            value: "RUNNING"
          },
          {
            action: "replace",
            path: "/status",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-c/instances/timus-test-for-probers-n2-std-4-idling",
            resourceType: "compute.googleapis.com/Instance",
            value: "TERMINATED"
          }
        ]
      }
    ]
  },
  description:
    "Save cost by stopping Idle VM 'timus-test-for-probers-n2-std-4-idling'.",
  etag: '"2e11293786a101ea"',
  lastRefreshTime: "2020-07-17T06:36:44Z",
  name:
    "projects/323016592286/locations/us-central1-c/recommenders/google.compute.instance.IdleResourceRecommender/recommendations/6df88342-8116-441c-beb7-ab66d18a3078",
  primaryImpact: {
    category: "COST",
    costProjection: {
      cost: {
        currencyCode: "USD",
        nanos: -182895999,
        units: "-140"
      },
      duration: "2592000s"
    }
  },
  recommenderSubtype: "STOP_VM",
  stateInfo: {
    state: "ACTIVE"
  }
} as RecommendationRaw;

export function freshSampleStopVMRawRecommendation(): RecommendationRaw {
  return deepCopy(sampleStopVMRawRecommendation) as RecommendationRaw;
}

const sampleDeleteDiskRecommendation: RecommendationRaw = {
  name:
    "projects/323016592286/locations/us-central1-a/recommenders/google.compute.disk.IdleResourceRecommender/recommendations/33d373d1-e6ad-45b8-991a-83d9dcdb5ea5",
  description:
    "Save cost by deleting idle persistent disk 'stanislawm-test-1'.",
  recommenderSubtype: "DELETE_DISK",
  primaryImpact: {
    category: "COST",
    costProjection: {
      cost: {
        currencyCode: "USD",
        nanos: -400000000
      },
      duration: "2592000s"
    }
  },
  content: {
    operationGroups: [
      {
        operations: [
          {
            action: "remove",
            path: "/",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-central1-a/disks/stanislawm-test-1",
            resourceType: "compute.googleapis.com/Disk"
          }
        ]
      }
    ]
  },
  stateInfo: {
    state: "ACTIVE"
  }
} as RecommendationRaw;

export function freshSampleDeleteDiskRawRecommendation(): RecommendationRaw {
  return deepCopy(sampleDeleteDiskRecommendation) as RecommendationRaw;
}
