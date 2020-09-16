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

import { enableFetchMocks } from "jest-fetch-mock";
enableFetchMocks(); // making it possible to mock fetch in this test suite
fetchMock.dontMock(); // not mocking fetch in every test by default

import {
  getRecommendationProject,
  getRecommendationResourceShortName,
  getRecommendationZone,
  getResourceConsoleLink,
  getRecommendationCostPerWeek
} from "@/store/data_model/recommendation_raw";
import { RecommendationExtra } from "@/store/data_model/recommendation_extra";
import { rootStoreFactory } from "@/store/root";
import {
  freshSampleRawRecommendation,
  freshSampleSnapshotRawRecommendation,
  freshSampleStopVMRawRecommendation,
  freshSampleDeleteDiskRawRecommendation,
  freshPerformanceRawRecommendation
} from "./sample_recommendation";

describe("Store", () => {
  test("addRecommendation", () => {
    const store = rootStoreFactory();

    const sampleRecommendation = freshSampleRawRecommendation();
    store.commit(
      "recommendationsStore/addRecommendation",
      sampleRecommendation
    );

    expect(store.state.recommendationsStore!.recommendations).toEqual([
      new RecommendationExtra(sampleRecommendation)
    ]);
  });
});

test("Getting the project that the recommendation references", () => {
  expect(getRecommendationProject(freshSampleRawRecommendation())).toEqual(
    "rightsizer-test"
  );
});

test("Getting the instance that the recommendation references", () => {
  expect(
    getRecommendationResourceShortName(freshSampleRawRecommendation())
  ).toEqual("alicja-test");
});

// zones are already tested by the tests for links, so they don't need to be as thorough
test("Getting the zone of the resource for CHANGE_MACHINE_TYPE", () => {
  expect(getRecommendationZone(freshSampleRawRecommendation())).toEqual(
    "us-east1-b"
  );
});

test("Getting the cost of the recommendation from primaryImpact", () => {
  expect(
    getRecommendationCostPerWeek(freshSampleRawRecommendation())
  ).toBeCloseTo(-17.0961);
});

test("Getting the cost of the recommendation from additionalImpact", () => {
  expect(
    getRecommendationCostPerWeek(freshPerformanceRawRecommendation())
  ).toBeCloseTo(16.9);
});

describe("Getting a Console link for the resource", () => {
  test("type: CHANGE_MACHINE_TYPE ", () => {
    expect(getResourceConsoleLink(freshSampleRawRecommendation())).toEqual(
      "https://console.cloud.google.com/compute/instancesDetail/zones/us-east1-b/instances/alicja-test?project=rightsizer-test"
    );
  });

  test("type: SNAPSHOT_AND_DELETE_DISK", () => {
    expect(
      getResourceConsoleLink(freshSampleSnapshotRawRecommendation())
    ).toEqual(
      "https://console.cloud.google.com/compute/disksDetail/zones/europe-west1-d/disks/vertical-scaling-krzysztofk-wordpress?project=rightsizer-test"
    );
  });

  test("type: DELETE_DISK", () => {
    expect(
      getResourceConsoleLink(freshSampleDeleteDiskRawRecommendation())
    ).toEqual(
      "https://console.cloud.google.com/compute/disksDetail/zones/us-central1-a/disks/stanislawm-test-1?project=rightsizer-test"
    );
  });

  test("type: STOP_VM", () => {
    expect(
      getResourceConsoleLink(freshSampleStopVMRawRecommendation())
    ).toEqual(
      "https://console.cloud.google.com/compute/instancesDetail/zones/us-central1-c/instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test"
    );
  });
});

describe("Fetching recommendations", () => {
  const responsesPrefix = [
    async () => {
      return {
        status: 201,
        body: "someRequestId123"
      };
    },
    JSON.stringify({ batchesProcessed: 12, numberOfBatches: 100 }),
    JSON.stringify({ batchesProcessed: 40, numberOfBatches: 100 }),
    JSON.stringify({ batchesProcessed: 98, numberOfBatches: 100 })
  ];
  test("Fetching works correctly when given response without errors", async () => {
    jest.setTimeout(30000);
    fetchMock.doMock();

    const responses = responsesPrefix.map(r => r);

    responses.push(
      JSON.stringify({ batchesProcessed: 12, numberOfBatches: 100 })
    );
    responses.push(
      JSON.stringify({ batchesProcessed: 40, numberOfBatches: 100 })
    );
    responses.push(
      JSON.stringify({ batchesProcessed: 98, numberOfBatches: 100 })
    );

    const sampleRecommendation = freshSampleRawRecommendation();
    responses.push(JSON.stringify({ recommendations: [sampleRecommendation] }));

    fetchMock.mockResponses(...responses);

    const store = rootStoreFactory();
    store.commit("recommendationsStore/resetRecommendations");
    await store.dispatch("recommendationsStore/fetchRecommendations");

    expect(store.state.recommendationsStore!.recommendations[0]).toEqual(
      new RecommendationExtra(sampleRecommendation)
    );
    expect(store.state.recommendationsStore!.recommendations.length).toEqual(1);
    fetchMock.dontMock();
  });

  test("Fetching works correctly when given response with errors", async () => {
    jest.setTimeout(30000);
    fetchMock.doMock();

    const responses = responsesPrefix.map(r => r);
    responses.push(async () => {
      return {
        status: 302,
        body: JSON.stringify({ errorMessage: "Something failed" })
      };
    });
    responses.push(
      JSON.stringify({ recommendations: [freshSampleRawRecommendation()] })
    );
    fetchMock.mockResponses(...responses);
    const store = rootStoreFactory();
    store.commit("recommendationsStore/resetRecommendations");
    await store.dispatch("recommendationsStore/fetchRecommendations");

    expect(store.state.recommendationsStore!.errorCode).toEqual(302);
    expect(store.state.recommendationsStore!.errorMessage).toEqual(
      "progress check failed: Found"
    );
    expect(store.state.recommendationsStore!.recommendations).toEqual([]);
    fetchMock.dontMock();
  });
});
