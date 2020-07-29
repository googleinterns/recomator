<!-- Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. -->

<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <h1>Recomator</h1>
    </v-app-bar>
    <v-main>
      <v-progress-linear
        v-if="$store.state.recommendationsStore.progress !== null"
        :value="$store.state.recommendationsStore.progress"
        data-name="main_progress_bar"
      />
      <v-container
        fluid
        v-if="$store.state.recommendationsStore.progress === null"
        data-name="main_container"
      >
        <h1>All the app content will go here</h1>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Vuetify from "vuetify";

Vue.use(Vuetify);

@Component
export default class Home extends Vue {}
</script>
