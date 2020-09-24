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
  <v-text-field
    v-model="resourceNameSearchText"
    append-icon="mdi-magnify"
    label="Filter resources"
    single-line
    hide-details
  />
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { IRootStoreState } from "../../store/root_state";

@Component
export default class ResourceFilter extends Vue {
  get resourceNameSearchText(): string {
    return (this.$store.state as IRootStoreState).coreTableStore!
      .resourceNameSearchText;
  }

  set resourceNameSearchText(text: string) {
    this.$store.commit("coreTableStore/setResourceNameSearchText", text);
  }
}
</script>
