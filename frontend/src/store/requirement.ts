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

import { Project } from "@/store/data_model/project";
import { delay } from "./utils/misc";
import { Module, MutationTree, ActionTree } from "vuex";
import { IRootStoreState } from "./root";
import { Requirement } from "./data_model/project";

const FETCH_WAIT_TIME = 500; // (1/2)s
const REQUIREMENT_LIST = [
  new Requirement("Cloud Resource Manager API", true),
  new Requirement("Compute Engine API", false),
  new Requirement("Service Usage API", true),
  new Requirement("Recommender API", false)
];

export interface IProjectsStoreState {
  projects: Project[];
  projectsSelected: Project[];
  requirements: boolean;
  recommendations: boolean;
  loaded: boolean;
}

export function projectsStoreStateFactory(): IProjectsStoreState {
  return {
    projects: [],
    projectsSelected: [],
    requirements: false,
    recommendations: false,
    loaded: false
  };
}

const mutations: MutationTree<IProjectsStoreState> = {
  // only entry point for projects
  addProject(state, project: Project): void {
    if (state.projects.filter(elt => elt.name === project.name).length !== 0) {
      throw "Duplicate project name";
    }

    state.projects.push(project);
  },

  endFetch(state): void {
    state.loaded = true;
  },

  setSelected(state, projects: Project[]): void {
    state.projectsSelected = projects;
  }
};

const actions: ActionTree<IProjectsStoreState, IRootStoreState> = {
  // Makes requests to the middleware and adds obtained projects to the store
  async fetchProjects(context): Promise<void> {
    await delay(FETCH_WAIT_TIME);

    const ProjectCount = 10;
    for (let i = 0; i < ProjectCount; i++) {
      const projectName = `Project ${i}`

      context.commit(
        "addProject",
        new Project(projectName, [])
      );
    }
    
    context.commit("endFetch");
  },

  async fetchRequirements(
    context, selectedProjects  ): Promise<void> {
    await delay(FETCH_WAIT_TIME);
    context.commit("setRequirements");
  },

  proceedToRequirements(context, selectedProjects) {
    context.dispatch("/projectsStore/fetchRequirements", selectedProjects);
  },

  proceedToRecommendations(context, selectedProjects) {
    context.dispatch("/recommendationsStore/fetchRecommendations", selectedProjects);
  }

};

export function requirementStoreFactory(): Module<
  IProjectsStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: requirementStoreStateFactory(),
    mutations: mutations,
    actions: actions
  };
}
