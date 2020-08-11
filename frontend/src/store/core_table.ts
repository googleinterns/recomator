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
import { RecommendationExtra } from "./model";
import {
  projectFilterAccepted,
  typeFilterAccepted,
  statusFilterAccepted
} from "./core_table_filter_utils";

export interface ICoreTableStoreState {
  projectsSelected: string[];
  typesSelected: string[];
  // status is the correct plural form, but this is clearer
  statusesSelected: string[];
  selected: RecommendationExtra[];
}

export function coreTableStoreStateFactory(): ICoreTableStoreState {
  return {
    projectsSelected: [],
    typesSelected: [],
    statusesSelected: [],
    selected: []
  };
}

const mutations: MutationTree<ICoreTableStoreState> = {
  setProjectsSelected(state, projects: string[]): void {
    state.projectsSelected = projects;
  },
  setTypesSelected(state, types: string[]): void {
    state.typesSelected = types;
  },
  setStatusesSelected(state, statuses: string[]): void {
    state.statusesSelected = statuses;
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

export const CoreTableStore = coreTableStoreFactory();

export function isRecommendationInResults(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
): boolean {
  return (
    projectFilterAccepted(tableState, recExtra) &&
    typeFilterAccepted(tableState, recExtra) &&
    statusFilterAccepted(tableState, recExtra)
  );
}
