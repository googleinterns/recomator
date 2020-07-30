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
        v-if="$store.state.recommendationsStore.progress !== null"
      />
      <v-container
        fluid
        data-name="main_container"
        v-if="$store.state.recommendationsStore.progress === null"
      >
        <v-row>
          <v-col>
            <GroupingHint />
          </v-col>
        </v-row>

        <v-row>
          <v-col>
            <CoreTable />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import GroupingHint from "@/components/grouping_hint.vue";
import CoreTable from "@/components/core_table.vue";

@Component({
  components: { GroupingHint, CoreTable }
})
export default class Home extends Vue {}
</script>
