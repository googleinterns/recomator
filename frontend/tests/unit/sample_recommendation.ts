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
        nanos: 268972762,
        units: "73"
      },
      duration: "2592000s"
    }
  },
  recommenderSubtype: "CHANGE_MACHINE_TYPE",
  stateInfo: {
    state: "CLAIMED"
  }
} as RecommendationRaw;

// only works for simple objects, maps and functions will be lost
function deepCopy(obj: object): object {
  return JSON.parse(JSON.stringify(obj));
}

export function freshSavingRawRecommendation(): RecommendationRaw {
  return deepCopy(savingRawRecommendation) as RecommendationRaw;
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
};

export function freshPerformanceRawRecommendation(): RecommendationRaw {
  return deepCopy(performanceRawRecommendation) as RecommendationRaw;
}

const performanceRawRecommendation: RecommendationRaw = {
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
    "Increase performance by changing machine type from custom-2-5120 to n1-standard-4.",
  name:
    "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd6d2137-14b7-499a-be95-a09fe0893911",
  primaryImpact: {
    category: "COST",
    costProjection: {
      cost: {
        currencyCode: "USD",
        units: "145"
      },
      duration: "2592000s"
    }
  },
  recommenderSubtype: "CHANGE_MACHINE_TYPE",
  stateInfo: {
    state: "CLAIMED"
  }
};

export function freshSampleRawRecommendation(): RecommendationRaw {
  return deepCopy(sampleRawRecommendation) as RecommendationRaw;
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
