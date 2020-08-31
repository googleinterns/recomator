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
  <!-- Shown in the first column on mobiles, and spread out otherwise -->
  <tr v-if="isMobile">
    <td>
      <ResourceFilter />
      <ProjectFilter />
      <TypeFilter />
      <DescriptionFilter />
      <CostFilter />
      <StatusFilter />
    </td>
  </tr>

  <tr v-else>
    <td />
    <td class="pb-5">
      <ResourceFilter />
    </td>
    <td>
      <ProjectFilter />
    </td>
    <td>
      <TypeFilter />
    </td>
    <td class="pb-5">
      <DescriptionFilter />
    </td>
    <td>
      <CostFilter />
    </td>
    <td>
      <StatusFilter />
    </td>
  </tr>
</template>

<script lang="ts">
import CostFilter from "./filters/CostFilter.vue";
import DescriptionFilter from "./filters/DescriptionFilter.vue";
import ProjectFilter from "./filters/ProjectFilter.vue";
import ResourceFilter from "./filters/ResourceFilter.vue";
import StatusFilter from "./filters/StatusFilter.vue";
import TypeFilter from "./filters/TypeFilter.vue";

import Vue from "vue";
import { Component } from "vue-property-decorator";

const FiltersRowProps = Vue.extend({
  props: {
    isMobile: {
      type: Boolean,
      required: true
    }
  }
});

@Component({
  components: {
    CostFilter,
    DescriptionFilter,
    ProjectFilter,
    ResourceFilter,
    StatusFilter,
    TypeFilter
  }
})
export default class FiltersRow extends FiltersRowProps {}
</script>
