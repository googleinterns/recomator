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
  category: string; // originally enum
  costProjection: CostProjection;
}

interface CostProjection {
  cost: Money;
  duration: string;
}

interface Money {
  currencyCode: string;
  units: string;
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

/* TODO: remove the placeholders and implement these

function getRecommendationProject(recommendation: Recommendation): string {
  return "rightsizer-test";
}

function getRecommendationResourceName(recommendation: Recommendation): string {
  return "timus-test-for-probers-n2-std-4-idling";
}*/

// Doesn't do much, but I think it is likely we will decide to show more clever descriptions later
function getRecomendationDescription(recommendation: Recommendation): string {
  return recommendation.description;
}

// "3.5$ per week"
function getRecommendationCostString(recommendation: Recommendation): string {
  console.assert(
    recommendation.primaryImpact.costProjection.cost.currencyCode === "USD",
    "Only USD supported, got %s",
    recommendation.primaryImpact.costProjection.cost.currencyCode
  );

  // As a month doesn't have a fixed number of seconds, weekly cost is used

  // example duration: "2592000s"
  const secs = parseInt(
    recommendation.primaryImpact.costProjection.duration.slice(0, -1)
  );
  // example units: "-73"
  const cost = parseInt(
    recommendation.primaryImpact.costProjection.cost.units
  );

  const costPerWeek = (cost * 60 * 60 * 24 * 7) / secs;

  return `${costPerWeek.toFixed(2)} USD per week`;
}
