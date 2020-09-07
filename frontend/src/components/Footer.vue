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
    <!-- Footer with 'Apply Selected' button -->
    <v-footer
      ref="footer"
      v-if="selectedRowsCount > 0"
      v-bind:fixed="true"
      color="primary"
      v-bind:app="true"
    >
      <v-container fluid>
        <v-row align="center" justify="space-around" justify-md="space-between">
          <v-btn fab small dark color="secondary" v-on:click="unselectAll">
            <v-icon dark>mdi-close</v-icon>
          </v-btn>
          <div>
            <div
              data-name="footer-summary"
              style="font-weight: bold; color: white"
            >
              {{ applyPart }}<br />
              <span v-if="spendingsPart.length > 0"
                >{{ spendingsPart }}<br
              /></span>
              <span v-if="savingsPart.length > 0">{{ savingsPart }}</span>
            </div>
          </div>
          <div>
            <v-btn
              data-name="footer-button"
              rounded
              color="secondary"
              style="font-weight: bold"
              dark
              v-on:click="dialog = true"
              >Apply Selected</v-btn
            >
          </div>
        </v-row>
      </v-container>
    </v-footer>

    <!-- Confirmation dialog -->
    <v-dialog ref="dialog" v-model="dialog" max-width="640px">
      <v-card data-name="dialog">
        <v-card-title class="headline">
          <v-row>
            <v-col
              >Are you sure you want to apply
              {{ selectedRowsCount }} recommendation{{
                selectedRowsCount == 1 ? "" : "s"
              }}?</v-col
            ></v-row
          >
        </v-card-title>

        <v-card-actions>
          <v-spacer />
          <v-btn
            data-name="yes-button"
            color="green white--text"
            v-on:click="
              dialog = false;
              applySelectedRecommendations();
            "
          >
            <v-icon>mdi-check</v-icon>
            Yes
          </v-btn>

          <v-btn
            data-name="cancel-button"
            color="primary white--text"
            v-on:click="dialog = false"
          >
            <v-icon>mdi-window-close</v-icon>
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import { Component } from "vue-property-decorator";
import { IRootStoreState } from "../store/root";
import { RecommendationExtra } from "../store/data_model/recommendation_extra";

@Component
export default class Footer extends Vue {
  dialog = false;

  // Helpers:

  get selectedRows(): RecommendationExtra[] {
    return (this.$store.state as IRootStoreState).coreTableStore!.selected;
  }

  get selectedRowsCount(): number {
    return this.selectedRows.length;
  }

  get savingsFromSelected(): number {
    return this.selectedRows
      .filter(recommendation => recommendation.costCol < 0)
      .reduce((acc, cur) => acc - cur.costCol, 0);
  }

  get spendingsFromSelected(): number {
    return this.selectedRows
      .filter(recommendation => recommendation.costCol > 0)
      .reduce((acc, cur) => acc + cur.costCol, 0);
  }

  get performanceSelectedCount(): number {
    return this.selectedRows
      .filter(recommendation => recommendation.costCol > 0)
      .reduce(acc => acc + 1, 0);
  }

  // Footer summary generators:

  get applyPart(): string {
    return `Apply ${this.selectedRowsCount} recommendation(s):`;
  }

  get savingsPart(): string {
    const savings = this.savingsFromSelected;
    if (savings === 0) {
      return "";
    }

    return `Use less resources and save ${savings.toFixed(2)}$ each week.`;
  }

  get spendingsPart(): string {
    if (this.performanceSelectedCount === 0) {
      return "";
    }

    return ` Increase performance of ${this.performanceSelectedCount} machine${
      this.performanceSelectedCount == 1 ? "" : "s"
    } by spending ${this.spendingsFromSelected.toFixed(2)}$ more each week.`;
  }

  // Handler of the 'Apply all selected' button
  applySelectedRecommendations(): void {
    this.$store.dispatch(
      "recommendationsStore/applyGivenRecommendations",
      this.selectedRows.map(row => row.name)
    );
    this.unselectAll();
  }

  unselectAll(): void {
    this.$store.commit("coreTableStore/setSelected", []);
  }
}
</script>
