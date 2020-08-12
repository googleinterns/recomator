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

import {
  RecommendationRaw,
  RecommendationExtra,
  getRecommendationProject,
  getRecommendationType,
  getInternalStatusMapping
} from "@/store/model";
import { delay, getServerAddress } from "./utils";
import { Module, MutationTree, ActionTree, GetterTree } from "vuex";
import { IRootStoreState } from "./root";

// TODO: move all of this to config in some next PR
const SERVER_ADDRESS: string = getServerAddress();
const FETCH_PROGRESS_WAIT_TIME = 100;
//CHANGE FROM THE FUTURE: const APPLY_PROGRESS_WAIT_TIME = 200;
const HTTP_OK_CODE = 200;

export interface IRecommendationsStoreState {
  recommendations: RecommendationExtra[];
  recommendationsByName: Map<string, RecommendationExtra>;
  errorCode: number | undefined;
  errorMessage: string | undefined;
  // % recommendations loaded, null if no fetching is happening
  progress: number | null;
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
    //TODO: add watching recommendations that are in progress at the app launch
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
  async fetchRecommendations(context): Promise<void> {
    if (context.state.progress !== null) {
      return;
    }

    context.commit("setProgress", 0);

    let response;
    let responseJson;
    let responseCode;

    for (;;) {
      response = await fetch(`${SERVER_ADDRESS}/recommendations`);
      responseJson = await response.json();
      responseCode = response.status;

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

    context.commit("endFetching");
  }
  /* [[FUTURE CHANGES, TO BE REVIEWED IN FURTHER PR-S TO KEEP THIS ONE SMALLER]]
  // we need names, not references so that we can find them in the state fast
  applyGivenRecommendations(context, selectedNames: string[]): void {
    // if selected has duplicates
    if (new Set(selectedNames).size !== selectedNames.length)
      throw "Duplicates found among given recommendation names";

    const selectedRecs = selectedNames.map(name =>
      context.state.recommendationsByName.get(name)
    );

    // if name not found
    if (selectedRecs.includes(undefined))
      throw "Name given doesn't match an existing recommendation";

    // If we find out that preparing the requests takes too long and
    // blocks the UI, we can wrap dispatches in setTimeout(...,0)
    selectedRecs.forEach(rec =>
      context.dispatch("_applySingleRecommendation", rec)
    );
  },
  // should return nearly immediately
  async _applySingleRecommendation(
    context,
    rec: RecommendationExtra
  ): Promise<void> {
    context.commit("setRecommendationStatus", {
      recName: rec.name,
      newStatus: "CLAIMED"
    });

    const response = await fetch(
      `${SERVER_ADDRESS}/recommendations/apply?name=${rec.name}`,
      { method: "POST" }
    );

    if (response.status === HTTP_OK_CODE) context.dispatch("_watchStatus", rec);
    else {
      context.commit("setRecommendationStatus", {
        recName: rec.name,
        newStatus: "FAILED"
      });
      context.commit("setRecommendationError", {
        recName: rec.name,
        header: `HTTP ERROR(${response.status})`,
        desc: `Couldn't reach the Recomator API:\n${response.statusText}`
      });
    }
  },
  //  should return nearly immediately, assumes "CLAIMED" status
  async _watchStatus(context, rec: RecommendationExtra): Promise<void> {
    if (rec.statusCol !== getInternalStatusMapping("CLAIMED"))
      throw "You can only watch claimed recommendations";

    const response = await fetch(
      `${SERVER_ADDRESS}/recommendations/checkStatus?name=${rec.name}`
    );
    if (response.status === HTTP_OK_CODE) {
      const responseJson = (await response.json()) as ICheckStatusResponse;
      switch (responseJson.status) {
        case "IN PROGRESS":
          // continue following the progress
          break;

        case "SUCCEEDED":
          context.commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "SUCCEEDED"
          });
          return;

        // Now we know it failed, tell the user why

        case "NOT APPLIED":
          context.commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED"
          });
          context.commit("setRecommendationError", {
            recName: rec.name,
            header: "Server hasn't acknowledged the request",
            desc:
              "Recomator API has not received the request to apply this recommendation. You can try sending it again."
          });
          return;

        case "FAILED":
          context.commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED"
          });
          context.commit("setRecommendationError", {
            recName: rec.name,
            header: "Applying recommendation failed server-side",
            desc: `${responseJson.errorMessage}`
          });
          return;

        default:
          context.commit("setRecommendationStatus", {
            recName: rec.name,
            newStatus: "FAILED"
          });
          context.commit("setRecommendationError", {
            recName: rec.name,
            header: `Bad status(${responseJson.status})`,
            desc: "Recomator API status not recognized."
          });
          return;
      }
    } else {
      // Non-200 HTTP code
      context.commit("setRecommendationStatus", {
        recName: rec.name,
        newStatus: "FAILED"
      });
      rec.errorHeader = `Status query failed (HTTP:${response.status})`;
      rec.errorHeader =
        "Failed to reach the Recomator API, recommendation status is unknown. We will try again in a moment.";
    }

    // ask the event loop to check the status once again in a bit
    setTimeout(
      () => context.dispatch("_watchStatus", rec),
      APPLY_PROGRESS_WAIT_TIME
    );
  }
}; 


interface ICheckStatusResponse {
  status: string;
  errorMessage: string;
}

END OF FUTURE CHANGES */
};

const getters: GetterTree<IRecommendationsStoreState, IRootStoreState> = {
  allProjects(state): string[] {
    const projects = state.recommendations.map(r =>
      getRecommendationProject(r)
    );
    return Array.from(new Set(projects));
  },
  allTypes(state): string[] {
    const types = state.recommendations.map(r => getRecommendationType(r));
    return Array.from(new Set(types));
  },
  allStatuses(state): string[] {
    const statuses = state.recommendations.map(r =>
      getInternalStatusMapping(r.stateInfo.state)
    );
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

export const RecommendationsStore = recommendationStoreFactory();
