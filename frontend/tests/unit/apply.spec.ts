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
import { RecommendationExtra } from "@/store/data_model/recommendation_extra";
import { getInternalStatusMapping } from "@/store/data_model/status_map";
import { freshSampleRawRecommendation } from "./sample_recommendations";
import { rootStoreFactory } from "@/store/root";
import { recommendationStoreFactory } from "@/store/recommendations";

let store: any, context: any, commit: any;
let dispatch: jest.Mock;

let firstRec: RecommendationExtra, secondRec: RecommendationExtra;

beforeAll(() => {
  enableFetchMocks(); // fetchMock instance visible
  fetchMock.dontMock(); // fetches only mocked at request
});

beforeEach(() => {
  fetchMock.doMock();
  store = rootStoreFactory();
  const modState = store.state.recommendationsStore!;

  // dispatch is now spied on and does nothing
  dispatch = jest.fn();

  // commit is converted to module form
  commit = (name: any, payload: any) => {
    store.commit("recommendationsStore/" + name, payload);
  };

  context = {
    dispatch: dispatch,
    commit: commit,
    state: modState
  };

  const recsRaw = [
    freshSampleRawRecommendation(),
    freshSampleRawRecommendation()
  ];
  recsRaw[0].name = "a/b";
  recsRaw[1].name = "a/b/c";
  store.commit("recommendationsStore/addRecommendation", recsRaw[0]);
  store.commit("recommendationsStore/addRecommendation", recsRaw[1]);
  firstRec = modState.recommendations[0];
  secondRec = modState.recommendations[1];
});

afterEach(() => {
  fetchMock.dontMock();
  // clear the spies
  fetchMock.mockClear();
  dispatch.mockClear();
});

test("applyGivenRecommendations action", async () => {
  const applier = recommendationStoreFactory().actions![
    "applyGivenRecommendations"
  ] as any;

  // Note: rejects/resolves fails (good) if the Promise actually doesn't

  // duplicates
  await expect(applier(context, ["a/b", "a/b"])).rejects.not.toBeNull();
  expect(dispatch).toBeCalledTimes(0);

  // non-existent
  dispatch.mockReset();
  await expect(applier(context, ["a/b", "a/b/c/d"])).rejects.not.toBeNull();
  expect(dispatch).toBeCalledTimes(0);

  // succeeding
  dispatch.mockReset();
  await expect(applier(context, ["a/b/c", "a/b"])).resolves.not.toBeNull();
  expect(dispatch).toBeCalledTimes(2);
  expect(dispatch.mock.calls[0][0]).toBe("applySingleRecommendation");
  expect(dispatch.mock.calls[1][0]).toBe("applySingleRecommendation");
  const passedRecNames = dispatch.mock.calls.map(call => call[1].name);
  expect(passedRecNames).toContain("a/b/c");
  expect(passedRecNames).toContain("a/b");
});

describe("applySingleRecommendation action", () => {
  let applier: any;
  beforeAll(() => {
    applier = recommendationStoreFactory().actions![
      "applySingleRecommendation"
    ] as any;
  });

  test("success", async () => {
    fetchMock.mockResponseOnce("");
    await applier(context, firstRec);

    expect((fetch as any).mock.calls.length).toEqual(1);
    expect((fetch as any).mock.calls[0][0].indexOf(firstRec.name)).not.toBe(-1);
    expect(firstRec.needsStatusWatcher).toBeTruthy();
    expect(firstRec.statusCol).toEqual(getInternalStatusMapping("CLAIMED"));
  });

  test("failure", async () => {
    fetchMock.mockResponseOnce(async () => {
      return { status: 404, body: "Endpoint not found" };
    });
    await applier(context, firstRec);

    expect((fetch as any).mock.calls.length).toEqual(1);
    expect(firstRec.needsStatusWatcher).toBeFalsy();
    expect(firstRec.statusCol).toEqual(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader).not.toBeNull();
  });
});

