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
import { IRootStoreState } from "./root";
import router from '../router/index';
import { authFetch } from './auth';
import { getBackendAddress } from '@/config';
import { start } from 'repl';

const BACKEND_ADDRESS: string = getBackendAddress();
const HTTP_OK_CODE = 200;

export interface IProjectsStoreState {
  projects: Project[];
  projectsSelected: Project[];
  loading: boolean;
  loaded: boolean;
}

export function projectsStoreStateFactory(): IProjectsStoreState {
  return {
    projects: [],
    projectsSelected: [],
    loading: false,
    loaded: false,
  };
}

const getters: GetterTree<IProjectsStoreState, IRootStoreState> = {
  selectedProjects(state): string[] {
    return state.projectsSelected.map(elt => elt.name);
  },
};

const mutations: MutationTree<IProjectsStoreState> = {
  // only entry point for projects
  addProject(state, project: Project): void {
    if (state.projects.filter(elt => elt.name === project.name).length !== 0) {
      throw "Duplicate project name";
    }

    state.projects.push(project);
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

  // First, select the projects (temporarily hard-coded)
  const response = await authFetch(`${BACKEND_ADDRESS}/projects`);
  const responseCode = response.status;

  if (responseCode !== HTTP_OK_CODE) {
    context.commit("setError", {
      errorCode: responseCode,
      errorMessage: `getting projects failed: ${response.statusText}`
    });
    return;
  }
    
  const responseJSON = await response.json();


    for (const project of responseJSON.projects) {
      context.commit("addProject", project);
    }

  context.commit("endFetching");
},
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
