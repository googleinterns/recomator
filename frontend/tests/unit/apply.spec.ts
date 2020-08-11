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
import { RecommendationExtra, getInternalStatusMapping } from "@/store/model";
import { freshSampleRawRecommendation } from "./sample_recommendation";
import { rootStoreFactory } from "@/store/root";
import { recommendationStoreFactory } from "@/store/recommendations";

let store: any, context: any, commit: any;
let dispatch: jest.Mock;

let firstRec: RecommendationExtra;

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
});

afterEach(() => {
  fetchMock.dontMock();
  // clear the spies
  fetchMock.mockClear();
  dispatch.mockClear();
});

test("applyGivenRecommendations action", () => {
  const applier = recommendationStoreFactory().actions![
    "applyGivenRecommendations"
  ] as any;
  expect(applier).not.toBeNull();

  // duplicates
  expect(() => {
    applier(context, ["a/b", "a/b"]);
  }).toThrowError();
  expect(dispatch).toBeCalledTimes(0);

  // non-existent
  dispatch.mockReset();
  expect(() => {
    applier(context, ["a/b", "a/b/c/d"]);
  }).toThrowError();
  expect(dispatch).toBeCalledTimes(0);

  dispatch.mockReset();
  expect(() => {
    applier(context, ["a/b/c", "a/b"]);
  }).not.toThrowError();
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
    expect(dispatch.mock.calls).toEqual([["watchStatus", firstRec]]);
    expect(firstRec.statusCol).toEqual(getInternalStatusMapping("CLAIMED"));
  });

  test("failure", async () => {
    fetchMock.mockResponseOnce(async () => {
      return { status: 404, body: "Endpoint not found" };
    });
    await applier(context, firstRec);

    expect((fetch as any).mock.calls.length).toEqual(1);
    expect(dispatch).toHaveBeenCalledTimes(0);
    expect(firstRec.statusCol).toEqual(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader).not.toBeNull();
  });
});

describe("watchStatus action", () => {
  let watchStatus: any;
  beforeAll(() => {
    watchStatus = recommendationStoreFactory().actions!["watchStatus"] as any;
  });
  beforeEach(() => {
    jest.useFakeTimers(); // mocked setTimeout
  });

  test("-> in progress", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "IN PROGRESS" }));
    await watchStatus(context, firstRec);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("CLAIMED"));

    // end the status updates cycle
    expect(setTimeout as any).toBeCalledTimes(1);
  });

  test("-> succeeded", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "SUCCEEDED" }));
    await watchStatus(context, firstRec);
    expect((fetch as any).mock.calls[0][0].indexOf(firstRec.name)).not.toBe(-1);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("SUCCEEDED"));

    // dispatch watchStatus in the future
    expect(setTimeout as any).toBeCalledTimes(0);
  });

  test("-> not applied", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "NOT APPLIED" }));
    await watchStatus(context, firstRec);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader!.startsWith("Server has")).toBeTruthy();

    expect(setTimeout as any).toBeCalledTimes(0);
  });

  test("-> failed", async () => {
    fetchMock.mockResponseOnce(
      JSON.stringify({
        status: "FAILED",
        errorMessage: "something bad happened"
      })
    );
    await watchStatus(context, firstRec);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader!.startsWith("Applying ")).toBeTruthy();
    expect(firstRec.errorDescription).toBe("something bad happened");

    expect(setTimeout as any).toBeCalledTimes(0);
  });

  test("-> gibberish", async () => {
    fetchMock.mockResponseOnce(JSON.stringify({ status: "%%%" }));
    await watchStatus(context, firstRec);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader!.startsWith("Bad status(")).toBeTruthy();

    expect(setTimeout as any).toBeCalledTimes(0);
  });

  test("server connection error", async () => {
    fetchMock.mockResponseOnce(async () => {
      return { status: 404, body: "Endpoint not found" };
    });
    await watchStatus(context, firstRec);
    expect(firstRec.statusCol).toBe(getInternalStatusMapping("FAILED"));
    expect(firstRec.errorHeader?.indexOf("404")).not.toBe(-1);

    // we want to try again in this case
    expect(setTimeout as any).toBeCalledTimes(1);
  });
});
