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

jest.mock("@/store/auth/auth_fetch");

import { enableFetchMocks } from "jest-fetch-mock";
import { RecommendationExtra } from "@/store/data_model/recommendation_extra";
import { getInternalStatusMapping } from "@/store/data_model/status_map";
import { freshSampleRawRecommendation } from "./sample_recommendation";
import { rootStoreFactory } from "@/store/root_store";
import { recommendationStoreFactory } from "@/store/recommendations_store";
import { getAuthFetch } from "@/store/auth/auth_fetch";

// Let's just act as if we were using a normal fetch (which will be mocked as well)
// that doesn't add the authentication token or redirect if auth failed.
// There is already a library for mocking fetch, so that is why we are not mocking
// authFetch as a whole.
(getAuthFetch as any).mockImplementation((_: any, maxRetries = 0) => {
  return async function authFetch(
    input: string,
    init?: RequestInit
  ): Promise<Response> {
    for (let i = 0; i < maxRetries + 1; i++) {
      const response = await fetch(input, init);

      if (response.ok) return response;
    }
    throw Error("redirected to error page");
  };
});

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

  // make sure that authFetch will not get an undefined token
  store.state.authStore.idToken = "idToken123";

  context = {
    dispatch: dispatch,
    commit: commit,
    state: modState,
    rootState: store.state
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
  (fetch as any).resetMocks();
  dispatch.mockClear();
});

test("applyGivenRecommendations action", async () => {
  const applier = recommendationStoreFactory().actions![
    "applyGivenRecommendations"
  ] as any;

  // Note: rejects/resolves fail if the Promise is resolved/rejected respectively

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
    fetchMock.mockResponseOnce(async () => {
      return { status: 201, body: "All good" };
    });
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
    fetchMock.mockResponseOnce(async () => {
      return { status: 200, body: "Unreachable success" };
    });

    firstRec.needsStatusWatcher = false;

    // Note: rejects/resolves fail if the Promise is resolved/rejected respectively
    await expect(applier(context, firstRec)).rejects.not.toBeNull();

    // 1 + 0 retries
    expect((fetch as any).mock.calls.length).toEqual(1);
    expect(firstRec.needsStatusWatcher).toBeFalsy();
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

  // expect(*).resolves/rejects throws if not resolves/rejects

  test("3 405 errors + success (1 + 3 retries)", async () => {
    for (let i = 0; i < 3; i++) {
      fetchMock.mockResponseOnce(async () => {
        return { status: 405, body: "Endpoint not found" };
      });
    }
    fetchMock.mockResponseOnce(JSON.stringify({ status: "SUCCEEDED" }));
    expect.assertions(2);
    await expect(checkStatusOnce(context, firstRec))
      .resolves.toBeFalsy()
      .then(() => {
        expect(firstRec.statusCol).toBe(getInternalStatusMapping("SUCCEEDED"));
      });
  });

  test("4 405 errors (1 + 3 retries)", async () => {
    for (let i = 0; i < 3; i++) {
      fetchMock.mockResponseOnce(async () => {
        return { status: 405, body: "Endpoint not found" };
      });
    }
    expect(checkStatusOnce(context, firstRec)).rejects.not.toBeNull();
  });
});
