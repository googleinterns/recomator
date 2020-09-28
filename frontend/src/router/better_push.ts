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

import VueRouter from "vue-router";

// VueRouter returns errors during router.push(*) if we redirect (using, let's say next) during
// the navigation, as apparently (the Vue.js team member claims) it is a potentially buggy case.
// That is why router-link passes an empty function as onComplete to router.push, as this changes the
// behaviour of router not to return a Promise
export function betterPush(router: VueRouter, routeName: string) {
  router.push({ name: routeName }, undefined, err => {
    console.log("Navigation warning:", err);
  });
}
