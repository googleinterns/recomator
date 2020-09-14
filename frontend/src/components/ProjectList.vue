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
  <v-card max-width="550" class="mx-auto">
    <v-toolbar color="primary" dark>
      <v-simple-checkbox
        class="mr-2"
        @click="selectAll"
        :indeterminate="!areAllSelected() && isSomeSelected()"
        :value="areAllSelected()"
      >
      </v-simple-checkbox>
      <v-toolbar-title> Select projects </v-toolbar-title>
      <v-spacer />
      <v-text-field
        v-if="searchEnabled"
        v-model="search"
        single-line
        hide-details
      >
        <template v-slot:append>
          <v-btn
            icon
            @click="
              searchEnabled = false;
              search = true;
            "
          >
            <v-icon> mdi-magnify </v-icon>
          </v-btn>
        </template>
      </v-text-field>

      <v-btn v-else icon @click="searchEnabled = true">
        <v-tooltip top transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs">mdi-magnify</v-icon>
          </template>
          Search projects
        </v-tooltip>
      </v-btn>
    </v-toolbar>
    <v-data-table
      v-model="selectedRows"
      :items="this.allRows"
      :hide-default-header="true"
      :headers="headers"
      :search="search"
      item-key="name"
      show-select
      class="elevation-1"
    >
    </v-data-table>

    <v-toolbar color="primary" dark>
      <v-spacer />
      <v-btn icon @click="getRequirements">
        <v-tooltip top transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs">mdi-equal-box</v-icon>
          </template>
          Proceed to testing requirements for the selected projects.
        </v-tooltip>
      </v-btn>

      <v-btn icon @click="getRecommendations">
        <v-tooltip top transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs">mdi-checkbox-marked-circle</v-icon>
          </template>
          Proceed to fetching recommendations from the selected projects.
        </v-tooltip>
      </v-btn>
    </v-toolbar>
  </v-card>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { IRootStoreState } from "../store/root";
import { Project } from "../store/data_model/project";

@Component({})
export default class ProjectList extends Vue {
  headers = [
    {
      value: "name"
    }
  ];

  searchEnabled = false;
  search = "";

  // Sync selected with the store
  get allRows(): Project[] {
    return (this.$store.state as IRootStoreState).projectsStore!.projects;
  }

  get selectedRows(): Project[] {
    return (this.$store.state as IRootStoreState).projectsStore!
      .projectsSelected;
  }

  set selectedRows(rows: Project[]) {
    this.$store.commit("projectsStore/setSelected", rows);
  }

  areAllSelected(): boolean {
    return this.allRows.length === this.selectedRows.length;
  }

  isSomeSelected(): boolean {
    return this.selectedRows.length > 0;
  }

  selectAll() {
    if (this.areAllSelected()) {
      this.selectedRows = [];
    } else {
      this.selectedRows = this.allRows;
    }
  }

  getRequirements() {
    this.$store.dispatch("projectsStore/proceedToRequirements")
  }

  getRecommendations() {
    this.$store.dispatch("projectsStore/proceedToRecommendations");
  }
}
</script>
