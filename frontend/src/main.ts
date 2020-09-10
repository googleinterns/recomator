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
import App from "./App.vue";
import router from "./router";
import store from "./store/root";
import vuetify from "./plugins/vuetify";

Vue.config.productionTip = false;

/* In order for this fetch to work with the fake middleware service,
run: `go run cmd/fake-service/*.go` first from the root folder.
It might initially help to run it repeatedly until the installing 
errors disappear, make sure that Go is in the latest version too. */

store.dispatch("projectsStore/fetchProjects");
// Asynchronously request and receive recommendations from the middleware
// store.dispatch("recommendationsStore/fetchRecommendations");
// Start status watchers
// store.dispatch("recommendationsStore/startCentralStatusWatcher");

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