describe("startCentralStatusWatcher action", () => {
  const action = recommendationStoreFactory().actions![
    "startCentralStatusWatcher"
  ] as any;

  test("only one instance at a time allowed", async () => {
    // make this fail on checkStatusOnceForAll
    dispatch = jest.fn(() => {
      throw "abcdef";
    });
    context.dispatch = dispatch;

    // make sure that both expects are called
    expect.assertions(2);

    // it should fail at the first dispatch the first time
    action(context).catch((error: any) => {
      expect(error.toString()).toBe("abcdef");
    });

    // it should fail earlier the second time
    action(context).catch((error: any) => {
      expect(error.toString()).not.toBe("abcdef");
    });
  });
});

describe("checkStatusOnceForAll action", () => {
  let action: any;
  beforeAll(() => {
    action = recommendationStoreFactory().actions![
      "checkStatusOnceForAll"
    ] as any;
  });

  test("checkStatusOnce: true => false", async () => {
    dispatch = jest
      .fn()
      .mockResolvedValueOnce(true)
      .mockResolvedValueOnce(false);
    context.dispatch = dispatch;

    firstRec.needsStatusWatcher = true;

    await action(context);
    expect(firstRec.needsStatusWatcher).toBeTruthy();
    await action(context);
    expect(firstRec.needsStatusWatcher).toBeFalsy();
  });

  test("ignores iff needsStatusWatcher: false", async () => {
    firstRec.needsStatusWatcher = false;
    secondRec.needsStatusWatcher = true;

    dispatch = jest.fn().mockResolvedValue(false);
    context.dispatch = dispatch;

    await action(context);
    expect(dispatch.mock.calls).toEqual([["checkStatusOnce", secondRec]]);
  });
});

describe("checkStatusOnce action", () => {
  let checkStatusOnce: any;
  beforeAll(() => {
    checkStatusOnce = recommendationStoreFactory().actions![
      "checkStatusOnce"
    ] as any;
  });
  beforeEach(() => {
    jest.useFakeTimers(); // mocked setTimeout
  });

  // applied by us/someone else
  test("-> in progress/claimed", async () => {
    for (const status of ["CLAIMED", "IN PROGRESS"]) {
      fetchMock.mockResponseOnce(JSON.stringify({ status: status }));
      const shouldContinue = await checkStatusOnce(context, firstRec);

      expect(firstRec.statusCol).toBe(getInternalStatusMapping("CLAIMED"));
      expect(shouldContinue).toBeTruthy();
    }
  });

  test("-> succeeded", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "SUCCEEDED" }));
    const shouldContinue = await checkStatusOnce(context, firstRec);

    expect((fetch as any).mock.calls[0][0].indexOf(firstRec.name)).not.toBe(-1);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("SUCCEEDED"));
    expect(shouldContinue).toBeFalsy();
  });

  test("-> active", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "ACTIVE" }));
    const shouldContinue = await checkStatusOnce(context, firstRec);

    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader!.startsWith("Server has")).toBeTruthy();
    expect(shouldContinue).toBeFalsy();
  });

  test("-> failed", async () => {
    fetchMock.mockResponseOnce(
      JSON.stringify({
        status: "FAILED",
        errorMessage: "something bad happened"
      })
    );
    const shouldContinue = await checkStatusOnce(context, firstRec);

    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader!.startsWith("Applying ")).toBeTruthy();
    expect(firstRec.errorDescription).toBe("something bad happened");
    expect(shouldContinue).toBeFalsy();
  });

  test("-> gibberish", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "%%%" }));
    const shouldContinue = await checkStatusOnce(context, firstRec);

    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader!.startsWith("Bad status(")).toBeTruthy();
    expect(shouldContinue).toBeFalsy();
  });

  test("server connection error", async () => {
    fetchMock.mockResponseOnce(async () => {
      return { status: 404, body: "Endpoint not found" };
    });
    const shouldContinue = await checkStatusOnce(context, firstRec);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader?.indexOf("404")).not.toBe(-1);
    expect(shouldContinue).toBeTruthy(); // try again in this case
  });
});
