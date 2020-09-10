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
      <v-toolbar-title> Select projects </v-toolbar-title>
      <v-spacer />
      <v-btn icon @click="acceptSelection">
        <v-icon>mdi-checkbox-marked-circle</v-icon>
      </v-btn>
    </v-toolbar>
    <v-data-table
      v-model="selectedRows"
      :items="this.allRows"
      :hide-default-header="true"
      :headers="headers"
      :expanded.sync="expanded"
      item-key="name"
      show-select
      show-expand
      class="elevation-1"
    >
      <template v-slot:expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <v-list class="pa-0 ma-1" dense>
            <v-list-item class="font-weight-bold pa-0 ma-1">
            <v-list-item-content  dense>Requirements:</v-list-item-content>
          </v-list-item>

          <v-list-item class="text-caption pa-0 ma-1" dense v-for="requirement in item.requirements" :key="requirement">
            <v-list-item-content  dense>{{ requirement }}</v-list-item-content>
          </v-list-item>
          </v-list>
        </td></template
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
import { Project } from "../store/data_model/project";

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
    {
      value: "data-table-expand",
    },
  ];

  expanded = [];

  // Sync selected with the store
  get allRows(): Project[] {
    return (this.$store.state as IRootStoreState).projectsStore!.projects;
  }

  get selectedRows(): string[] {
    return (this.$store.state as IRootStoreState).projectsStore!
      .projectsSelected;
  }

  set selectedRows(rows: string[]) {
    this.$store.commit("projectsStore/setSelected", rows);
  }

  acceptSelection() {
    this.$store.commit("projectsStore/chooseProjects");
  }
}
</script>
