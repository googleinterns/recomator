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

jest.mock("@/store/utils/core_table_filter_utils");

import {
  descriptionFilterAccepted,
  statusFilterAccepted,
  projectFilterAccepted,
  typeFilterAccepted,
  resourceFilterAccepted,
  costFilterAccepted
} from "@/store/utils/core_table_filter_utils";

import {
  coreTableStoreStateFactory,
  isRecommendationInResults
} from "@/store/core_table";

import { freshSampleRawRecommendation } from "./sample_recommendation";
import { RecommendationExtra } from "@/store/recommendation_extra";

describe("Filtering aggregate (resource name, type, status...)", () => {
  const tableState = coreTableStoreStateFactory();
  const recommendation = freshSampleRawRecommendation();
  const passesAll = () =>
    isRecommendationInResults(
      tableState,
      new RecommendationExtra(recommendation)
    ) as boolean;

  test("project filter", () => {
    const numberOfFilters = 6;
    let isOn = new Array(numberOfFilters).fill(true);
    (descriptionFilterAccepted as any).mockImplementation(() => isOn[0]);
    (projectFilterAccepted as any).mockImplementation(() => isOn[1]);
    (typeFilterAccepted as any).mockImplementation(() => isOn[2]);
    (resourceFilterAccepted as any).mockImplementation(() => isOn[3]);
    (statusFilterAccepted as any).mockImplementation(() => isOn[4]);
    (costFilterAccepted as any).mockImplementation(() => isOn[5]);

    for (let i = 0; i < numberOfFilters; i++) {
      expect(passesAll()).toBeTruthy();
      isOn[i] = false;
      expect(passesAll()).toBeFalsy();
      isOn[i] = true;
    }

    isOn = new Array(numberOfFilters).fill(false);
    expect(passesAll()).toBeFalsy();
  });
});
