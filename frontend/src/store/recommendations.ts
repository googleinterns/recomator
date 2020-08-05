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

import {
  Recommendation,
  RecommendationExtra,
  getRecommendationProject,
  getRecommendationType
} from "@/store/model";
import { delay, getServerAddress } from "./utils";
import { Module, MutationTree, ActionTree, GetterTree } from "vuex";
import { IRootStoreState } from "./root";

const SERVER_ADDRESS: string = getServerAddress();
const REQUEST_DELAY = 100;
const HTTP_OK_CODE = 200;

export interface IRecommendationsStoreState {
  recommendations: Recommendation[];
  selected: RecommendationExtra[];

  errorCode: number | undefined;
  errorMessage: string | undefined;
  // % recommendations loaded, null if no fetching is happening
  progress: number | null;
}

export function recommendationsStoreStateFactory(): IRecommendationsStoreState {
  return {
    recommendations: [],
    selected: [],
    progress: null,
    errorCode: undefined,
    errorMessage: undefined
  };
}

const mutations: MutationTree<IRecommendationsStoreState> = {
  addRecommendation(state, recommendation: Recommendation): void {
    state.recommendations.push(recommendation);
  },
  setSelected(state, selected: RecommendationExtra[]) {
    state.selected = selected;
  },
  endFetching(state) {
    state.progress = null;
  },
  setProgress(state, progress: number) {
    state.progress = progress;
  },
  resetRecommendations(state) {
    state.recommendations = [];
  },
  setError(state, errorInfo: { errorCode: number; errorMessage: string }) {
    state.errorCode = errorInfo.errorCode;
    state.errorMessage = errorInfo.errorMessage;
  }
};

const actions: ActionTree<IRecommendationsStoreState, IRootStoreState> = {
  async fetchRecommendations(context): Promise<void> {
    if (context.state.progress !== null) {
      return;
    }

    context.commit("setProgress", 0);

    let response;
    let responseJson;
    let responseCode;

    for (;;) {
      response = await fetch(`${SERVER_ADDRESS}/recommendations`);
      responseJson = await response.json();
      responseCode = response.status;

      if (responseCode !== HTTP_OK_CODE) {
        context.commit("setError", {
          errorCode: responseCode,
          errorMessage: responseJson.errorMessage
        });

        context.commit("endFetching");
        return;
      }

      if (responseJson.recommendations !== undefined) {
        break;
      }

      context.commit(
        "setProgress",
        Math.floor(
          (100 * responseJson.batchesProcessed) / responseJson.numberOfBatches
        )
      );

      await delay(REQUEST_DELAY);
    }

    for (const recommendation of responseJson.recommendations) {
      context.commit("addRecommendation", recommendation);
    }

    context.commit("endFetching");
  }
};

const getters: GetterTree<IRecommendationsStoreState, IRootStoreState> = {
  allProjects(state): string[] {
    const projects = state.recommendations.map(r =>
      getRecommendationProject(r)
    );
    return Array.from(new Set(projects));
  },
  allTypes(state): string[] {
    const projects = state.recommendations.map(r => getRecommendationType(r));
    return Array.from(new Set(projects));
  }
};

export function recommendationStoreFactory(): Module<
  IRecommendationsStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: recommendationsStoreStateFactory(),
    mutations: mutations,
    actions: actions,
    getters: getters
  };
}

export const RecommendationsStore = recommendationStoreFactory();
