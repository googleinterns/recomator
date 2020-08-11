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
  <tr>
    <td />
    <td />
    <td>
      <v-combobox
        v-model="projectsSelected"
        :items="allProjects"
        label="Select projects"
        multiple
      >
      </v-combobox>
    </td>
    <td>
      <v-combobox
        v-model="typesSelected"
        :items="allTypes"
        label="Select types"
        multiple
      >
      </v-combobox>
    </td>
    <td />
    <td></td>
    <td>
      <v-combobox
        v-model="statusesSelected"
        :items="allStatuses"
        label="Select status"
        multiple
      >
      </v-combobox>
    </td>
  </tr>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { IRootStoreState } from "../store/root";
@Component
export default class FiltersRow extends Vue {
  // projects

  get allProjects(): string[] {
    // We could cache these, but filtering is the bottleneck so there is no point to bother
    return this.$store.getters["recommendationsStore/allProjects"];
  }

  get projectsSelected(): string[] {
    return (this.$store.state as IRootStoreState).coreTableStore!
      .projectsSelected;
  }

  set projectsSelected(projects: string[]) {
    this.$store.commit("coreTableStore/setProjectsSelected", projects);
  }

  // types

  get allTypes(): string[] {
    return this.$store.getters["recommendationsStore/allTypes"];
  }

  get typesSelected(): string[] {
    return (this.$store.state as IRootStoreState).coreTableStore!.typesSelected;
  }

  set typesSelected(types: string[]) {
    this.$store.commit("coreTableStore/setTypesSelected", types);
  }

  // statuses

  get statusesSelected(): string[] {
    return (this.$store.state as IRootStoreState).coreTableStore!
      .statusesSelected;
  }

  set statusesSelected(statuses: string[]) {
    this.$store.commit("coreTableStore/setStatusesSelected", statuses);
  }

  get allStatuses(): string[] {
    return this.$store.getters["recommendationsStore/allStatuses"];
  }
}
</script>
