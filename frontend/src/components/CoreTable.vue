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
  <v-data-table
    dense
    :headers="headers"
    :items="filteredRecommendations"
    show-group-by
    v-on:update:group-by="onGroupByUpdated"
    :items-per-page="itemsPerPage"
    :single-select="false"
    v-model="$store.state.recommendationsStore.selected"
    show-select
    item-key="name"
  >
    <!-- customFilter prop is not used, because its implementation executes it for each property -->
    <template v-slot:body.prepend>
      <FiltersRow />
    </template>

    <!-- TODO: re-enable this warnings once these are implemented -->

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:group.summary="props">
      <!-- TODO: Group summary -->
    </template>

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:item.resource="{ item }">
      <ResourceCell :rowRecommendation="item" />
    </template>

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:item.project="{ item }">
      <ProjectCell :rowRecommendation="item" />
    </template>

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:item.recommenderSubtype="{ item }">
      <TypeCell :rowRecommendation="item" />
    </template>

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:item.description="{ item }">
      <!-- TODO: Description column -->
    </template>

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:item.cost="{ item }">
      <!-- TODO: Savings/Cost column -->
    </template>

    <!-- eslint-disable-next-line vue/no-unused-vars -->
    <template v-slot:item.apply="{ item }">
      <!-- TODO: Apply/status column -->
    </template>
  </v-data-table>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import FiltersRow from "@/components/FiltersRow.vue";
import ResourceCell from "@/components/ResourceCell.vue";
import ProjectCell from "@/components/ProjectCell.vue";
import TypeCell from "@/components/TypeCell.vue";
import { IRootStoreState } from "../store/root";
import { ICoreTableStoreState } from "../store/core_table";
import {
  Recommendation,
  getRecommendationProject,
  getRecommendationType,
  getRecommendationResourceShortName,
  getRecomendationDescription
} from "../store/model";

@Component({
  components: {
    FiltersRow,
    ResourceCell,
    ProjectCell,
    TypeCell
  }
})
export default class CoreTable extends Vue {
  headers = [
    { text: "Resource", value: "resource", groupable: false, sortable: true },
    { text: "Project", value: "project", groupable: true, sortable: true },
    {
      text: "Type",
      value: "recommenderSubtype",
      groupable: true,
      sortable: true
    },
    {
      text: "Description",
      value: "description",
      sortable: false,
      groupable: false
    },
    {
      text: "Savings/cost per week",
      value: "savingsAndCost",
      groupable: false
    },
    { text: "", value: "applyAndStatus", groupable: false, sortable: false }
  ];

  get filteredRecommendations() {
    const rootStoreState = this.$store.state as IRootStoreState;
    return rootStoreState.recommendationsStore!.recommendations.filter(
      (recommendation: Recommendation) =>
        CoreTable.filterPredicate(
          rootStoreState.coreTableStore!,
          recommendation
        )
    );
  }

  itemsPerPage = 10;
  // TODO: grouping should temporarily increase items shown to all

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onGroupByUpdated(groupByCategories: string[]): void {
    // TODO: update itemsPerPage
  }

  // TODO: once there is a new non-empty groupBy, close (toggle) all opened projects/types

  // returns true if the recommendation should be included in the results
  static filterPredicate(
    coreTableStoreState: ICoreTableStoreState,
    rec: Recommendation
  ): boolean {
    return (
      // project filter
      (coreTableStoreState.projectsSelected.length === 0 ||
        coreTableStoreState.projectsSelected.includes(
          getRecommendationProject(rec)
        )) &&
      // types selected
      (coreTableStoreState.typesSelected.length === 0 ||
        coreTableStoreState.typesSelected.includes(
          getRecommendationType(rec)
        )) &&
      // resource name search
      (coreTableStoreState.resourceNameSearchText.length === 0 ||
        getRecommendationResourceShortName(rec).indexOf(
          coreTableStoreState.resourceNameSearchText
        ) !== -1) &&
      // description search
      (coreTableStoreState.descriptionSearchText.length === 0 ||
        getRecomendationDescription(rec).indexOf(
          coreTableStoreState.descriptionSearchText
        ) !== -1)
    );
  }
}
</script>
