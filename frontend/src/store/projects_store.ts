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
import { Module, MutationTree, ActionTree } from "vuex";
import { IRootStoreState } from "./root_state";
import { getAuthFetch } from "./auth/auth_fetch";
import { getBackendAddress } from "@/config";
import {
  IProjectsStoreState,
  projectsStoreStateFactory
} from "./projects_state";

const BACKEND_ADDRESS: string = getBackendAddress();

const mutations: MutationTree<IProjectsStoreState> = {
  addProject(state, project: string): void {
    if (state.projects.filter(elt => elt.name === project).length !== 0) {
      throw "Duplicate recommendation name";
    }

    state.projects.push(new Project(project));
  },

  startFetch(state): void {
    state.loading = true;
  },

  endFetch(state): void {
    state.loaded = true;
  },

  setSelected(state, projects: Project[]): void {
    state.projectsSelected = projects;
  },

  resetProjects(state): void {
    state.projects = [];
    state.projectsSelected = [];
    state.loaded = false;
    state.loading = false;
  }
};

const actions: ActionTree<IProjectsStoreState, IRootStoreState> = {
  async fetchProjects(context) {
    // one fetch at a time only
    if (context.state.loading === true) {
      return;
    }
    context.commit("resetProjects");
    context.commit("startFetch");

    const authFetch = getAuthFetch(context.rootState, 3);

    // First, select the projects
    const response = await authFetch(`${BACKEND_ADDRESS}/projects`);
    const responseJSON = await response.json();

    if (responseJSON !== null) {
      for (const project of responseJSON.projects) {
        context.commit("addProject", project);
      }
    }
    context.commit("endFetch");
  },

  saveSelectedProjects(context) {
    window.localStorage.setItem(
      "project_list",
      JSON.stringify(context.state.projectsSelected)
    );
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
    actions: actions
  };
}
