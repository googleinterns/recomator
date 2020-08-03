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
import { mount } from "@vue/test-utils";
import { rootStoreFactory } from "@/store/root";
import Home from "@/views/Home.vue";
import vuetify from "@/plugins/vuetify";

// Helpful for debugs: console.log(wrapper.html())
//
// await Vue.nextTick() is used to wait for Vue to actually update the components
//  as recommended here: https://vue-test-utils.vuejs.org/guides/#testing-asynchronous-behavior
//
// We might consider switching to local Vue instances/vuetifys
//  (there is a function for that in test-utils), as it is recommended
//  but there are issues with that, which have fixes described here:
//   https://vuetifyjs.com/en/getting-started/unit-testing/

describe("Home", () => {
  test("Progress bar and main container visibility", async () => {
    const fakeStore = rootStoreFactory();
    const wrapper = mount(Home, { store: fakeStore, vuetify: vuetify });

    function progressBarExists(): boolean {
      return wrapper.find("[data-name=main_progress_bar]").exists();
    }

    function mainContainerExists(): boolean {
      return wrapper.find("[data-name=main_container]").exists();
    }

    // simulate typical progress updates

    for (let i = 0; i <= 100; i += 10) {
      fakeStore.commit("recommendationsStore/setProgress", i);
      await Vue.nextTick();
      expect(progressBarExists()).toBeTruthy();
      expect(mainContainerExists()).toBeFalsy();
    }

    fakeStore.commit("recommendationsStore/setProgress", null);
    await Vue.nextTick();
    expect(progressBarExists()).toBeFalsy();
    expect(mainContainerExists()).toBeTruthy();
  });
});
