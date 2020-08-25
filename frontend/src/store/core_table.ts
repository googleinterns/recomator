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
import { IRootStoreState } from "./root";
import { RecommendationExtra } from "./data_model/recommendation_extra";
import {
  descriptionFilterAccepted,
  projectFilterAccepted,
  resourceFilterAccepted,
  typeFilterAccepted,
  statusFilterAccepted,
  costFilterAccepted
} from "./utils/core_table_filter_utils";

export interface ICoreTableStoreState {
  resourceNameSearchText: string;
  descriptionSearchText: string;
  projectsSelected: string[];
  typesSelected: string[];
  statusesSelected: string[];
  costCategoriesSelected: string[];
  selected: RecommendationExtra[];
}

export function coreTableStoreStateFactory(): ICoreTableStoreState {
  return {
    resourceNameSearchText: "",
    descriptionSearchText: "",
    projectsSelected: [],
    typesSelected: [],
    statusesSelected: [],
    costCategoriesSelected: [],
    selected: []
  };
}

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

export function isRecommendationInResults(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
): boolean {
  return (
    projectFilterAccepted(tableState, recExtra) &&
    typeFilterAccepted(tableState, recExtra) &&
    statusFilterAccepted(tableState, recExtra) &&
    costFilterAccepted(tableState, recExtra) &&
    resourceFilterAccepted(tableState, recExtra) &&
    descriptionFilterAccepted(tableState, recExtra)
  );
}

export const costCategoriesNames = { costs: "Costs", gains: "Savings" };
