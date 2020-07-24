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
import { delay } from "./utils";

const ServerAddress = "http://localhost:8082";
const RequestDelay = 100;

@Module
export default class extends VuexModule {
  recommendations: Record<string, Recommendation> = {};

  @Mutation
  addRecommendation(recommendation: Recommendation) {
    // prevent overwriting
    console.assert(
      !(recommendation.name in this.recommendations),
      `recommendation name ${recommendation.name} is present in the store already`
    );
    this.recommendations[recommendation.name] = recommendation;
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
    let response;
    let responseJson;

    for (;;) {
      response = await fetch(`${ServerAddress}/recommendations`);
      responseJson = await response.json();

      if (responseJson.recommendations !== undefined) {
        break;
      }

      await delay(RequestDelay);
    }

    for (const recommendation of responseJson.recommendations) {
      this.context.commit("addRecommendation", recommendation);
    }
  }
}
