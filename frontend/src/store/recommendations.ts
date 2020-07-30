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

import { Recommendation } from "@/store/model";
import { delay, getServerAddress } from "./utils";
import { ReferenceWrapper } from "./reference";
import { Module, MutationTree, ActionTree } from "vuex";
import { IRootStoreState } from "./root";

const SERVER_ADDRESS: string = getServerAddress();
const REQUEST_DELAY = 100;
const HTTP_OK_CODE = 200;

export interface IRecommendationsStoreState {
  recommendations: Record<string, Recommendation>;
  errorCode: number | undefined;
  errorMessage: string | undefined;
  // % recommendations loaded, null if no fetching is happening
  progress: number | null;
}

export function recommendationsStoreStateFactory(): IRecommendationsStoreState {
  return {
    recommendations: {},
    progress: null,
    errorCode: undefined,
    errorMessage: undefined
  };
}

const mutations: MutationTree<IRecommendationsStoreState> = {
  addRecommendation(state, recommendation: Recommendation): void {
    // prevent overwriting
    console.assert(
      !(recommendation.name in state.recommendations),
      `recommendation name ${recommendation.name} is present in the store already`
    );
    state.recommendations[recommendation.name] = recommendation;
  },
  tryStartFetching(state, result: ReferenceWrapper<boolean>) {
    // Only one fetching may be in progress at once
    if (state.progress !== null) {
      result.setValue(true);
      return;
    }

    result.setValue(false);
    state.progress = 0;
  },
  endFetching(state) {
    state.progress = null;
  },
  setProgress(state, progress: number) {
    state.progress = progress;
  },
  setError(state, errorCodeAndMessage: [number, string]) {
    state.errorCode = errorCodeAndMessage[0];
    state.errorMessage = errorCodeAndMessage[1];
  }
};

const actions: ActionTree<IRecommendationsStoreState, IRootStoreState> = {
  getRecommendation(context, recommendationName: string): Recommendation {
    console.assert(
      recommendationName in context.state.recommendations,
      `recommendation name ${recommendationName} is not present in the store`
    );
    return context.state.recommendations[recommendationName];
  },

  async fetchRecommendations(context): Promise<void> {
    const isAnotherFetchInProgress = new ReferenceWrapper(false);
    context.commit("recommendationsStore/tryStartFetching", isAnotherFetchInProgress);

    if (isAnotherFetchInProgress.getValue()) {
      return;
    }

    let response;
    let responseJson;
    let responseCode;

    for (;;) {
      response = await fetch(`${SERVER_ADDRESS}/recommendations`);
      responseJson = await response.json();
      responseCode = response.status;

      if (responseCode !== HTTP_OK_CODE) {
        context.commit("recommendationsStore/setError", responseCode, responseJson.errorMessage);

        context.commit("recommendationsStore/endFetching");

        return;
      }

      if (responseJson.recommendations !== undefined) {
        break;
      }

      context.commit(
        "recommendationsStore/setProgress",
        Math.floor(
          (100 * responseJson.batchesProcessed) / responseJson.numberOfBatches
        )
      );

      await delay(REQUEST_DELAY);
    }

    for (const recommendation of responseJson.recommendations) {
      context.commit("recommendationsStore/addRecommendation", recommendation);
    }

    context.commit("recommendationsStore/endFetching");
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
    actions: actions
  };
}

export const RecommendationsStore = recommendationStoreFactory();
