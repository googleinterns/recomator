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
    <AppBar>
      <v-btn
        v-if="$store.state.requirementsStore.progress === null"
        icon
        @click="getProjectSelection"
      >
        <v-tooltip left transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs" color="white">mdi-cog</v-icon>
          </template>
          Change selected projects
        </v-tooltip>
      </v-btn>
    </AppBar>
    <v-main>
      <ProgressWithHeader
        v-if="$store.state.requirementsStore.progress !== null"
        :progress="$store.state.requirementsStore.progress"
        header="Loading requirements..."
        data-name="requirement_progress_bar"
      />

      <v-container
        fluid
        v-if="$store.state.requirementsStore.progress === null"
      >
        <v-row>
          <v-col>
            <ProjectsWithRequirements />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import ProjectsWithRequirements from "@/components/ProjectsWithRequirements.vue";
import AppBar from "@/components/AppBar.vue";
import ProgressWithHeader from "@/components/ProgressWithHeader.vue";
import { betterPush } from "./../router/better_push";

@Component({
  components: { ProjectsWithRequirements, ProgressWithHeader, AppBar }
})
export default class Requirements extends Vue {
  getProjectSelection() {
    betterPush(this.$router, "ProjectsWithInit");
  }
}
</script>
