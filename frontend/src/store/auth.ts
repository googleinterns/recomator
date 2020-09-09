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
import router from "../router/index";
import store from "../store/root";

export interface IAuthStoreState {
  idToken?: string;
}

export function authStoreStateFactory(): IAuthStoreState {
  return {
    idToken: undefined
  };
}

const mutations: MutationTree<IAuthStoreState> = {
  setIDToken(state, token: string) {
    state.idToken = token;
  }
};

export function authStoreFactory(): Module<IAuthStoreState, IRootStoreState> {
  return {
    namespaced: true,
    state: authStoreStateFactory(),
    mutations: mutations
  };
}

// adds the idToken to the request, redirects to google sign-in if auth failed
// assumes init, inits.headers are in Record<string, string> format
export async function authFetch(
  input: string,
  init?: RequestInit
): Promise<Response> {
  // first, make sure that init, init.headers exist
  if (init == undefined) init = {};
  if (init.headers == undefined) init.headers = {};

  init.headers = Object.assign({}, init.headers as Record<string, string>, {
    Authorization: `Bearer ${store.state.authStore!.idToken}`
  });

  const response = await fetch(input, init);
  const responseCode = response.status;

  // access denied
  if (responseCode === 401) {
    // redirect to GoogleSignIn
    router.push({ name: "GoogleSignIn" });
  }

  return response;
}
