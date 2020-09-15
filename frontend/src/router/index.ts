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
import Requirements from "../views/Requirements.vue";
import Recommendations from "../views/Recommendations.vue";
import store from "../store/root";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Home",
    component: Home,
    beforeEnter(_, __, next) {
      // Asynchronously request and receive projects from the middleware
      store.dispatch("projectsStore/fetchProjects");
      next();
    }
  },

  {
    path: "/requirements",
    name: "Requirements",
    component: Requirements,
    beforeEnter(_, __, next) {
      // Asynchronously request and receive requirements from the middleware
      store.dispatch("requirementsStore/fetchRequirements", store.state.projectsStore?.projectsSelected);
      next();
    }
  },

  {
    path: "/recommendations",
    name: "Recommendations",
    component: Recommendations,
    beforeEnter(_, __, next) {
      // Asynchronously request and receive recommendations from the middleware
      store.dispatch("recommendationsStore/fetchRecommendations");// , store.state.projectsStore?.projectsSelected);
      next();
    }
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
