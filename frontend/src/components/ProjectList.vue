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
  <v-card max-width="475" class="mx-auto">
    <v-toolbar color="primary" dark>
      <v-simple-checkbox class="mr-2" @click="selectAll" :indeterminate="!areAllSelected() && isSomeSelected()" :value="areAllSelected()"> </v-simple-checkbox>
      <v-toolbar-title> Select projects </v-toolbar-title>
      <v-spacer />
      <v-btn icon @click="acceptSelection">
        <v-tooltip top transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs">mdi-checkbox-marked-circle</v-icon>
          </template>
          Proceed to fetching recommendations from the selected projects.
        </v-tooltip>
      </v-btn>
    </v-toolbar>
    <v-data-table
      v-model="selectedRows"
      :items="this.allRows"
      :hide-default-header="true"
      :headers="headers"
      item-key="name"
      show-select
      class="elevation-1"
    >
    </v-data-table>
  </v-card>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import FiltersRow from "@/components/FiltersRow.vue";
import ResourceCell from "@/components/ResourceCell.vue";
import ProjectCell from "@/components/ProjectCell.vue";
import DescriptionCell from "@/components/DescriptionCell.vue";
import TypeCell from "@/components/TypeCell.vue";
import SavingsCostCell from "@/components/SavingsCostCell.vue";
import ApplyAndStatusCell from "@/components/ApplyAndStatusCell.vue";
import { getInternalStatusMapping } from "../store/data_model/status_map";
import { IRootStoreState } from "../store/root";
import { RecommendationExtra } from "../store/data_model/recommendation_extra";
import { Project, Requirement } from "../store/data_model/project";

@Component({
  components: {
    FiltersRow,
    ResourceCell,
    ProjectCell,
    TypeCell,
    DescriptionCell,
    SavingsCostCell,
    ApplyAndStatusCell,
  },
})
export default class ProjectList extends Vue {
  headers = [
    {
      value: "name",
    },
  ];

  expanded = [];

  // Sync selected with the store
  get allRows(): string[] {
    return (this.$store.state as IRootStoreState).projectsStore!.projects;
  }

  get selectedRows(): string[] {
    return (this.$store.state as IRootStoreState).projectsStore!
      .projectsSelected;
  }

  set selectedRows(rows: string[]) {
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

  acceptSelection() {
    this.$store.commit("projectsStore/chooseProjects");
  }
}
</script>
