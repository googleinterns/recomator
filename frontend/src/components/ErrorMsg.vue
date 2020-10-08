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
  <v-app>
    <AppBar />
    <v-main>
      <v-card id="errorCard" persistent max-width="700">
        <v-card-title id="elem" class="justify-center">
          <v-icon large color="error">mdi-alert-circle</v-icon>
          Error: {{ header }}
        </v-card-title>
        <v-expansion-panels accordion multiple v-model="openedIndices">
          <v-expansion-panel v-for="(value, name) in body" :key="name">
            <v-expansion-panel-header>{{ name }}</v-expansion-panel-header>
            <v-expansion-panel-content>
              {{ value }}
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-card>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { Component } from "vue-property-decorator";
import AppBar from "./AppBar.vue";

const ErrorMsgProps = Vue.extend({
  props: {
    header: {
      type: String,
      required: true
    },
    body: {
      type: Object as PropType<Record<string, string>>,
      required: true
    }
  }
});

@Component({
  components: {
    AppBar
  }
})
export default class ErrorMsg extends ErrorMsgProps {
  openedIndices: number[] = [];
  mounted() {
    const numberOfSections = Object.keys(this.body).length;

    // [0,1,...,numberOfSections-1]
    this.openedIndices = Array.from(Array(numberOfSections).keys());
  }
}
</script>

<style lang="scss">
#errorCard {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
#elem {
  -webkit-hyphens: none;
  -moz-hyphens: none;
  -ms-hyphens: none;
  hyphens: none;
}
</style>
