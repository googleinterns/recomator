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
        v-if="$store.state.recommendationsStore.progress === null"
        tile
        @click="getProjectSelection"
        color="primary"
      >
        <v-tooltip left transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs" left color="white">mdi-pencil</v-icon>
          </template>
          Change selected projects
        </v-tooltip>
        Edit projects
      </v-btn>
    </AppBar>
    <v-main>
      <ProgressWithHeader
        v-if="$store.state.recommendationsStore.progress !== null"
        :progress="$store.state.recommendationsStore.progress"
        header="Loading recommendations..."
        data-name="main_progress_bar"
      />

      <v-container
        fluid
        data-name="main_container"
        v-if="$store.state.recommendationsStore.progress === null"
      >
        <PermissionDialog />
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
import AppBar from "@/components/AppBar.vue";
import CoreTable from "@/components/CoreTable.vue";
import Footer from "@/components/Footer.vue";
import PermissionDialog from "@/components/PermissionDialog.vue";

import ProgressWithHeader from "@/components/ProgressWithHeader.vue";
import { betterPush } from "./../router/better_push";

@Component({
  components: {
    CoreTable,
    Footer,
    ProgressWithHeader,
    AppBar,
    PermissionDialog
  }
})
export default class Home extends Vue {
  getProjectSelection() {
    betterPush(this.$router, "ProjectsWithInit");
  }
}
</script>
