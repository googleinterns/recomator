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

import { RecommendationExtra } from "./data_model/recommendation_extra";

export interface IRecommendationsStoreState {
  recommendations: RecommendationExtra[];
  failedProjects: string[];
  recommendationsByName: Map<string, RecommendationExtra>;
  requestId: string;
  cancel: boolean; // if we want to cancel fetching
  progress: number | null; // % recommendations loaded, null if no fetching is happening
  centralStatusWatcherRunning: boolean;
}

export function recommendationsStoreStateFactory(): IRecommendationsStoreState {
  return {
    recommendations: [],
    failedProjects: [],
    recommendationsByName: new Map<string, RecommendationExtra>(),
    requestId: "null",
    cancel: false,
    progress: null,
    centralStatusWatcherRunning: false
  };
}
