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
    :items="$store.state.recommendationsStore.recommendations"
    show-group-by
    v-on:update:group-by="onGroupByUpdated"
    :items-per-page="itemsPerPage"
    :single-select="false"
    v-model="$store.state.recommendationsStore.selected"
    show-select
    item-key="instance"
  />
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class CoreTable extends Vue {
  private headers = [
    { text: "Name", value: "instance", groupable: false, sortable: true },
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

  private itemsPerPage = 10;
  // TODO: grouping should temporarily increase items shown to all

  private onGroupByUpdated(groupByCategories: string[]): void {
    // TODO: update itemsPerPage
  }

  // TODO: once there is a new non-empty groupBy, close (toggle) all opened projects/types
}
</script>
