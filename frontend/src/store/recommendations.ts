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
import { delay, getServerAddress } from "./utils/misc";
import { getInternalStatusMapping } from "@/store/data_model/status_map";
import { Module, MutationTree, ActionTree, GetterTree } from "vuex";
import { IRootStoreState } from "./root";

// TODO: move all of this to config in some next PR
const SERVER_ADDRESS: string = getServerAddress();
const FETCH_PROGRESS_WAIT_TIME = 100;
const APPLY_PROGRESS_WAIT_TIME = 200;
const HTTP_OK_CODE = 200;

export interface IRecommendationsStoreState {
  recommendations: RecommendationExtra[];
  recommendationsByName: Map<string, RecommendationExtra>;
  errorCode: number | undefined;
  errorMessage: string | undefined;
  progress: number | null; // % recommendations loaded, null if no fetching is happening
}

export function recommendationsStoreStateFactory(): IRecommendationsStoreState {
  return {
    recommendations: [],
    recommendationsByName: new Map<string, RecommendationExtra>(),
    progress: null,
    errorCode: undefined,
    errorMessage: undefined
  };
}

const mutations: MutationTree<IRecommendationsStoreState> = {
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
  }
};

const actions: ActionTree<IRecommendationsStoreState, IRootStoreState> = {
  // Makes requests to the middleware and adds obtained recommendations to the store
  // - it can only be run once in the app lifetime,
  //    as there is no way to stop status watchers atm
  async fetchRecommendations(context): Promise<void> {
    // one fetch at a time only
    if (context.state.progress !== null) {
      return;
    }
    context.commit("resetRecommendations");
    context.commit("setProgress", 0);

    // send /recommendations requests until data received
    let responseJson: any;
    for (;;) {
      const response = await fetch(`${SERVER_ADDRESS}/recommendations`);
      responseJson = await response.json();
      const responseCode = response.status;

      if (responseCode !== HTTP_OK_CODE) {
        context.commit("setError", {
          errorCode: responseCode,
          errorMessage: responseJson.errorMessage
        });

        context.commit("endFetching");
        return;
      }

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

    for (const recommendation of responseJson.recommendations) {
      context.commit("addRecommendation", recommendation);
    }

    // start watchers for all added recommendations already in progress
    for (const recommendation of context.state.recommendations) {
      if (recommendation.statusCol == getInternalStatusMapping("CLAIMED"))
        context.dispatch("watchStatus", recommendation);
    }

    context.commit("endFetching");
  },

  applyGivenRecommendations(
    { dispatch, state },
    selectedNames: string[]
  ): void {
    // if selected has duplicates
    if (new Set(selectedNames).size !== selectedNames.length)
      throw "Duplicates found among given recommendation names";

    const selectedRecs = selectedNames.map(name =>
      state.recommendationsByName.get(name)
    );

    // if name not found
    if (selectedRecs.includes(undefined))
      throw "Name given doesn't match an existing recommendation";

    // If we find out that preparing the requests takes too long and
    //  blocks the UI, we can wrap dispatches in setTimeout(...,0)
    selectedRecs.forEach(rec => dispatch("applySingleRecommendation", rec));
  },

  // should return nearly immediately
  async applySingleRecommendation(
    { commit, dispatch },
    rec: RecommendationExtra
  ): Promise<void> {
    commit("setRecommendationStatus", {
      recName: rec.name,
      newStatus: "CLAIMED"
    });

    const response = await fetch(
      `${SERVER_ADDRESS}/recommendations/apply?name=${rec.name}`,
      { method: "POST" }
    );

    // If server accepted the request, watch the status. Otherwise, save the error.
    if (response.status === HTTP_OK_CODE) dispatch("watchStatus", rec);
    else {
      commit("setRecommendationStatus", {
        recName: rec.name,
        newStatus: "FAILED"
      });
      commit("setRecommendationError", {
        recName: rec.name,
        header: `HTTP ERROR(${response.status})`,
        desc: `Couldn't reach the Recomator API:\n${response.statusText}`
      });
    }
  },
  // If there is a "CLAIMED" recommendation, its status will change once it
  //  has finished being applied. Therefore, we want to follow it and update the store
  //  (which will, in turn, automatically update the UI).
  async watchStatus({ dispatch }, rec: RecommendationExtra): Promise<void> {
    for (;;) {
      const shouldContinue = await dispatch("watchStatusOnce", rec);
      if (!shouldContinue) break;
      // ask the browser to do something else for a bit and then resume
      await delay(APPLY_PROGRESS_WAIT_TIME);
    }
  },
  // Should return nearly immediately,
  //  the promise encapsulates the answer to: do we want to continue watching this recommendation?
  async watchStatusOnce(
    { commit },
    rec: RecommendationExtra
  ): Promise<boolean> {
    const response = await fetch(
      `${SERVER_ADDRESS}/recommendations/checkStatus?name=${rec.name}`
    );

    if (response.status === HTTP_OK_CODE) {
      const responseJson = (await response.json()) as ICheckStatusResponse;
      switch (responseJson.status) {
        case "IN PROGRESS":
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "CLAIMED"
          });
          // continue following the progress
          return true;

        case "SUCCEEDED":
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "SUCCEEDED"
          });
          return false;

        // Now we know it failed, tell the user why and return: false

        case "NOT APPLIED":
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED"
          });
          commit("setRecommendationError", {
            recName: rec.name,
            header: "Server hasn't acknowledged the request",
            desc:
              "Recomator API has not received the request to apply this recommendation. You can try applying it again."
          });
          return false;

        case "FAILED":
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED"
          });
          commit("setRecommendationError", {
            recName: rec.name,
            header:
              "Applying recommendation failed server-side. You can try applying it again.",
            desc: `${responseJson.errorMessage}`
          });
          return false;

        default:
          commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED"
          });
          commit("setRecommendationError", {
            recName: rec.name,
            header: `Bad status(${responseJson.status})`,
            desc: "Recomator API status not recognized."
          });
          return false;
      }
    } else {
      // Non-200 HTTP code
      commit("setRecommendationStatus", {
        recName: rec.name,
        newStatus: "FAILED"
      });
      commit("setRecommendationError", {
        recName: rec.name,
        header: `Status query failed (HTTP:${response.status})`,
        desc:
          "Failed to reach the Recomator API, recommendation status is unknown. We will try again in a moment."
      });
    }
    return true; // Continue watching the status, maybe this is a temporary connection error
  }
};

interface ICheckStatusResponse {
  status: string;
  errorMessage: string;
}

const getters: GetterTree<IRecommendationsStoreState, IRootStoreState> = {
  // Used for calculating filter choices

  allProjects(state): string[] {
    const projects = state.recommendations.map(r => r.projectCol);
    return Array.from(new Set(projects));
  },
  allTypes(state): string[] {
    const types = state.recommendations.map(r => r.typeCol);
    return Array.from(new Set(types));
  },
  allStatuses(state): string[] {
    const statuses = state.recommendations.map(r => r.statusCol);
    return Array.from(new Set(statuses));
  }
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
    getters: getters
  };
}
