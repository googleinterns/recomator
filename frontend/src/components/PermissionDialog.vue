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
    <v-dialog v-model="permissionDialogOpened" max-width="750px">
      <v-card>
        <v-card-title class="headline">
          <v-spacer />

          <v-col>
            <v-row>
              Loading partially failed
            </v-row>
          </v-col>
          <v-spacer />
        </v-card-title>
        <v-card-text>
          We were unable to fetch recommendations for the projects listed below.
          It may be caused by some requirements not being satisfied. For more
          details check permissions.
          <v-card class="ma-5">
            <v-list dense>
              <v-list-item v-for="item in failedProjects" :key="item">
                {{ item }}
              </v-list-item>
            </v-list>
          </v-card>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color="primary"
            dark
            v-on:click="
              permissionDialogOpened = false;
              getPermissions();
            "
          >
            <v-icon>mdi-equal-box</v-icon>
            Permissions
          </v-btn>
          <v-btn
            color="primary"
            dark
            v-on:click="permissionDialogOpened = false"
          >
            <v-icon>mdi-window-close</v-icon>
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style>
.v-card__text,
.v-card__title {
  word-break: normal;
}
</style>

<script lang="ts">
import Vue from "vue";
import { Component } from "vue-property-decorator";

import { IRootStoreState } from "../store/root_state";

@Component
export default class PermissionsDialog extends Vue {
  getPermissions() {
    this.$router.push("/requirements");
  }
  get failedProjects() {
    return (this.$store.state as IRootStoreState).recommendationsStore!
      .failedProjects;
  }

  permissionDialogOpened = this.failedProjects.length !== 0;
}
</script>
