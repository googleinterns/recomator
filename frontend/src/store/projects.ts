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
import { Module, MutationTree, ActionTree, GetterTree } from "vuex";
import { IRootStoreState } from './root';

const FETCH_WAIT_TIME = 500; // (1/2)s
const REQUIREMENT_LIST = [
  "Cloud Resource Manager API", 
  "Compute Engine API",
  "Service Usage API",
  "Recommender API",
]

export interface IProjectsStoreState {
  projects: Project[];
  projectsSelected: string[];
  chosen: boolean;
  loaded: boolean;
}

export function projectsStoreStateFactory(): IProjectsStoreState {
  return {
    projects: [],
    projectsSelected: [],
    chosen: false,
    loaded: false,
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

  setSelected(state, projects: string[]): void {
    state.projectsSelected = projects;
  },

  chooseProjects(state): void {
    state.chosen = true;
  }
};

const actions: ActionTree<IProjectsStoreState, IRootStoreState> = {
  // Makes requests to the middleware and adds obtained recommendations to the store
  async fetchProjects(context): Promise<void> {
    await delay(FETCH_WAIT_TIME);

    const ProjectCount = 10;
    for (let i = 0; i < ProjectCount; i++) {
      const projectName = `Project ${i}`;
      const projectRequirements = await context.dispatch("fetchRequirements", projectName);

      context.commit("addProject", new Project(projectName, projectRequirements));
    }

    context.commit("endFetch");
  },

  async fetchRequirements(context, projectName: string): Promise<string[]> {
    await delay(FETCH_WAIT_TIME);
    return [REQUIREMENT_LIST[parseInt(projectName[projectName.length - 1], 10)]];
  }
};

export function projectStoreFactory(): Module<
  IProjectsStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: projectsStoreStateFactory(),
    mutations: mutations,
    actions: actions,
  };
}
