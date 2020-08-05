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
      store.state.recommendationsStore!.recommendations.includes(
        sampleRecommendation
      )
    ).toBe(true);
    // TODO add recommendation content checks as well
  });
});

describe("Recommendation-type objects", () => {
  test("Getting the project that the recommendation references", () => {
    expect(Model.getRecommendationProject(sampleRecommendation)).toEqual(
      "rightsizer-test"
    );
  });
});

describe("Fetching recommendations", () => {
  test("Getting the instance that the recommendation references",  () => {
    expect(
      Model.getRecommendationResourceShortName(sampleRecommendation)
    ).toEqual("alicja-test");
  });

  test("Fetching works correctly when given response without errors", async () => {
    jest.setTimeout(30000);
    fetchMock.doMock();

    const responses = [];
    responses.push(
      JSON.stringify({ batchesProcessed: 12, numberOfBatches: 100 })
    );
    responses.push(
      JSON.stringify({ batchesProcessed: 40, numberOfBatches: 100 })
    );
    responses.push(
      JSON.stringify({ batchesProcessed: 98, numberOfBatches: 100 })
    );
    responses.push(JSON.stringify({ recommendations: [sampleRecommendation] }));
    fetchMock.mockResponses(...responses);
    const store = rootStoreFactory();
    store.commit("recommendationsStore/resetRecommendations");
    await store.dispatch("recommendationsStore/fetchRecommendations");

    expect(store.state.recommendationsStore!.recommendations[0]).toEqual(
      sampleRecommendation
    );
    expect(store.state.recommendationsStore!.recommendations.length).toEqual(1);
    fetchMock.dontMock();
  });

  test("Fetching works correctly when given response with errors", async () => {
    jest.setTimeout(30000);
    fetchMock.doMock();

    const responses = [];
    responses.push(
      JSON.stringify({ batchesProcessed: 12, numberOfBatches: 100 })
    );
    responses.push(
      JSON.stringify({ batchesProcessed: 40, numberOfBatches: 100 })
    );
    responses.push(
      JSON.stringify({ batchesProcessed: 98, numberOfBatches: 100 })
    );
    responses.push(async () => {
      return {
        status: 302,
        body: JSON.stringify({ errorMessage: "Something failed" })
      };
    });
    responses.push(JSON.stringify({ recommendations: [sampleRecommendation] }));
    fetchMock.mockResponses(...responses);
    const store = rootStoreFactory();
    store.commit("recommendationsStore/resetRecommendations");
    await store.dispatch("recommendationsStore/fetchRecommendations");

    expect(store.state.recommendationsStore!.errorCode).toEqual(302);
    expect(store.state.recommendationsStore!.errorMessage).toEqual(
      "Something failed"
    );
    expect(store.state.recommendationsStore!.recommendations).toEqual([]);
    fetchMock.dontMock();
  });
});
