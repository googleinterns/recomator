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

import { ProjectRequirement } from "@/store/data_model/project_with_requirement";
import { delay } from "./utils/misc";
import { Module, MutationTree, ActionTree } from "vuex";
import { IRootStoreState } from "./root";
import { Requirement } from "./data_model/project_with_requirement";
import { Project } from "./data_model/project";

const FETCH_WAIT_TIME = 500; // (1/2)s
const REQUIREMENT_LIST = [
  new Requirement("Cloud Resource Manager API", true, "xxx"),
  new Requirement("Compute Engine API", false, "xxx"),
  new Requirement("Service Usage API", true, "xxx"),
  new Requirement("Recommender API", false, "xxx")
];

export interface IRequirementsStoreState {
  projects: ProjectRequirement[];
  projectsSelected: ProjectRequirement[];
  progress: null | number;
  errorCode: undefined | number;
  errorMessage: undefined | string;
  display: boolean;
}

export function requirementsStoreStateFactory(): IRequirementsStoreState {
  return {
    projects: [],
    projectsSelected: [],
    progress: null,
    errorCode: undefined,
    errorMessage: undefined,
    display: false
  };
}

const mutations: MutationTree<IRequirementsStoreState> = {
  // only entry point for projects
  addProjectRequirement(state, project: ProjectRequirement): void {
    if (state.projects.filter(elt => elt.name === project.name).length !== 0) {
      throw "Duplicate project name";
    }

    state.projects.push(project);
  },

  startFetch(state): void {
    state.progress = 0;
  },

  endFetch(state): void {
    state.progress = null;
  },

  setSelected(state, projects: ProjectRequirement[]): void {
    state.projectsSelected = projects;
  },
  setProgress(state, progress: number) {
    state.progress = progress;
  },

  setError(state, errorInfo: { errorCode: number; errorMessage: string }) {
    state.errorCode = errorInfo.errorCode;
    state.errorMessage = errorInfo.errorMessage;
  }
};

const actions: ActionTree<IRequirementsStoreState, IRootStoreState> = {
  // Makes requests to the middleware and adds obtained projects to the store
  async fetchRequirements(context, selectedProjects: Project[]): Promise<void> {
    context.commit("startFetch");
    let i = 0;

    for (const project of selectedProjects) {
      i++;
      await delay(2*FETCH_WAIT_TIME);
      context.commit("setProgress", i / selectedProjects.length * 100)
      context.commit(
        "addProjectRequirement",
        new ProjectRequirement(project.name, REQUIREMENT_LIST)
      );
    }

    context.commit("endFetch");
  },

  proceedToRecommendations(context, selectedProjects) {
    context.dispatch(
      "/recommendationsStore/fetchRecommendations",
      selectedProjects
    );
  }
};

export function requirementStoreFactory(): Module<
  IRequirementsStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: requirementsStoreStateFactory(),
    mutations: mutations,
    actions: actions
  };
}
