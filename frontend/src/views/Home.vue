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
      <div class="text-center">
        <v-card
          color="primary"
          dark
          style="position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%);"
          v-if="$store.state.recommendationsStore.progress !== null"
          hide-overlay
          persistent
          width="300"
        >
          <v-card-title class="justify-center">
            Loading recommendations...
          </v-card-title>
          <v-card-text>
            <v-progress-linear
              color="white"
              :indeterminate="$store.state.recommendationsStore.progress === 0"
              :value="$store.state.recommendationsStore.progress"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </div>
      <v-container
        fluid
        data-name="main_container"
        v-if="$store.state.recommendationsStore.progress === null"
      >
        <v-row>
          <v-col>
            <CoreTable />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
    <Footer data-name="main-footer" />
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import CoreTable from "@/components/CoreTable.vue";
import Footer from "@/components/Footer.vue";

@Component({
  components: { CoreTable, Footer }
})
export default class Home extends Vue {}
</script>
