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

import { Module, MutationTree } from "vuex";
import { IRootStoreState } from "./root_state";
import { RecommendationExtra } from "./data_model/recommendation_extra";
import {
  ICoreTableStoreState,
  coreTableStoreStateFactory
} from "./core_table_state";

const mutations: MutationTree<ICoreTableStoreState> = {
  setResourceNameSearchText(state, text: string): void {
    state.resourceNameSearchText = text;
  },
  setDescriptionSearchText(state, text: string): void {
    state.descriptionSearchText = text;
  },
  setProjectsSelected(state, projects: string[]): void {
    state.projectsSelected = projects;
  },
  setTypesSelected(state, types: string[]): void {
    state.typesSelected = types;
  },
  setStatusesSelected(state, statuses: string[]): void {
    state.statusesSelected = statuses;
  },
  setCostCategoriesSelected(state, costCategories: string[]): void {
    state.costCategoriesSelected = costCategories;
  },
  setSelected(state, selected: RecommendationExtra[]) {
    state.selected = selected;
  },
  setCurrentlySelectable(state, currentlySelectable: RecommendationExtra[]) {
    state.currentlySelectable = currentlySelectable;
  },
  selectAllSelectable(state) {
    for (const row of state.currentlySelectable) {
      if (!state.selected.includes(row)) {
        state.selected.push(row);
      }
    }
  },
  unselectAllSelectable(state) {
    for (const row of state.currentlySelectable) {
      state.selected = state.selected.filter(item => item !== row);
    }
  },
  select(state, toSelect: RecommendationExtra) {
    if (!state.selected.includes(toSelect)) {
      state.selected.push(toSelect);
    }
  },
  unselect(state, toSelect: RecommendationExtra) {
    state.selected = state.selected.filter(item => item !== toSelect);
  }
};

export function coreTableStoreFactory(): Module<
  ICoreTableStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: coreTableStoreStateFactory(),
    mutations: mutations
  };
}
