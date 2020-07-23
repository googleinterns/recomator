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

import "@/store/model";
import { extractFromResource } from "./utility";

import { Module, VuexModule, Mutation, Action } from "vuex-module-decorators";

@Module
export default class extends VuexModule {
  recommendations: Record<string, Recommendation> = {};

  @Mutation
  ADD_RECOMMENDATION(recommendation: Recommendation) {
    // prevent overwriting
    console.assert(
      !(recommendation.name in this.recommendations),
      `recommendation name ${recommendation.name} is present in the store already`
    );
    this.recommendations[recommendation.name] = recommendation;
  }

  @Action
  getRecommendationResource(recommendationName: string): string {
    return this.recommendations[recommendationName].content.operationGroups[0]
      .operations[0].resource;
  }

  @Action
  getRecommendationInstance(recommendationName: string): string {
    const resource = this.recommendations[recommendationName].content
      .operationGroups[0].operations[0].resource;

    return extractFromResource("instances", resource);
  }

  @Action
  getRecommendationProject(recommendationName: string): string {
    const resource = this.recommendations[recommendationName].content
      .operationGroups[0].operations[0].resource;

    return extractFromResource("projects", resource);
  }

  @Action
  getRecomendationDescription(recommendationName: string): string {
    return this.recommendations[recommendationName].description;
  }

  /* "3.5$ per week"
  @Action
  getRecommendationCostString(recommendationName: string): string {
    const recommendation = this.recommendations[recommendationName];
  }
  */

  @Action
  getRecommendation(recommendationName: string): Recommendation {
    console.assert(
      recommendationName in this.recommendations,
      `recommendation name ${recommendationName} is not present in the store`
    );
    return this.recommendations[recommendationName];
  }

  @Action
  fetchRecommendations() {
    //TODO
  }
}
