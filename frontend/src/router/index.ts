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

import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";
import store from "../store/root";
import { IRootStoreState } from "../store/root";
import { getBackendAddress } from "../config";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Home",
    component: Home,
    beforeEnter(_, __, next) {
      const token = (store.state as IRootStoreState).authStore!.idToken;
      // redirect to google sign in if we don't have a token
      if (token == undefined) {
        next({ name: "GoogleSignIn" });
        return;
      }

      // Asynchronously request and receive recommendations from the middleware
      store.dispatch("recommendationsStore/fetchRecommendations");
      // Start status watchers
      store.dispatch("recommendationsStore/startCentralStatusWatcher");
      next();
    }
  },
  {
    path: "/auth",
    name: "AuthCodeReceiver",
    async beforeEnter(to, _, next) {
      const authCode = to.query.code;
      if (authCode == undefined) {
        next(new Error("The auth code is missing"));
        return;
      }

      // exchange the authCode for a token
      const response = await fetch(
        `${getBackendAddress()}/auth?code=${authCode}`
      );
      const responseCode = response.status;

      if (responseCode !== 200) {
        // we could redirect to googlesign-in instead, but
        // this might introduce an infinite loop
        next(new Error("Failed to connect to the backend"));
        return;
      }

      // the request was successful, extract the token
      const responseJson = await response.json();
      const token = responseJson.token;

      store.commit("authStore/setIDToken", token);
      next({ name: "Home" });
    }
  },
  {
    // internal endpoint to redirect to Google sign-in
    path: "/google-sign-in",
    name: "GoogleSignIn",
    beforeEnter() {
      window.location.href = `${getBackendAddress()}/redirect`;
    }
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
