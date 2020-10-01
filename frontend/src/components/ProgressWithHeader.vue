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
  <v-card color="primary" dark id="progressBar" persistent width="375">
    <v-card-title class="d-flex justify-center">
      <v-row align="center" justify="space-around">
        {{ header }} <slot></slot>
      </v-row>
    </v-card-title>
    <v-card-text>
      <v-progress-linear
        color="white"
        :indeterminate="indeterminate"
        :value="progress"
        class="mb-0"
      ></v-progress-linear>
    </v-card-text>
  </v-card>
</template>
<script lang="ts">
import Vue from "vue";
import { Component } from "vue-property-decorator";

const ProgressWithHeaderProps = Vue.extend({
  props: {
    progress: {
      type: Number,
      required: true
    },
    header: {
      type: String,
      required: true
    }
  }
});

@Component
export default class ProgressWithHeader extends ProgressWithHeaderProps {
  // progress tends to be 0 for the first few seconds, therefore we first show
  // the progress bar in the indeterminate state (progress animation independent of actual progress)
  get indeterminate(): boolean {
    return this.progress === 0;
  }
}
</script>

<style lang="scss">
#progressBar {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
