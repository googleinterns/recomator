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

import Vue from "vue";
import Vuex, { StoreOptions, Store, GetterTree } from "vuex";

import { IRootStoreState } from "./root_state";

import { recommendationStoreFactory } from "./recommendations_store";
import { requirementStoreFactory } from "./requirements_store";
import { coreTableStoreFactory } from "./core_table_store";
import { authStoreFactory } from "./auth_store";

import { isRecommendationInResults } from "./core_table_filters/aggregate";
import { RecommendationExtra } from "./data_model/recommendation_extra";
import { projectStoreFactory } from "./projects_store";

Vue.use(Vuex);

const getters: GetterTree<IRootStoreState, IRootStoreState> = {
  filteredRecommendationsWithExtras(state): RecommendationExtra[] {
    return state.recommendationsStore!.recommendations.filter(
      (recExtra: RecommendationExtra) =>
        isRecommendationInResults(state.coreTableStore!, recExtra)
    );
  },

  selectedProjects(state): string[] {
    return state.projectsStore!.projectsSelected.map(elt => elt.name);
  },

  failedProjects(state): string[] {
    return state.recommendationsStore!.failedProjects.map(elt => elt);
  }
};

export function rootStoreFactory(): Store<IRootStoreState> {
  const storeOptions: StoreOptions<IRootStoreState> = {
    state: {},
    getters: getters,
    modules: {
      recommendationsStore: recommendationStoreFactory(),
      requirementsStore: requirementStoreFactory(),
      projectsStore: projectStoreFactory(),
      coreTableStore: coreTableStoreFactory(),
      authStore: authStoreFactory()
    }
  };
  return new Store<IRootStoreState>(storeOptions);
}

export default rootStoreFactory();
