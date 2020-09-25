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
    ref="check"
    dense
    :headers="headers"
    :items="this.$store.getters['filteredRecommendationsWithExtras']"
    :single-select="false"
    show-select
    v-model="selectedRows"
    v-on:current-items="setSelectableRows"
    item-key="name"
    :footer-props="{ itemsPerPageOptions: [10, 100, -1] }"
  >
    <!-- ^customFilter prop is not used, because its implementation executes it for each property -->
    <template v-slot:header.data-table-select="{}">
      <v-simple-checkbox
        :value="areAllSelected()"
        :indeterminate="!areAllSelected() && isSomeSelected()"
        v-on:input="selectAll"
      />
    </template>

    <template v-slot:item.data-table-select="{ item }">
      <v-simple-checkbox
        v-on:input="select(item, isSelected(item))"
        :value="isSelected(item) && isActive(item)"
        :disabled="!isActive(item)"
      />
    </template>

    <!-- The row with filters just above the data -->
    <template v-slot:body.prepend="{ isMobile }">
      <FiltersRow :isMobile="isMobile" />
    </template>

    <!-- Mappings of cell implementations to column slots: -->
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
import { getInternalStatusMapping } from "../store/data_model/status_map";
import { IRootStoreState } from "../store/root_state";
import { RecommendationExtra } from "../store/data_model/recommendation_extra";

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
  // headers ending with "Col" have values that are bound to corresponding columns
  //  for example, Resource will take RecommendationExtra.resourceCol for sorting
  headers = [
    {
      text: "Resource",
      value: "resourceCol",
      sortable: true
    },
    { text: "Project", value: "projectCol", sortable: true },
    {
      text: "Type",
      value: "typeCol",
      sortable: true
    },
    {
      text: "Description",
      value: "description",
      sortable: false
    },
    {
      text: "Savings/cost per week",
      value: "costCol",
      sortable: true
    },
    { text: "", value: "statusCol", sortable: false }
  ];

  get allRows(): RecommendationExtra[] {
    return this.$store.getters["filteredRecommendationsWithExtras"];
  }

  // Sync selected with the store
  get selectedRows(): RecommendationExtra[] {
    return (this.$store.state as IRootStoreState).coreTableStore!.selected;
  }

  set selectedRows(rows: RecommendationExtra[]) {
    this.$store.commit("coreTableStore/setSelected", rows);
  }

  get selectableRows(): RecommendationExtra[] {
    return (this.$store.state as IRootStoreState).coreTableStore!
      .currentlySelectable;
  }

  setSelectableRows(rows: RecommendationExtra[]) {
    this.$store.commit(
      "coreTableStore/setCurrentlySelectable",
      rows.filter(item => this.isActive(item))
    );
  }

  // If select parameter is true, selects all the rows on the current page
  // Otherwise unselects all rows on the current page
  selectAll(select: boolean) {
    select
      ? this.$store.commit("coreTableStore/selectAllSelectable")
      : this.$store.commit("coreTableStore/unselectAllSelectable");
  }

  // Checks if everything is selected on the current page
  areAllSelected(): boolean {
    return (
      this.selectableRows.every(item => this.isSelected(item)) &&
      this.selectableRows.length !== 0
    );
  }

  // Checks if anything is selected on the current page
  isSomeSelected(): boolean {
    return this.selectableRows.some(item => this.isSelected(item));
  }

  isActive(recommendation: RecommendationExtra) {
    return recommendation.statusCol === getInternalStatusMapping("ACTIVE");
  }

  isSelected(row: RecommendationExtra): boolean {
    return this.selectedRows.includes(row);
  }

  // If unselect is set to true, results in row being not selected.
  // If it is set to false, then results in row being selected.
  select(row: RecommendationExtra, unselect: boolean) {
    unselect
      ? this.$store.commit("coreTableStore/unselect", row)
      : this.$store.commit("coreTableStore/select", row);
  }
}
</script>
