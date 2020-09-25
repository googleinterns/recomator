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
    <v-data-table
      :items="this.allRows"
      :hide-default-header="true"
      :headers="headers"
      item-key="name"
      class="text-center elevation-1"
    >
      <template v-slot:header="">
        <thead>
          <tr>
            <th colspan="1" class="text-center">Project</th>
            <th colspan="4" class="text-center">APIs</th>
            <th colspan="8" class="text-center">VM permissions</th>
            <th colspan="4" class="text-center">Disks permissions</th>
            <th colspan="3" class="text-center">Other permissions</th>
          </tr>
        </thead>
      </template>

      <template
        v-for="requirement in requirementList"
        v-slot:[`item.${requirement}`]="item"
      >
        <v-tooltip
          :key="requirement"
          v-if="item.item.satisfiesRequirement(requirement)"
          top
          transition="none"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs" color="green"
              >mdi-check-bold</v-icon
            >
          </template>
          Requirement for {{ requirement }} is satisfied.
        </v-tooltip>
        <v-tooltip
          :key="requirement"
          v-else-if="!item.item.hasRequirement(requirement)"
          top
          transition="none"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs" color="grey">mdi-help</v-icon>
          </template>
          One of the requirements needed for checking the
          {{ requirement }} requirement is not satisfied.
        </v-tooltip>
        <v-tooltip :key="requirement" v-else top transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-on="on" v-bind="attrs" color="red"
              >mdi-alert-circle</v-icon
            >
          </template>
          {{ item.item.getErrorMessage(requirement) }}
        </v-tooltip>
      </template>
    </v-data-table>
  </div>
</template>

<style scoped src="./ProjectsWithRequirements.vue">
table th + th {
  border-left: 1px solid #dddddd;
}
</style>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { IRootStoreState } from "../store/root_state";
import { ProjectRequirement } from "../store/data_model/project_with_requirements";

@Component({})
export default class ProjectList extends Vue {
  requirementList = [
    "Service Usage API",
    "Compute Engine API",
    "Cloud Resource Manager API",
    "Recommender API",

    "compute.instances.setMachineType",
    "compute.instances.start",
    "compute.instances.stop",
    "compute.instances.get",
    "recommender.computeInstanceIdleResourceRecommendations.list",
    "recommender.computeInstanceMachineTypeRecommendations.list",
    "recommender.computeInstanceIdleResourceRecommendations.update",
    "recommender.computeInstanceMachineTypeRecommendations.update",

    "compute.disks.createSnapshot or compute.snapshots.create",
    "compute.disks.delete",
    "recommender.computeDiskIdleResourceRecommendations.list",
    "recommender.computeDiskIdleResourceRecommendations.update",

    "compute.regions.list",
    "compute.zones.list",
    "serviceusage.services.get"
  ];

  headers = ([] as { value: string; align?: string }[]).concat(
    [{ value: "name" }],
    this.requirementList.map(reqName => {
      return { value: reqName, align: "center" };
    })
  );

  get allRows(): ProjectRequirement[] {
    return (this.$store.state as IRootStoreState).requirementsStore!.projects;
  }

  getRecommendations() {
    this.$store.dispatch("projectsStore/saveSelectedProjects");
    this.$router.push("homeWithInit");
  }

  getProjects() {
    this.$router.push("projects");
  }
}
</script>
