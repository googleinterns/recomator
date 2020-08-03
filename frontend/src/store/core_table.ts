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

export interface ICoreTableStoreState {
  resourceNameSearchText: string;
  descriptionSearchText: string;
  projectsSelected: Array<string>;
  typesSelected: Array<string>;
  // status is the correct plural form, but this is clearer
  statusesSelected: Array<string>;
}

export function coreTableStoreStateFactory(): ICoreTableStoreState {
  return {
    resourceNameSearchText: "",
    descriptionSearchText: "",
    projectsSelected: [],
    typesSelected: [],
    statusesSelected: []
  };
}

const mutations: MutationTree<ICoreTableStoreState> = {
  setResourceNameSearchText(state, text: string): void {
    state.resourceNameSearchText = text;
  },
  setDescriptionSearchText(state, text: string): void {
    state.descriptionSearchText = text;
  },
  setProjectsSelected(state, projects: Array<string>): void {
    state.projectsSelected = projects;
  },
  setTypesSelected(state, types: Array<string>): void {
    state.typesSelected = types;
  },
  setStatusesSelected(state, statuses: Array<string>): void {
    state.statusesSelected = statuses;
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