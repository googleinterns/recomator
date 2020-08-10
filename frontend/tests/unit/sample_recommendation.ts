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

import { RecommendationRaw } from "@/store/model";

const sampleRawRecommendation: RecommendationRaw = {
  content: {
    operationGroups: [
      {
        operations: [
          {
            path: "/machineType",
            resource:
              "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test",
            resourceType: "compute.googleapis.com/Instance"
          },
          {
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
    "projects/323016592286/locations/us-east1-b/recommenders/google.compute.instance.MachineTypeRecommender/recommendations/6dfd692f-14b7-499a-be95-a09fe0893911",
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

export function freshSampleRawRecommendation(): RecommendationRaw {
  // deep copy
  return JSON.parse(
    JSON.stringify(sampleRawRecommendation)
  ) as RecommendationRaw;
}
