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
    :items="this.$store.getters['filteredRecommendationsWithExtras']"
    show-group-by
    v-on:update:group-by="onGroupByUpdated"
    :items-per-page="itemsPerPage"
    :single-select="false"
    v-model="selectedRows"
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

    <template v-slot:item.resourceCol="{ item }">
      <ResourceCell :rowRecommendation="item" />
    </template>

    <template v-slot:item.projectCol="{ item }">
      <ProjectCell :rowRecommendation="item" />
    </template>

    <template v-slot:item.typeCol="{ item }">
      <TypeCell :rowRecommendation="item" />
    </template>

    <template v-slot:item.description="{ item }">
      <DescriptionCell :rowRecommendation="item" />
    </template>

    <template v-slot:item.costCol="{ item }">
      <SavingsCostCell :rowRecommendation="item" />
    </template>

    <template v-slot:item.statusCol="{ item }">
      <ApplyAndStatusCell :rowRecommendation="item" />
    </template>
  </v-data-table>
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
import { IRootStoreState } from "../store/root";
import { RecommendationExtra } from "../store/model";

@Component({
  components: {
    FiltersRow,
    ResourceCell,
    ProjectCell,
    TypeCell,
    DescriptionCell,
    SavingsCostCell,
    ApplyAndStatusCell
  }
})
export default class CoreTable extends Vue {
  // headers ending with "Col" have values that are bound to corresponding properties
  //  for example, Resource will take RecommendationExtra.resourceCol for sorting
  headers = [
    {
      text: "Resource",
      value: "resourceCol",
      groupable: false,
      sortable: true
    },
    { text: "Project", value: "projectCol", groupable: true, sortable: true },
    {
      text: "Type",
      value: "typeCol",
      groupable: true,
      sortable: true
    },
    {
      text: "Description",
      value: "description",
      groupable: false,
      sortable: false
    },
    {
      text: "Savings/cost per week",
      value: "costCol",
      groupable: false,
      sortable: true
    },
    { text: "", value: "statusCol", groupable: false, sortable: false }
  ];

  itemsPerPage = 10;
  // TODO: grouping should temporarily increase items shown to all

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onGroupByUpdated(groupByCategories: string[]): void {
    // TODO: update itemsPerPage
  }

  // TODO: once there is a new non-empty groupBy, close (toggle) all opened projects/types

  get selectedRows(): RecommendationExtra[] {
    return (this.$store.state as IRootStoreState).recommendationsStore!
      .selected;
  }

  set selectedRows(selected: RecommendationExtra[]) {
    this.$store.commit("recommendationsStore/setSelected", selected);
  }
}
</script>
