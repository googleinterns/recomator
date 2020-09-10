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

import { RecommendationRaw } from "@/store/data_model/recommendation_raw";
import { RecommendationExtra } from "@/store/data_model/recommendation_extra";
import { freshSampleRawRecommendation } from "./sample_recommendation";
import {
  RatioCounter,
  getOrSetDefault,
  TrainingData,
  similarityAvg,
  similaritySort,
  TrainingDataHandler
} from "@/store/smart_sort/similarity";

test("getOrSetDefault", () => {
  const map: Record<string, RatioCounter> = {};
  expect(getOrSetDefault(map, "foo")).toEqual([0, 0]);
  getOrSetDefault(map, "foo2");
  getOrSetDefault(map, "foo")[1] = 7;

  expect(getOrSetDefault(map, "foo2")).toEqual([0, 0]);
  expect(getOrSetDefault(map, "foo")).toEqual([0, 7]);
  expect(getOrSetDefault(map, "foo3")).toEqual([0, 0]);
});

const mockTrainingData = new TrainingData();
mockTrainingData.projectCounters = { projA: [30, 100], projB: [10, 20] };
mockTrainingData.typeCounters = { typeA: [15, 95], typeB: [25, 25] };

function setProjectAndType(
  recRaw: RecommendationRaw,
  proj: string,
  type: string
) {
  recRaw.content.operationGroups[0].operations[0].resource = `//compute.googleapis.com/projects/${proj}/zones/us-east1-b/instances/alicja-test`;
  recRaw.recommenderSubtype = type;
}

describe("similarityAvg", () => {
  let recRaw: RecommendationRaw;
  beforeEach(() => {
    recRaw = freshSampleRawRecommendation();
  });

  test("for unseen project and type", () => {
    const rec = new RecommendationExtra(recRaw);
    expect(similarityAvg(mockTrainingData, 45, rec)).toEqual(2);
  });

  test("for projB and typeA", () => {
    setProjectAndType(recRaw, "projB", "typeA");
    const rec = new RecommendationExtra(recRaw);
    expect(similarityAvg(mockTrainingData, 40, rec)).toBeCloseTo(
      2 - 10 / 40 - 15 / 40
    );
  });
});

test("similaritySort", () => {
  const recsRaw = [
    freshSampleRawRecommendation(),
    freshSampleRawRecommendation(),
    freshSampleRawRecommendation()
  ];
  recsRaw[0].description = "first";
  setProjectAndType(recsRaw[0], "projA", "typeA"); // 2-30/40-15/40
  recsRaw[1].description = "second";
  setProjectAndType(recsRaw[1], "projB", "typeA"); // 2-10/40-15/40
  recsRaw[2].description = "third";
  setProjectAndType(recsRaw[2], "projB", "typeB"); // 2-10/40-25/40

  const recs = recsRaw.map(rec => new RecommendationExtra(rec));

  similaritySort(recs, mockTrainingData);

  expect(recs.map(rec => rec.description)).toEqual([
    "first",
    "third",
    "second"
  ]);
});

describe("TrainingDataHandler add/applyAdded", () => {
  let handler: TrainingDataHandler;
  let rec: RecommendationExtra;
  beforeEach(() => {
    handler = new TrainingDataHandler();

    // assing a deep copy
    handler.data = JSON.parse(JSON.stringify(mockTrainingData));

    const recRaw = freshSampleRawRecommendation();
    setProjectAndType(recRaw, "projA", "typeB");
    rec = new RecommendationExtra(recRaw);
  });
  test("add recommendation to training data", () => {
    handler.addRecommendation(rec);
    expect(handler.data.projectCounters[rec.projectCol]).toEqual([30, 101]);
    expect(handler.data.typeCounters[rec.typeCol]).toEqual([25, 26]);
  });
  test("mark added recommendation as applied", () => {
    handler.addRecommendation(rec);
    handler.applyAddedRecommendation(rec);
    expect(handler.data.projectCounters[rec.projectCol]).toEqual([31, 101]);
    expect(handler.data.typeCounters[rec.typeCol]).toEqual([26, 26]);
  });
});

test("TrainingDataHandler save and load", () => {
  const handler = new TrainingDataHandler();
  const trainingData = JSON.parse(
    JSON.stringify(mockTrainingData)
  ) as TrainingData;
  handler.data = trainingData;

  handler.saveToLocalStorage();

  // now make sure that it actually saved the data somewhere
  handler.data.projectCounters["projX"] = [7, 7];

  handler.loadFromLocalStorage();

  expect(handler.data).toEqual(mockTrainingData);
});
