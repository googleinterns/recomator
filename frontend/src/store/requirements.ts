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
import { getBackendAddress } from "@/config";
import { authFetch } from "./auth";
import router from '@/router';

const BACKEND_ADDRESS: string = getBackendAddress();
const FETCH_PROGRESS_WAIT_TIME = 100; // (1/10)s
const APPLY_PROGRESS_WAIT_TIME = 10000; // 10s
const HTTP_OK_CODE = 200;

export interface IRequirementsStoreState {
  projects: ProjectRequirement[];
  progress: null | number;
  errorCode: undefined | number;
  errorMessage: undefined | string;
  requestId: string;
  display: boolean;
}

export function requirementsStoreStateFactory(): IRequirementsStoreState {
  return {
    projects: [],
    progress: null,
    errorCode: undefined,
    errorMessage: undefined,
    requestId: "",
    display: false,
  };
}

const mutations: MutationTree<IRequirementsStoreState> = {
  // only entry point for projects
  startFetch(state): void {
    state.progress = 0;
  },

  addRequirement(state, requirement): void {
    if (state.projects.filter(elt => elt.name === requirement.name).length !== 0) {
      throw "Duplicate requirement name";
    }

    state.projects.push(requirement);
  },

  resetRequirements(state): void {
    state.projects = [];
  },

  endFetch(state): void {
    state.progress = null;
  },

  setProgress(state, progress: number) {
    state.progress = progress;
  },

  setRequestId(state, id: string) {
    state.requestId = id;
  },

  setError(state, errorInfo: { errorCode: number; errorMessage: string }) {
    state.errorCode = errorInfo.errorCode;
    state.errorMessage = errorInfo.errorMessage;
  },
};

const actions: ActionTree<IRequirementsStoreState, IRootStoreState> = {
  // Makes requests to the middleware and adds obtained projects to the store
  async fetchRequirements(context): Promise<void> {
    // one fetch at a time only
    if (context.state.progress !== null) {
      return;
    }
    context.commit("resetRequirements");
    context.commit("startFetch");

    // First, select the projects (temporarily hard-coded)
    const response = await authFetch(`${BACKEND_ADDRESS}/requirements`, {
      body: JSON.stringify({
        projects: context.getters("projectsStore/selectedProjects", {
          root: true,
        }),
      }),
      method: "POST",
    });
    const responseCode = response.status;
    // 201 = Created (Success)
    if (responseCode !== 201) {
      context.commit("setError", {
        errorCode: responseCode,
        errorMessage: `selecting projects failed: ${response.statusText}`,
      });
      return;
    }

    // An id has just been assigned for our us/project selection combination
    // that we need to refer to in future requests
    context.commit("setRequestId", await response.text());

    // send /recommendations requests until data received
    let responseJson: any;
    for (;;) {
      const response = await authFetch(
        `${BACKEND_ADDRESS}/requirements?request_id=${context.state.requestId}`
      );
      const responseCode = response.status;

      if (responseCode !== HTTP_OK_CODE) {
        context.commit("setError", {
          errorCode: responseCode,
          errorMessage: `progress check failed: ${response.statusText}`,
        });

        context.commit("endFetching");
        return;
      }

      responseJson = await response.json();

      if (responseJson.requirements !== undefined) {
        break;
      }

      context.commit(
        "setProgress",
        Math.floor(
          (100 * responseJson.batchesProcessed) / responseJson.numberOfBatches
        )
      );

      await delay(FETCH_PROGRESS_WAIT_TIME);
    }

    const requirementList = responseJson.map((elt :any) => new ProjectRequirement(elt.project, elt.requiements));
    for (const requirement of requirementList) {
      context.commit("addRequirement", requirement);
    }

    context.commit("endFetching");
  },

  proceedToRecommendations(context) {
    router.push("recommendations");
  },
};

export function requirementStoreFactory(): Module<
  IRequirementsStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: requirementsStoreStateFactory(),
    mutations: mutations,
    actions: actions,
  };
}