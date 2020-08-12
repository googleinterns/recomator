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
  <div>
    <v-footer
      v-if="selectedRowsCount > 0"
      v-bind:fixed="true"
      color="rgba(255, 200, 20, 100)"
    >
    <v-row class="px-5">
      <v-spacer></v-spacer>

      <div class="font-weight-black">
        {{ footerMessage }}
      </div>
      <v-spacer></v-spacer>
      </v-row>
      <v-row class="px-5">
        <v-spacer></v-spacer>
        <v-btn rounded color="primary" dark v-on:click="applySelectedRecommendations()">Apply Selected Recommendations</v-btn>
        <v-spacer></v-spacer>
      </v-row>
    </v-footer>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import { Component } from "vue-property-decorator";
import { IRootStoreState } from "../store/root";
import { RecommendationExtra } from "../store/model";

@Component
export default class FooterSummary extends Vue {
get selectedRows(): RecommendationExtra[] {
    return (this.$store.state as IRootStoreState).coreTableStore!.selected;
  }

  get selectedRowsCount(): number {
    return this.selectedRows.length;
  }

  get savingsFromSelected(): number {
    const selectedRecommendations = this.selectedRows;

    let result = 0;
    for (const recommendation of selectedRecommendations) {
      if (recommendation.costCol < 0) {
        result -= recommendation.costCol; // Minus to have absolute value
      }
    }

    return result;
  }

  get spendingsFromSelected(): number {
    const selectedRecommendations = this.selectedRows;

    let result = 0;
    for (const recommendation of selectedRecommendations) {
      if (recommendation.costCol > 0) {
        result += recommendation.costCol;
      }
    }

    return result;
  }

  get performanceSelectedCount(): number {
    const selectedRecommendations = this.selectedRows;

    let result = 0;
    for (const recommendation of selectedRecommendations) {
      if (recommendation.costCol > 0) {
        result++;
      }
    }

    return result;
  }

  get applyPart(): string {
    const count = this.selectedRowsCount;
    if (count === 1) {
      return `Apply ${count} recommendation.`;
    }

    return `Apply ${count} recommendations.`;
  }

  get savingsPart(): string {
    const savings = this.savingsFromSelected;
    if (savings === 0) {
      return "";
    }

    return ` Save ${savings.toFixed(
      2
    )}$ each week by not using unnecessary resources.`;
  }

  get spendingsPart(): string {
    const count = this.performanceSelectedCount;
    const spendings = this.spendingsFromSelected;

    if (count === 0) {
      return "";
    }

    if (count === 1) {
      return ` Increase performance of ${count} machine, by spending ${spendings.toFixed(
        2
      )}$ more each week.`;
    }

    return ` Increase performance of ${count} machines, by spending ${spendings.toFixed(
      2
    )}$ more each week.`;
  }

  get footerMessage(): string {
    return this.applyPart + this.savingsPart + this.spendingsPart;
  }

  applySelectedRecommendations(): void {
    this.$store.dispatch("recommendationsStore/applyGivenRecommendations", this.selectedRows);
  }
}
</script>
