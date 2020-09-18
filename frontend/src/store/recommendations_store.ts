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

import { RecommendationRaw } from "@/store/data_model/recommendation_raw";
import { RecommendationExtra } from "@/store/data_model/recommendation_extra";
import { delay } from "./utils/misc";
import { getInternalStatusMapping } from "@/store/data_model/status_map";
import { Module, MutationTree, ActionTree, GetterTree } from "vuex";
import { IRootStoreState } from "./root_state";
import { getBackendAddress } from "../config";
import { getAuthFetch } from "./auth/auth_fetch";
import { similaritySort, trainingDataHandler } from "./smart_sort/similarity";
import {
  IRecommendationsStoreState,
  recommendationsStoreStateFactory,
} from "./recommendations_state";

// TODO: move all of this to config in some next PR
const BACKEND_ADDRESS: string = getBackendAddress();
const FETCH_PROGRESS_WAIT_TIME = 100; // (1/10)s
const APPLY_PROGRESS_WAIT_TIME = 10000; // 10s
const HTTP_OK_CODE = 200;

const mutations: MutationTree<IRecommendationsStoreState> = {
  // only entry point for recommendations
  addRecommendation(state, recommendation: RecommendationRaw): void {
    const extended = new RecommendationExtra(recommendation);
    if (state.recommendationsByName.get(extended.name) !== undefined)
      throw "Duplicate recommendation name";
    state.recommendations.push(extended);
    state.recommendationsByName.set(extended.name, extended);
  },

  endFetching(state) {
    state.progress = null;
  },
  setProgress(state, progress: number) {
    state.progress = progress;
  },
  resetRecommendations(state) {
    state.recommendations = [];
    state.recommendationsByName.clear();
  },
  setError(state, errorInfo: { errorCode: number; errorMessage: string }) {
    state.errorCode = errorInfo.errorCode;
    state.errorMessage = errorInfo.errorMessage;
  },
  setRequestId(state, requestId: string) {
    state.requestId = requestId;
  },
  doSimilaritySort(state) {
    similaritySort(state.recommendations, trainingDataHandler.data);
  },
  setRecommendationStatus(
    state,
    payload: { recName: string; newStatus: string }
  ) {
    const rec = state.recommendationsByName.get(payload.recName);
    if (rec == undefined)
      throw `Attempting to access an inexistent recommendation`;
    rec.statusCol = getInternalStatusMapping(payload.newStatus);
  },
  setRecommendationError(
    state,
    payload: { recName: string; header: string; desc: string }
  ) {
    const rec = state.recommendationsByName.get(payload.recName);
    if (rec == undefined)
      throw `Attempting to access an inexistent recommendation`;
    rec.errorHeader = payload.header;
    rec.errorDescription = payload.desc;
  },
  setRecommendationNeedsStatusWatcher(
    state,
    payload: { recName: string; needs: boolean }
  ) {
    const rec = state.recommendationsByName.get(payload.recName);
    if (rec == undefined)
      throw `Attempting to access an inexistent recommendation`;
    rec.needsStatusWatcher = payload.needs;
  },
  registerCentralStatusWatcher(state) {
    if (state.centralStatusWatcherRunning)
      throw `More than once central status watcher`;
    state.centralStatusWatcherRunning = true;
  },
};

