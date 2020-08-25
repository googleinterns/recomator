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
  <v-combobox
    v-model="costCategoriesSelected"
    :items="allCostCategories"
    label="Savings/Costs"
    multiple
  >
  </v-combobox>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { IRootStoreState } from "../../store/root";
import { costCategoriesNames } from "../../store/core_table";

@Component
export default class CostFilter extends Vue {
  // ["Gains", "Costs"] or some synonyms of these
  get allCostCategories(): string[] {
    return Object.values(costCategoriesNames);
  }

  get costCategoriesSelected(): string[] {
    return (this.$store.state as IRootStoreState).coreTableStore!
      .costCategoriesSelected;
  }
  set costCategoriesSelected(costCategories: string[]) {
    this.$store.commit(
      "coreTableStore/setCostCategoriesSelected",
      costCategories
    );
  }
}
</script>
