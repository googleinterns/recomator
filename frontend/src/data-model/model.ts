/* Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the License);
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an AS IS BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

// Follows data model from:
//    https://cloud.google.com/recommender/docs/reference/rest/v1beta1/projects.locations.recommenders.recommendations

interface Recommendation {
  name: string;
  description: string;
  recommenderSubtype: string;
  primaryImpact: Impact;
  content: RecommendationContent;
  stateInfo: RecommendationStateInfo;
}

interface Impact {
  category: string; // Originally enum
  costProjection: CostProjection;
}

interface CostProjection {
  cost: Money;
  duration: Duration;
}

interface Money {
  currencyCode: string;
  units: number;
}

interface Duration {
  seconds: number; // originally int64
}

interface RecommendationStateInfo {
  state: string; // originally enum
}

interface RecommendationContent {
  operationGroups: OperationGroupsList;
}

interface OperationGroupsList {
  [index: number]: OperationGroup;
}

interface OperationGroup {
  operations: OperationsList;
}

interface OperationsList {
  [index: number]: Operation;
}

interface Operation {
  resource: string;
  resourceType: string;
  path: string;
}