const actions: ActionTree<IRecommendationsStoreState, IRootStoreState> = {
  // Makes requests to the middleware and adds obtained recommendations to the store
  async fetchRecommendations(context): Promise<void> {
    // one fetch at a time only
    if (context.state.progress !== null) {
      return;
    }
    context.commit("resetRecommendations");
    context.commit("setProgress", 0);

    // First, select the projects
    const authFetch = getAuthFetch(context.rootState);
    const response = await authFetch(`${BACKEND_ADDRESS}/recommendations`, {
      body: JSON.stringify({
        projects: context.rootGetters["selectedProjects"],
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
      const authFetch = getAuthFetch(context.rootState);
      const response = await authFetch(
        `${BACKEND_ADDRESS}/recommendations?request_id=${context.state.requestId}`
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

      if (responseJson.recommendations !== undefined) {
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

    if (responseJson.recommendations !== null) {
      for (const recommendation of responseJson.recommendations)
        context.commit("addRecommendation", recommendation);
      for (const recommendation of context.state.recommendations)
        trainingDataHandler.addRecommendation(recommendation);

      context.commit("doSimilaritySort");
    }

    context.commit("endFetching");
  },

  async applyGivenRecommendations(
    { dispatch, state },
    selectedNames: string[]
  ): Promise<void> {
    // if selected has duplicates
    if (new Set(selectedNames).size !== selectedNames.length)
      throw "Duplicates found among given recommendation names";

    const selectedRecs = selectedNames.map((name) =>
      state.recommendationsByName.get(name)
    );

    // if name not found
    if (selectedRecs.includes(undefined))
      throw "Name given doesn't match an existing recommendation";

    // update training data; if a recommenadation fails and is re-applied,
    // it will be counted twice but we don't expect too many failures,
    // so this should have a negligible effect.
    for (const rec of selectedRecs)
      trainingDataHandler.applyAddedRecommendation(rec!);
    trainingDataHandler.saveToLocalStorage();

    // Apply all, waiting for them to send one by one
    for (const rec of selectedRecs)
      await dispatch("applySingleRecommendation", rec!);
  },

  // should return nearly immediately
  async applySingleRecommendation(
    { commit, rootState },
    rec: RecommendationExtra
  ): Promise<void> {
    commit("setRecommendationStatus", {
      recName: rec.name,
      newStatus: "CLAIMED",
    });

    const authFetch = getAuthFetch(rootState);
    const response = await authFetch(
      `${BACKEND_ADDRESS}/recommendations/apply?name=${rec.name}`,
      { method: "POST" }
    );

    // If server accepted the request, watch the status. Otherwise, save the error.
    if (response.status === 201)
      commit("setRecommendationNeedsStatusWatcher", {
        recName: rec.name,
        needs: true,
      });
    else {
      commit("setRecommendationNeedsStatusWatcher", {
        recName: rec.name,
        needs: false,
      });
      commit("setRecommendationStatus", {
        recName: rec.name,
        newStatus: "FAILED",
      });
      commit("setRecommendationError", {
        recName: rec.name,
        header: `HTTP ERROR(${response.status})`,
        desc: `Couldn't reach the Recomator API:\n${response.statusText}`,
      });
    }
  },
  // If there is a "CLAIMED" recommendation, its status will change once it
  // has finished being applied. Therefore, we want to follow it and update the store
  // (which will, in turn, automatically update the UI).
  //
  // It is safe to add to/clear the recommendations list during execution,
  // as we are operating on a copy
  async startCentralStatusWatcher({ commit, dispatch }): Promise<void> {
    commit("registerCentralStatusWatcher");
    for (;;) {
      await dispatch("checkStatusOnceForAll");
      // ask the browser to do something else for a bit and then resume
      await delay(APPLY_PROGRESS_WAIT_TIME);
    }
  },
  async checkStatusOnceForAll({ state, commit, dispatch }): Promise<void> {
    // make a copy first to make the array constant inside the for loop
    const recsCopy = state.recommendations.map((rec) => rec);

    // check status of all recommendations that need to be watched
    // (waits for a response before starting a new request)
    for (const rec of recsCopy) {
      if (!rec.needsStatusWatcher) continue;

      const shouldContinue = await dispatch("checkStatusOnce", rec);
      commit("setRecommendationNeedsStatusWatcher", {
        recName: rec.name,
        needs: shouldContinue,
      });
    }
  },
  // Iff awaited, waits for the status check to finish.
  // Assumes the recommendation has been applied already (by us/Pantheon/something else).
  // Returns: do we want to continue watching this recommendation?
  async checkStatusOnce(
    { commit, rootState },
    rec: RecommendationExtra
  ): Promise<boolean> {
    // send the request
    const authFetch = getAuthFetch(rootState);
    const response = await authFetch(
      `${BACKEND_ADDRESS}/recommendations/checkStatus?name=${rec.name}`
    );

    // handle all possible types of responses
    if (response.status === HTTP_OK_CODE) {
      const responseJson = (await response.json()) as ICheckStatusResponse;

      switch (responseJson.status) {
        case "CLAIMED": // applied somewhere else (e.g. Pantheon, us earlier
        case "IN PROGRESS": // applied by us
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "CLAIMED",
          });
          // continue following the progress
          return true;

        case "SUCCEEDED":
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "SUCCEEDED",
          });
          return false;

        // Now we know it failed, save the error message telling the user why

        case "ACTIVE": // shouldn't be ACTIVE if it has been applied
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED",
          });
          commit("setRecommendationError", {
            recName: rec.name,
            header: "Server hasn't acknowledged the request",
            desc:
              "The recommendation is still active, so the Recomator API has not received the request to apply this recommendation. You can try applying it again.",
          });
          return false;

        case "FAILED":
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED",
          });
          commit("setRecommendationError", {
            recName: rec.name,
            header:
              "Applying recommendation failed server-side. You can try applying it again.",
            desc: `${responseJson.errorMessage}`,
          });
          return false;

        default:
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED",
          });
          commit("setRecommendationError", {
            recName: rec.name,
            header: `Bad status(${responseJson.status})`,
            desc: "Recomator API status not recognized.",
          });
          return false;
      }
    } else {
      // Non-200 HTTP code
      commit("setRecommendationStatus", {
        recName: rec.name,
        newStatus: "FAILED",
      });
      commit("setRecommendationError", {
        recName: rec.name,
        header: `Status query failed (HTTP:${response.status})`,
        desc:
          "Failed to reach the Recomator API, recommendation status is unknown. We will try again in a moment.",
      });
    }
    return true; // Continue watching the status, maybe this is a temporary connection error
  },
};

interface ICheckStatusResponse {
  status: string;
  errorMessage: string;
}

const getters: GetterTree<IRecommendationsStoreState, IRootStoreState> = {
  // Used for calculating filter choices, sorted to avoid changes in order
  // caused by a different set of recommendations in a filtered view
  // (for example. the first type found might change very frequently)

  allProjects(state): string[] {
    const projects = state.recommendations.map((r) => r.projectCol).sort();
    return Array.from(new Set(projects));
  },
  allTypes(state): string[] {
    const types = state.recommendations.map((r) => r.typeCol).sort();
    return Array.from(new Set(types));
  },
  allStatuses(state): string[] {
    const statuses = state.recommendations.map((r) => r.statusCol).sort();
    return Array.from(new Set(statuses));
  },
};

export function recommendationStoreFactory(): Module<
  IRecommendationsStoreState,
  IRootStoreState
> {
  return {
    namespaced: true,
    state: recommendationsStoreStateFactory(),
    mutations: mutations,
    actions: actions,
    getters: getters,
  };
}
