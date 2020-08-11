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

import {
  coreTableStoreStateFactory,
  ICoreTableStoreState
} from "@/store/core_table";

import {
  projectFilterAccepted,
  typeFilterAccepted,
  statusFilterAccepted
} from "@/store/core_table_filter_utils";
import { freshSampleRawRecommendation } from "./sample_recommendation";
import {
  RecommendationExtra,
  RecommendationRaw,
  getInternalStatusMapping
} from "@/store/model";

describe("Filtering predicates individually", () => {
  // Not using actual Vuex instances here, so we don't need to use mutations
  let tableState: ICoreTableStoreState;
  let recommendation: RecommendationRaw;
  const extra = () => new RecommendationExtra(recommendation);

  beforeEach(() => {
    tableState = coreTableStoreStateFactory();
    recommendation = freshSampleRawRecommendation();
  });

  test("project filter", () => {
    recommendation.content.operationGroups[0].operations[0].resource =
      "//compute.googleapis.com/projects/Facebook/zones/us-east1-b/instances/bob-vm0";

    tableState.projectsSelected = [];
    expect(projectFilterAccepted(tableState, extra())).toBeTruthy();

    tableState.projectsSelected = ["tiktok", "instagram", "Facebook"];
    expect(projectFilterAccepted(tableState, extra())).toBeTruthy();

    // Now let's change the project to "instagram"
    recommendation.content.operationGroups[0].operations[0].resource =
      "//compute.googleapis.com/projects/instagram/zones/us-east1-b/instances/bob-vm0";

    tableState.projectsSelected = ["tiktok", "instagram", "snapchat"];
    expect(projectFilterAccepted(tableState, extra())).toBeTruthy();

    tableState.projectsSelected = ["tiktok", "snapchat"];
    expect(projectFilterAccepted(tableState, extra())).toBeFalsy();

    tableState.projectsSelected = ["instagram".toUpperCase()];
    expect(projectFilterAccepted(tableState, extra())).toBeFalsy();
  });

  // These are very similar, so I will not test them as thoroughly,
  //   as they have nearly identical implementations and it is easy to spot differences.
  // I also don't want to introduce any more abstractions so that the code
  //   is still easy to read.
  test("type filter", () => {
    recommendation.recommenderSubtype = "SWITCH_TO_AWS";

    tableState.typesSelected = [];
    expect(typeFilterAccepted(tableState, extra())).toBeTruthy();

    tableState.typesSelected = ["CHANGE_MACHINE_TYPE"];
    expect(typeFilterAccepted(tableState, extra())).toBeFalsy();

    recommendation.recommenderSubtype = "CHANGE_MACHINE_TYPE";
    // what if one of the types is empty?
    tableState.typesSelected = ["", "SWITCH_TO_AWS", "CHANGE_MACHINE_TYPE"];
    expect(typeFilterAccepted(tableState, extra())).toBeTruthy();
  });

  test("status filter", () => {
    recommendation.stateInfo.state = "FAILED"; // this maps to "Failed"

    tableState.statusesSelected = [];
    expect(statusFilterAccepted(tableState, extra())).toBeTruthy();

    // what if statuses repeat ?
    tableState.statusesSelected = [
      "SUCCEEDED",
      "SUCCEEDED",
      "CLAIMED"
    ].map(status => getInternalStatusMapping(status));
    expect(statusFilterAccepted(tableState, extra())).toBeFalsy();

    recommendation.stateInfo.state = "CLAIMED"; // this maps to "In progress"
    expect(statusFilterAccepted(tableState, extra())).toBeTruthy();
  });
});
