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
    <AppBar />
    <v-main>
      <ProgressWithHeader
        v-if="!$store.state.projectsStore.loaded"
        :progress="0"
        header="Loading projects..."
        data-name="project_progress_bar"
      />

      <v-container fluid v-if="$store.state.projectsStore.loaded">
        <v-row>
          <v-col>
            <ProjectList />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import ProjectList from "@/components/ProjectList.vue";
import AppBar from "../components/AppBar.vue";
import ProgressWithHeader from "@/components/ProgressWithHeader.vue";

@Component({
  components: { ProjectList, AppBar, ProgressWithHeader }
})
export default class Project extends Vue {}
</script>
