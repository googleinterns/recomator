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
    <v-toolbar color="primary" dark>
      <v-toolbar-title> Project requirements </v-toolbar-title>
      <v-spacer />
        <v-tooltip top transition="none">
          <template v-slot:activator="{ on, attrs }">
            <v-btn @click="getRecommendations" style="font-weight: bold" rounded depressed small v-on="on" v-bind="attrs" color="secondary" class="white--text">
              Recommendations
              <v-icon color="white">mdi-checkbox-marked-circle</v-icon>
            </v-btn>
          </template>
          Proceed to fetching recommendations from the selected projects.
        </v-tooltip>
    </v-toolbar>
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
            <th colspan="2" class="text-center">Other permissions</th>
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
          {{ getUnknownTooltip(item.item, requirement) }}
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
import { ProjectRequirement } from "../store/data_model/project_with_requirement";

@Component({})
export default class ProjectList extends Vue {
  requirementList = [
    "Service Usage API and services.get permission",
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

    "compute.disks.createSnapshot",
    "compute.snapshots.create",
    "compute.disks.delete",
    "recommender.computeDiskIdleResourceRecommendations.list",
    "recommender.computeDiskIdleResourceRecommendations.update",

    "compute.regions.list",
    "compute.zones.list",
  ];

  headers = [
    { value: "name" },

    { value: "Service Usage API and services.get permission" },
    { value: "Compute Engine API" },
    { value: "Cloud Resource Manager API" },
    { value: "Recommender API" },

    { value: "compute.instances.setMachineType" },
    { value: "compute.instances.start" },
    { value: "compute.instances.stop" },
    { value: "compute.instances.get" },
    { value: "recommender.computeInstanceIdleResourceRecommendations.update" },
    { value: "recommender.computeInstanceMachineTypeRecommendations.update" },
    { value: "recommender.computeInstanceIdleResourceRecommendations.list" },
    { value: "recommender.computeInstanceMachineTypeRecommendations.list" },

    { value: "compute.disks.createSnapshot" },
    { value: "compute.snapshots.create" },
    { value: "compute.disks.delete" },
    { value: "recommender.computeDiskIdleResourceRecommendations.list" },
    { value: "recommender.computeDiskIdleResourceRecommendations.update" },

    { value: "compute.regions.list" },
    { value: "compute.zones.list" },
  ].map((elt) => {
    if (elt.value === "name") {
      return elt;
    } else {
      return Object.assign(elt, { align: "center" });
    }
  });

  // Sync selected with the store
  get allRows(): ProjectRequirement[] {
    return (this.$store.state as IRootStoreState).requirementsStore!.projects;
  }

  getUknownText(
    projectRequirement: ProjectRequirement,
    requirement: string
  ): string {
    if (requirement === "compute.disks.createSnapshot") {
      if (projectRequirement.satisfiesRequirement("compute.snapshots.create")) {
        return `Requirement ${requirement} is not needed, as requirement compute.snapshots.create is satisfied`;
      }
    } else {
      if (requirement === "compute.snapshots.create") {
        if (
          projectRequirement.satisfiesRequirement(
            "compute.disks.createSnapshot"
          )
        ) {
          return `Requirement ${requirement} is not needed, as requirement compute.disks.createSnapshot is satisfied`;
        }
      } else {
        return `One of the requirements needed for checking the
          ${requirement} requirement is not satisfied.`;
      }
    }
  }

  getRecommendations() {
    this.$router.push("recommendations");
  }
}
</script>
