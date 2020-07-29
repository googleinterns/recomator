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

import * as Model from "@/store/model";
import { rootStoreFactory } from "@/store/root";
import sampleRecommendation from "./sample_recommendation";

describe("Store", () => {
  test("addRecommendation", () => {
    const store = rootStoreFactory();

    store.commit(
      "recommendationsStore/addRecommendation",
      sampleRecommendation
    );
    // For some very weird reason, using expect().toHaveProperty() only works if the name is short
    expect(
      sampleRecommendation.name in
        store.state.recommendationsStore!.recommendations
    ).toBe(true);
    // TODO add recommendation content checks as well
  });
});

describe("Recommendation-type objects", () => {
  test("Getting the project that the recommendation references", async () => {
    expect(Model.getRecommendationProject(sampleRecommendation)).toEqual(
      "rightsizer-test"
    );
  });

  test("Getting the instance that the recommendation references", async () => {
    expect(
      Model.getRecommendationResourceShortName(sampleRecommendation)
    ).toEqual("alicja-test");
  });
});
