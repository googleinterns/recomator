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
        :value="$store.state.recommendationsStore.progress"
        data-name="main_progress_bar"
        v-if="
          $store.state.projectsStore.chosen &&
            $store.state.recommendationsStore.progress !== null
        "
      />
      <v-progress-linear
        v-if="!$store.state.projectsStore.loaded"
        indeterminate
      />
      <v-row>
        <v-col>
          <ProjectList
            v-if="
              $store.state.projectsStore.loaded &&
                !$store.state.projectsStore.chosen
            "
          />
        </v-col>
      </v-row>

      <v-container
        fluid
        data-name="main_container"
        v-if="
          $store.state.projectsStore.chosen &&
            $store.state.recommendationsStore.progress === null
        "
      >
        <v-row>
          <v-col>
            <CoreTable v-if="$store.state.projectsStore.chosen" />
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
import ProjectList from "@/components/ProjectList.vue";
import Footer from "@/components/Footer.vue";

@Component({
  components: { CoreTable, Footer, ProjectList }
})
export default class Home extends Vue {}
</script>
