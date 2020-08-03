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

import CoreTable from "@/components/CoreTable.vue";
import { rootStoreFactory } from "@/store/root";
import sampleRecommendation from "./sample_recommendation";

describe("Core Table", () => {
  test("Filtering results by resource name", async () => {
    const fakeStore = rootStoreFactory();
    const newSampleRecommendation = JSON.parse(
      JSON.stringify(sampleRecommendation)
    );
    const filterPredicate = (CoreTable.prototype.constructor as any)
      .filterPredicate;

    // Check if we accept resoruce name: bob-vm0 when searching for bob
    fakeStore.commit("coreTableStore/setResourceNameSearchText", "bob");
    newSampleRecommendation.content.operationGroups[0].operations[0].resource =
      "//compute.googleapis.com/projects/search/zones/us-east1-b/instances/bob-vm0";
    fakeStore.commit(
      "recommendationsStore/addRecommendation",
      newSampleRecommendation
    );
    expect(
      filterPredicate(fakeStore.state.coreTableStore!, newSampleRecommendation)
    ).toBeTruthy();

    // Check if we accept resoruce name: alice-vm0 when searching for bob
    newSampleRecommendation.content.operationGroups[0].operations[0].resource =
      "//compute.googleapis.com/projects/search/zones/us-east1-b/instances/alice-vm0";
    expect(
      filterPredicate(fakeStore.state.coreTableStore!, newSampleRecommendation)
    ).toBeFalsy();

    // Check if we accept resoruce name: alice-vm0 when searching for alice
    //  - also checking if updating the search text works
    fakeStore.commit("coreTableStore/setResourceNameSearchText", "alice");
    expect(
      filterPredicate(fakeStore.state.coreTableStore!, newSampleRecommendation)
    ).toBeTruthy();
  });
});
