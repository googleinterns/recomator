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
  <v-app-bar app color="primary" dark>
    <v-btn text color="white" class="text-capitalize" @click="getHomePage()">
      <h1>Recomator</h1>
    </v-btn>
    <v-spacer />
    <slot />
  </v-app-bar>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { betterPush } from "./../router/better_push";

@Component({})
export default class AppBar extends Vue {
  getHomePage() {
    this.$store.dispatch("projectsStore/saveSelectedProjects");
    betterPush(this.$router, "HomeWithInit");
  }
}
</script>
