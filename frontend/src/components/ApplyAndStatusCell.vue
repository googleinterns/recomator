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
    <!-- Apply button/circular progress -->
    <v-btn
      rounded
      color="primary"
      v-if="checkStatus('ACTIVE') || checkStatus('CLAIMED')"
      v-on:click="applyRecommendation()"
      :loading="checkStatus('CLAIMED')"
      dark
      small
      block
      >Apply</v-btn
    >

    <v-btn
      rounded
      small
      label
      v-if="checkStatus('SUCCEEDED')"
      color="green"
      block
    >
      <v-icon color="white" left dark>mdi-check-circle</v-icon>
      {{ rowRecommendation.statusCol }}
    </v-btn>

    <v-dialog
      v-if="checkStatus('FAILED')"
      v-model="errorDialogOpened"
      max-width="600px"
    >
      <template v-slot:activator="{ on }">
        <v-btn color="red darken-2" v-on="on" small rounded block>
          <v-icon left color="white">mdi-alert-box</v-icon>
          Show Error
        </v-btn>
      </template>
      <v-card>
        <v-card-title class="headline">
          {{ rowRecommendation.errorHeader }}
        </v-card-title>
        <v-card-text>
          {{ rowRecommendation.errorDescription }}
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" dark v-on:click="applyRecommendation()">
            Retry
          </v-btn>
          <v-btn color="primary" dark v-on:click="errorDialogOpened = false">
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
<script lang="ts">
import Vue, { PropType } from "vue";
import { Component } from "vue-property-decorator";
import {
  RecommendationExtra,
  throwIfInvalidStatus,
  getInternalStatusMapping
} from "../store/model";

const ApplyAndStatusCellProps = Vue.extend({
  props: {
    rowRecommendation: {
      type: Object as PropType<RecommendationExtra>,
      required: true
    }
  }
});

@Component
export default class ApplyAndStatusCell extends ApplyAndStatusCellProps {
  errorDialogOpened = false;

  // we don't want a typo in status name to go unnoticed
  checkStatus(statusName: string) {
    throwIfInvalidStatus(statusName);
    return (
      this.rowRecommendation.statusCol === getInternalStatusMapping(statusName)
    );
  }

  applyRecommendation(): void {
    this.$store.dispatch("recommendationsStore/applyGivenRecommendations", [
      this.rowRecommendation.name
    ]);
  }
}
</script>
