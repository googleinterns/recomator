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
  <td>
    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
        <v-chip
          :color="costColour"
          dark
          v-bind="attrs"
          block="true"
          v-on="on"
          display="fill"
        >
          {{ (cost >= 0 ? "" : "+") + Math.abs(costRounded) }}$
        </v-chip>
      </template>
      {{
        cost >= 0
          ? `Save ${costRounded}$ per week by applying this recommendation`
          : `Applying this recommendation will cost an additional ${-costRounded}$ per week`
      }}
    </v-tooltip>
  </td>
</template>
<script lang="ts">
import Vue from "vue";
import { Component } from "vue-property-decorator";
import { getRecommendationCostPerWeek } from "../store/model";

const SavingsCostCellProps = Vue.extend({
  props: ["rowRecommendation"]
});

@Component
export default class SavingsCostCell extends SavingsCostCellProps {
  get cost(): number {
    return getRecommendationCostPerWeek(this.rowRecommendation);
  }
  get costRounded(): string {
    return this.cost.toFixed(2);
  }
  get costColour(): string {
    return this.cost > 0 ? "orange" : "green";
  }
}
</script>
