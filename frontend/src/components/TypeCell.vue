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
    <v-chip class="ma-2" color="white" label>
      <v-icon left>{{ iconName() }}</v-icon>
      {{ recommenderSubtype() }}
    </v-chip>
  </td>
</template>
<script lang="ts">
import Vue from "vue";
import { Component } from "vue-property-decorator";
import { getRecommendationType } from "../store/model";

const TypeProps = Vue.extend({
  props: ["rowRecommendation"]
});

@Component
export default class TypeCell extends TypeProps {
  recommenderSubtype(): string {
    return getRecommendationType(this.rowRecommendation);
  }
  iconName(): string {
    switch (this.recommenderSubtype()) {
      case "CHANGE_MACHINE_TYPE":
        return "mdi-move-resize-variant";
      case "STOP_VM":
        return "mdi-monitor-off";
      case "INCREASE_PERFORMANCE":
        return "mdi-monitor-screenshot";
      case "SNAPSHOT_AND_DELETE_DISK":
        return "mdi-harddisk-remove";
      default:
        return "mdi-cloud-question";
    }
  }
}
</script>
