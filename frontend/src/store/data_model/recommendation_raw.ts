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

// RecommendationRaw definition and parsers
// Follows data model from:
//    https://cloud.google.com/recommender/docs/reference/rest/v1beta1/projects.locations.recommenders.recommendations

export interface RecommendationRaw {
  name: string;
  description: string;
  recommenderSubtype: string;
  primaryImpact: Impact;
  content: RecommendationContent;
  stateInfo: RecommendationStateInfo;
}

export interface Impact {
  category: string; // originally enum
  costProjection: CostProjection;
}

export interface CostProjection {
  cost: Money;
  duration: string;
}

export interface Money {
  currencyCode: string;
  units?: string;
  nanos?: number;
}

export interface RecommendationStateInfo {
  state: string; // originally enum
}

export interface RecommendationContent {
  operationGroups: OperationGroupsList;
}

export interface OperationGroupsList {
  [index: number]: OperationGroup;
}

export interface OperationGroup {
  operations: OperationsList;
}

export interface OperationsList {
  [index: number]: Operation;
}

export interface Operation {
  action: string;
  resource: string;
  resourceType: string;
  path: string;
  // This is a part of a union field, which other type ValueMatcher is not currently in use
  // This might well fail to parse for non-standard (not seen in recommendations now) operations
  value?: string | AddOperationValue;
}

// Operation value used for snapshots
export interface AddOperationValue {
  name: string;
  source_disk: string;
  storage_locations: string[];
}

// Checks if the operation looks like the first operation in a snapshot recommendation
function isAddOperationValue(
  value: string | AddOperationValue | undefined
): value is AddOperationValue {
  if (typeof value !== "object") {
    return false;
  }
  const valueExpectedProperties = ["name", "source_disk", "storage_locations"];

  for (const property in value) {
    if (Object.prototype.hasOwnProperty.call(value, property)) {
      if (!valueExpectedProperties.includes(property)) {
        return false;
      }
    }
  }

  for (const property of valueExpectedProperties) {
    if (!Object.prototype.hasOwnProperty.call(value, property)) {
      return false;
    }
  }

  if (typeof value.name !== "string") {
    return false;
  }

  if (typeof value.source_disk !== "string") {
    return false;
  }

  if (typeof value.storage_locations !== "object") {
    return false;
  }

  for (const field of value.storage_locations) {
    if (typeof field !== "string") {
      return false;
    }
  }

  return true;
}

// -> "add"
export function getRecommendationFirstAction(
  recommendation: RecommendationRaw
): string {
  return recommendation.content.operationGroups[0].operations[0].action;
}

// -> "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test"
export function getRecommendationFirstResource(
  recommendation: RecommendationRaw
): string {
  return recommendation.content.operationGroups[0].operations[0].resource;
}

// -> "compute.googleapis.com/Snapshot"
export function getRecommendationFirstResourceType(
  recommendation: RecommendationRaw
): string {
  return recommendation.content.operationGroups[0].operations[0].resourceType;
}

export function getRecommendationFirstValue(
  recommendation: RecommendationRaw
): string | AddOperationValue | undefined {
  return recommendation.content.operationGroups[0].operations[0].value;
}

// (b, /a/b/c/d/e) => c
export function extractFromResource(
  property: string,
  resource: string
): string {
  const sliceLen = `/${property}/`.length;
  const pattern = `/${property}/[^/]*`;
  const regex = new RegExp(pattern);

  const found = regex.exec(resource);
  if (found === null) {
    throw `couldn't parse identifier: ${resource}`;
  }

  const result = found[0].slice(sliceLen);
  return result;
}

// Returns a name to identify the related resource by, regardless of recommendation type.
// -> "timus-test-for-probers-n2-std-4-idling"
// -> "shcheshnyak-disk"
export function getRecommendationResourceShortName(
  recommendation: RecommendationRaw
): string {
  const action = getRecommendationFirstAction(recommendation);

  switch (action) {
    case "add": {
      const value = getRecommendationFirstValue(recommendation);
      if (isAddOperationValue(value)) {
        return extractFromResource("disks", value.source_disk);
      }

      throw "the given value parameter doesn't match the action";
    }
    case "remove": {
      const resource = getRecommendationFirstResource(recommendation);
      return extractFromResource("disks", resource);
    }
    case "replace": {
      const resource = getRecommendationFirstResource(recommendation);
      return extractFromResource("instances", resource);
    }
    case "test": {
      const resource = getRecommendationFirstResource(recommendation);
      return extractFromResource("instances", resource);
    }
    default:
      throw "the given recommendation contains an unsupported action";
  }
}

// -> "rightsizer-test"
export function getRecommendationProject(
  recommendation: RecommendationRaw
): string {
  const resource = getRecommendationFirstResource(recommendation);
  return extractFromResource("projects", resource);
}

// -> "Save cost by changing machine type from n1-standard-4 to custom-2-5120."
export function getRecomendationDescription(
  recommendation: RecommendationRaw
): string {
  return recommendation.description;
}

// "3.5" ($ per week)
export function getRecommendationCostPerWeek(
  recommendation: RecommendationRaw
): number {
  const costObject = recommendation.primaryImpact.costProjection.cost;
  console.assert(
    costObject.currencyCode === "USD",
    "Only USD supported, got %s",
    costObject.currencyCode
  );

  // As a month doesn't have a fixed number of seconds, weekly cost is used
  // example duration: "2592000s"
  const secs = parseInt(
    recommendation.primaryImpact.costProjection.duration.slice(0, -1)
  );
  // example units: "-73", example nanos: 4200000000 => -73.42
  // Sometimes we will only get the 'nanos' field
  //  (only observed for snaphshots for now)

  let cost = 0;

  if (costObject.units !== undefined) cost += parseInt(costObject.units!);
  if (costObject.nanos !== undefined)
    cost += costObject.nanos! / (1000 * 1000 * 1000);

  return (cost * 60 * 60 * 24 * 7) / secs;
}

// "CHANGE_MACHINE_TYPE", "INCREASE_PERFORMANCE", ...
export function getRecommendationType(recommendation: RecommendationRaw) {
  return recommendation.recommenderSubtype;
}
