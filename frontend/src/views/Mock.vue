<template>
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

  <v-app>
    <v-app-bar app color="primary" dark></v-app-bar>
    <v-main>
      <h1>Recomator</h1>
      <v-simple-table>
        <thead>
          <tr>
            <th v-for="header in headers1" v-bind:key="header.value">
              <span> {{ header.text }} </span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(recommendation, index) in recommendations1"
            v-bind:key="index"
          >
            <td>{{ recommendation.type }}</td>
            <td>{{ recommendation.cost }}</td>
            <td>{{ recommendation.path }}</td>
            <td>
              <v-btn rounded color="primary" :disabled="false" dark x-small
                >Apply Recommendation</v-btn
              >
            </td>
            <td>{{ recommendation.status }}</td>
          </tr>
        </tbody>
      </v-simple-table>
    </v-main>
  </v-app>
</template>

<script lang="ts">
class Recommendation {
  type: string;
  cost: string;
  path: string;
  status: string;

  constructor(type: string, cost: string, path: string, status: string) {
    this.type = type;
    this.cost = cost;
    this.path = path;
    this.status = status;
  }
}

export default {
  name: "Mock",
  data: () => {
    return {
      headers1: [
        { text: "Type", align: "start", sortable: false, value: "type" },
        { text: "Cost", value: "cost" },
        { text: "VM path", value: "path" },
        { text: "Apply", value: "apply" },
        { text: "Status", value: "status" }
      ],
      recommendations1: [new Recommendation("xxx", "xxx", "xxx", "applicable")],
      headers2: [],
      recommendations2: []
    };
  }
};
</script>
