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

import { extractFromResource } from "./utils";

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
  // This is a part of a union field, which other version
  // (ValueMatcher) is not currently in use
  value?: string | AddOperationValue;
}

export interface AddOperationValue {
  name: string;
  source_disk: string;
  storage_locations: string[];
}

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
export function getExampleRecommendationAction(
  recommendation: RecommendationRaw
): string {
  return recommendation.content.operationGroups[0].operations[0].action;
}

// -> "//compute.googleapis.com/projects/rightsizer-test/zones/us-east1-b/instances/alicja-test"
export function getExampleRecommendationResource(
  recommendation: RecommendationRaw
): string {
  return recommendation.content.operationGroups[0].operations[0].resource;
}

// -> "compute.googleapis.com/Snapshot"
export function getExampleRecommendationResourceType(
  recommendation: RecommendationRaw
): string {
  return recommendation.content.operationGroups[0].operations[0].resourceType;
}

export function getExampleRecommendationValue(
  recommendation: RecommendationRaw
): string | AddOperationValue | undefined {
  return recommendation.content.operationGroups[0].operations[0].value;
}

// -> "timus-test-for-probers-n2-std-4-idling"
export function getRecommendationResourceShortName(
  recommendation: RecommendationRaw
): string {
  const action = getExampleRecommendationAction(recommendation);

  switch (action) {
    case "add": {
      const value = getExampleRecommendationValue(recommendation);
      if (isAddOperationValue(value)) {
        return extractFromResource("disks", value.source_disk);
      }

      throw "the given value parameter doesn't match the action";
    }
    case "remove": {
      const resource = getExampleRecommendationResource(recommendation);
      return extractFromResource("disks", resource);
    }
    case "replace": {
      const resource = getExampleRecommendationResource(recommendation);
      return extractFromResource("instances", resource);
    }
    case "test": {
      const resource = getExampleRecommendationResource(recommendation);
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
  const resource = getExampleRecommendationResource(recommendation);
  return extractFromResource("projects", resource);
}

// TODO: remove ignoring Eslint, once these methods are actually used somewhere

// Doesn't do much, but I think it is likely we will decide to show more clever descriptions later
export function getRecomendationDescription(
  recommendation: RecommendationRaw
): string {
  return recommendation.description;
}

// "3.5" ($ per week)
// eslint-disable-next-line @typescript-eslint/no-unused-vars
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

export const internalStatusMap: Record<string, string> = {
  ACTIVE: "Applicable",
  CLAIMED: "In progress",
  SUCCEEDED: "Success",
  FAILED: "Failed",
  DISMISSED: "Dismissed"
};

export function throwIfInvalidStatus(statusName: string): void {
  if (!(statusName in internalStatusMap))
    throw `invalid status name passed: ${statusName}`;
}

export function getInternalStatusMapping(statusName: string): string {
  throwIfInvalidStatus(statusName);
  return internalStatusMap[statusName];
}

// All data maintained for each recommendation
export class RecommendationExtra implements RecommendationRaw {
  // These should not be modified (including inner fields) outside of tests:
  readonly name: string;
  readonly description: string;
  readonly recommenderSubtype: string;
  readonly primaryImpact: Impact;
  readonly content: RecommendationContent;
  readonly stateInfo: RecommendationStateInfo; // original status

  // need to remember them so that v-data-table knows what to sort by
  readonly costCol: number;
  readonly projectCol: string;
  readonly resourceCol: string;
  readonly typeCol: string;

  // These can be modified:
  statusCol: string; // follows the current recommendation status
  errorHeader?: string;
  errorDescription?: string;

  constructor(rec: RecommendationRaw) {
    this.name = rec.name;
    this.description = rec.description;
    this.recommenderSubtype = rec.recommenderSubtype;
    this.primaryImpact = rec.primaryImpact;
    this.content = rec.content;
    this.stateInfo = rec.stateInfo;

    this.costCol = getRecommendationCostPerWeek(rec);
    this.projectCol = getRecommendationProject(rec);
    this.resourceCol = getRecommendationResourceShortName(rec);
    this.typeCol = getRecommendationType(rec);
    this.statusCol = getInternalStatusMapping(rec.stateInfo.state);
  }
}
