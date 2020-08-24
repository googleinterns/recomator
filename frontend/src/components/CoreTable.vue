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
    :single-select="false"
    show-select
    v-model=selectedRows
    v-on:pagination="pagination"
    item-key="name"
    :footer-props="{ itemsPerPageOptions: [10, 100, -1] }"
  >
    <!-- ^customFilter prop is not used, because its implementation executes it for each property -->
    <template v-slot:header.data-table-select="{ }">
      <v-simple-checkbox
      :value="areAllSelected()"
        v-on:input="selectAll"
      />
    </template>


    <template v-slot:item.data-table-select="{ item }">
      <v-simple-checkbox
        v-on:input="select(item, isSelected(item))"
        :value="isSelected(item) && isActiveRecommendation(item)"
        :disabled="!isActiveRecommendation(item)"
      />
    </template>

     <template v-slot:body.prepend="{ isMobile }">
      <FiltersRow :isMobile="isMobile" />
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
import { RecommendationExtra, getInternalStatusMapping } from "../store/model";
import { IRootStoreState } from "../store/root";

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

  // Array that keeps all recommendations that are selected
  // The array in the store keeps only the recommendations that
  // are selected and are applicable
  firstRow = 0;
  lastRow = Infinity;
  selectableCount = 0;
  selectableRows = new Array<RecommendationExtra>();

  get allRows(): RecommendationExtra[] {
    return this.$store.getters['filteredRecommendationsWithExtras'];
  }

  get selectedRows(): RecommendationExtra[] {
    return (this.$store.state as IRootStoreState).coreTableStore!.selected;
  }

  set selectedRows(selected: RecommendationExtra[]) {
    this.$store.commit(
      "coreTableStore/setSelected",
      selected);
  }

  pagination(event: any) {
    this.firstRow = event.pageStart;
    this.lastRow = event.pageStop;
    this.selectableRows = this.allRows.slice(this.firstRow, this.lastRow).filter(item => this.isActiveRecommendation(item));
    this.selectableCount = this.selectableRows.reduce(acc => acc + 1, 0);
  }

  selectAll(value: boolean) {
    if (value) {
      this.selectedRows = this.selectableRows;
    } else {
      this.selectedRows = [];
    }
  }

  areAllSelected(): boolean {
    return this.selectedRows.length === this.selectableCount;
  }

  isActiveRecommendation(item: RecommendationExtra) {
    return item.statusCol === getInternalStatusMapping("ACTIVE");
  }

  isSelected(item: RecommendationExtra): boolean {
    return this.selectedRows.includes(item);
  }

  select(item: RecommendationExtra, value: boolean) {
    if (value) {
      this.selectedRows = this.selectedRows.filter(elt => elt != item);
    } else {
      const copySelected = this.selectedRows.slice();
      copySelected.push(item);
      this.selectedRows = copySelected;
    }
  }
}
</script>
