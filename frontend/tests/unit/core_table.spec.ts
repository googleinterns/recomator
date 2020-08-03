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

// Helpful for debugs: console.log(wrapper.html())
//
// await Vue.nextTick() is used to wait for Vue to actually update the components
//  as recommended here: https://vue-test-utils.vuejs.org/guides/#testing-asynchronous-behavior
//
// We might consider switching to local Vue instances/vuetifys
//  (there is a function for that in test-utils), as it is recommended
//  but there are issues with that, which have fixes described here:
//   https://vuetifyjs.com/en/getting-started/unit-testing/

describe("Core Table", () => {
  test("Filtering results", async () => {
    const fakeStore = rootStoreFactory();
    const newSampleRecommendation = JSON.parse(
      JSON.stringify(sampleRecommendation)
    );

    fakeStore.commit("coreTableStore/setResourceNameSearchText", "bob");
    newSampleRecommendation.name =
      "//compute.googleapis.com/projects/search/zones/us-east1-b/instances/bob-vm0";
    fakeStore.commit(
      "recommendationsStore/addRecommendation",
      newSampleRecommendation
    );
    expect(
      CoreTable.prototype.constructor.filterPredicate(
        fakeStore.state.coreTableStore!,
        newSampleRecommendation
      )
    );

    newSampleRecommendation.name =
      "//compute.googleapis.com/projects/search/zones/us-east1-b/instances/alice-vm0";
    expect(
      !CoreTable.prototype.constructor.filterPredicate(
        fakeStore.state.coreTableStore!,
        newSampleRecommendation
      )
    );

    fakeStore.commit("coreTableStore/setResourceNameSearchText", "alice");
    expect(
      CoreTable.prototype.constructor.filterPredicate(
        fakeStore.state.coreTableStore!,
        newSampleRecommendation
      )
    );
  });
});
