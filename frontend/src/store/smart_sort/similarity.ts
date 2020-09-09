/* Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the License);
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an AS IS BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

import { RecommendationExtra } from "../data_model/recommendation_extra";

// Used to represent applied/total ratio for a certain category
export type RatioCounter = [number, number];
export function getOrSetDefault(
  map: Record<string, RatioCounter>,
  key: string
): RatioCounter {
  if (map[key] == undefined) map[key] = [0, 0];
  return map[key];
}

// The formula for similarity average is rewritten in an easier to compute form
// For example, for projects, we would normally take the average of
// (0 if same project, 1 if different) across all applied recommendations.
// Now, this is just (# applied and with a different project)/(# applied).
export function similarityAvg(
  data: TrainingData,
  appliedSize: number,
  rec: RecommendationExtra
): number {
  const appliedWithSameProject = getOrSetDefault(
    data.projectCounters,
    rec.projectCol
  )[0];
  const partForProject = 1 - appliedWithSameProject / appliedSize;

  const appliedWithSameType = getOrSetDefault(
    data.typeCounters,
    rec.typeCol
  )[0];
  const partForType = 1 - appliedWithSameType / appliedSize;

  return partForProject + partForType;
}

// Orders (in place) by average similarity to applied recommendations.
// The summary of seen recommendations is stored in localStorage in the browser.
// Similarity between two recommendations is defined to be the number of features
// that are different between them. For example, if they have different projects (1)
// but the same type (0), then the similarity is 1 + 0 = 1 between these two.
export function similaritySort(
  recommendations: RecommendationExtra[],
  data: TrainingData
): void {
  // count the total number of applied recommendations in training data so far
  const appliedSize = Object.values(data.projectCounters)
    .map((rc: RatioCounter) => rc[0])
    .reduce((a, b) => a + b);

  const compareFn = function(
    fir: RecommendationExtra,
    sec: RecommendationExtra
  ): number {
    const firSimilarityAvg = similarityAvg(data, appliedSize, fir);
    const secSimilarityAvg = similarityAvg(data, appliedSize, sec);

    if (firSimilarityAvg === secSimilarityAvg) return 0;
    else return firSimilarityAvg < secSimilarityAvg ? -1 : 1;
  };

  recommendations.sort(compareFn);
}

export class TrainingData {
  projectCounters: Record<string, RatioCounter> = {};
  typeCounters: Record<string, RatioCounter> = {};
}

export class TrainingDataHandler {
  data: TrainingData = new TrainingData();
  loadFromLocalStorage(): void {
    const loadedDataString = window.localStorage.getItem("training_data");
    this.data =
      loadedDataString == null
        ? new TrainingData()
        : (JSON.parse(loadedDataString) as TrainingData);
  }
  // should be called after every batch of updates
  saveToLocalStorage(): void {
    window.localStorage.setItem("training_data", JSON.stringify(this.data));
  }
  addRecommendation(rec: RecommendationExtra): void {
    getOrSetDefault(this.data.projectCounters, rec.projectCol)[1]++;
    getOrSetDefault(this.data.typeCounters, rec.typeCol)[1]++;
  }
  applyAddedRecommendation(rec: RecommendationExtra): void {
    getOrSetDefault(this.data.projectCounters, rec.projectCol)[0]++;
    getOrSetDefault(this.data.typeCounters, rec.typeCol)[0]++;
  }
}

export const trainingDataHandler = new TrainingDataHandler();
trainingDataHandler.loadFromLocalStorage();
