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
import { Module, VuexModule, Mutation, Action } from "vuex-module-decorators";
import { delay, getServerAddress } from "./utils";

const SERVER_ADDRESS: string = getServerAddress();
const REQUEST_DELAY = 100;
const HTTP_OK_CODE = 200;

@Module
export default class extends VuexModule {
  recommendations: Record<string, Recommendation> = {};
  progress: number | null = null; // percent of recommendations loaded, if null, then loading recommendations is not in progress
  errorCode: number | undefined = undefined;
  errorMessage: string | undefined = undefined;

  @Mutation
  addRecommendation(recommendation: Recommendation) {
    // prevent overwriting
    console.assert(
      !(recommendation.name in this.recommendations),
      `recommendation name ${recommendation.name} is present in the store already`
    );
    this.recommendations[recommendation.name] = recommendation;
  }

  @Mutation
  endFetching() {
    this.progress = null;
  }

  @Mutation
  setProgress(progress: number) {
    this.progress = progress;
  }

  @Mutation
  setError(errorInfo: { errorCode: number; errorMessage: string }) {
    this.errorCode = errorInfo.errorCode;
    this.errorMessage = errorInfo.errorMessage;
  }

  @Mutation
  resetRecommendations() {
    this.recommendations = {};
  }

  @Action
  getRecommendation(recommendationName: string): Recommendation {
    console.assert(
      recommendationName in this.recommendations,
      `recommendation name ${recommendationName} is not present in the store`
    );
    return this.recommendations[recommendationName];
  }

  @Action
  async fetchRecommendations() {
    if (this.progress !== null) {
      return;
    }

    this.context.commit("setProgress", 0);

    let response;
    let responseJson;
    let responseCode;

    for (;;) {
      response = await fetch(`${SERVER_ADDRESS}/recommendations`);

      responseJson = await response.json();

      responseCode = response.status;

      if (responseCode !== HTTP_OK_CODE) {
        this.context.commit("setError", {
          errorCode: responseCode,
          errorMessage: responseJson.errorMessage
        });

        this.context.commit("endFetching");

        return;
      }

      if (responseJson.recommendations !== undefined) {
        break;
      }

      this.context.commit(
        "setProgress",
        Math.floor(
          (100 * responseJson.batchesProcessed) / responseJson.numberOfBatches
        )
      );

      await delay(REQUEST_DELAY);
    }

    for (const recommendation of responseJson.recommendations) {
      this.context.commit("addRecommendation", recommendation);
    }

    this.context.commit("endFetching");
  }
}
