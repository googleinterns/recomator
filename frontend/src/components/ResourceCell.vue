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
    <!-- Stackoverflow points to a potential security problem with opening the
         link in a new tab; the newly opened tab will have access to our data. 
         As long as the link is to Google Console, this shouldn't be a problem. 
         Related Stackoverflow question:
         https://stackoverflow.com/questions/17711146/how-to-open-link-in-new-tab-on-html 
    -->
    <a :href="consoleLink" target="_blank" rel="noopener noreferrer">
      {{ shortName }}
    </a>
  </div>
</template>
<script lang="ts">
import Vue, { PropType } from "vue";
import { Component } from "vue-property-decorator";
import { RecommendationExtra } from "../store/data_model/recommendation_extra";
import { getResourceConsoleLink } from "../store/data_model/recommendation_raw";

const ResourceCellProps = Vue.extend({
  props: {
    rowRecommendation: {
      type: Object as PropType<RecommendationExtra>,
      required: true
    }
  }
});

@Component
export default class ResourceCell extends ResourceCellProps {
  get shortName() {
    return this.rowRecommendation.resourceCol;
  }
  get consoleLink() {
    return getResourceConsoleLink(this.rowRecommendation);
  }
}
</script>
