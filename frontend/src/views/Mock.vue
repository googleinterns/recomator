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
    <v-app-bar app color="primary" dark>
      <h1>Recomator</h1>
    </v-app-bar>
    <v-progress-linear v-if="!successfullyLoaded" :value="progressPercentage">
    </v-progress-linear>
    <v-main v-if="successfullyLoaded">
      <v-card class="pa-5">
        <h2>{{ summary.toString() }}</h2>
        <v-btn rounded color="primary" dark small
          >Apply All Recomendations</v-btn
        >
      </v-card>
      <v-simple-table>
        <thead>
          <tr>
            <th v-for="header in headers" v-bind:key="header.value">
              {{ header.text }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(recommendation, index) in recommendations"
            v-bind:key="index"
          >
            <td class="text-left">{{ recommendation.project }}</td>
            <td class="text-left">
              <a :href="recommendation.path"> {{ recommendation.name }} </a>
            </td>
            <td class="text-left">{{ recommendation.getDescription() }}</td>
            <td class="text-left">
              <v-btn
                rounded
                color="primary"
                v-if="recommendation.applicable()"
                dark
                x-small
                >Apply Recommendation</v-btn
              >
              <v-progress-circular
                color="primary"
                indeterminate="true"
                v-if="recommendation.inProgress()"
              ></v-progress-circular>
              <div
                v-if="
                  !recommendation.applicable() && !recommendation.inProgress()
                "
              >
                {{ recommendation.status }}
              </div>
            </td>
          </tr>
        </tbody>
      </v-simple-table>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

class Recommendation {
  type: string;
  cost: string;
  path: string;
  project: string;
  name: string;
  status: string;

  constructor(
    type: string,
    cost: string,
    path: string,
    project: string,
    name: string,
    status: string
  ) {
    this.type = type; // UPPERCASED, contained in the set {RESIZE, REMOVE, PERFORMANCE} (doing this should be possible)
    this.cost = cost; // monthly cost, parsed to be positive (I assume, that the number would always be negative for resize and delete, and positive for performance)
    this.path = path;
    this.project = project;
    this.name = name;
    this.status = status;
  }

  applicable(): boolean {
    return this.status == "ACTIVE";
  }

  inProgress(): boolean {
    return this.status == "__INPROGRESS";
  }

  getDescription(): string {
    switch (this.type) {
      case "RESIZE": {
        return `Resize the VM to save ${this.cost} a month.`;
      }
      case "REMOVE": {
        return `Delete the VM to save ${this.cost} a month.`; // Should something be added about saving the machine's state?
      }
      case "PERFORMANCE": {
        return `Increase performance of the VM by spending additional ${this.cost} a month.`; // Should we be putting the cost here?
      }
      default: {
        return `Unknown recommendation.`;
      }
    }
  }
}

class Summary {
  private recommendationCount: number;
  private moneySaved: string;

  constructor(recommendationList: Recommendation[]) {
    this.recommendationCount = recommendationList.length;
    this.moneySaved = Summary.calculateSavings(recommendationList);
  }

  private static getCurrency(cost: string): string {
    return cost.slice(-1); //TODO take all the not numeric characters from the end? It would be good to write some regex.
  }

  private static costToNumber(cost: string): number {
    return parseFloat(cost.slice(0, -1));
  }

  private static calculateSavings(
    recommendationList: Recommendation[]
  ): string {
    let result = 0;
    let currency = "";

    for (const recommendation of recommendationList) {
      if (["RESIZE", "REMOVE", "PERFORMANCE"].includes(recommendation.type)) {
        const cost: number = Summary.costToNumber(recommendation.cost);
        currency = Summary.getCurrency(recommendation.cost);

        if (["RESIZE", "REMOVE"].includes(recommendation.type)) {
          result += cost;
        } else {
          result -= cost;
        }
      }
    }

    return result.toFixed(2) + currency;
  }

  public toString(): string {
    const moneySavedCount = Summary.costToNumber(this.moneySaved);
    if (moneySavedCount > 0) {
      return `Apply ${this.recommendationCount} recommendations to save ${this.moneySaved} every month.`;
    }

    return `Spend ${this.moneySaved.slice(
      1
    )} more each month to increase the performance by applying ${
      this.recommendationCount
    } recommendations.`;
  }
}

@Component
export default class Mock extends Vue {
  private headers = [
    { text: "Project", value: "name" },
    { text: "Name", value: "name" },
    { text: "Description", align: "start", sortable: false, value: "type" },
    { text: "Apply", value: "apply" }
  ];

  private recommendations = [
    new Recommendation(
      "RESIZE",
      "6.50$",
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-test",
      "timus-test-for-probers-n2-std-4-bored",
      "ACTIVE"
    ),
    new Recommendation(
      "REMOVE",
      "10.00$",
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-prod",
      "timus-test-for-probers-n2-std-4-very-bored",
      "SUCCEEDED"
    ),
    new Recommendation(
      "SECURITY",
      "Security recommendations don't have a cost",
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-prod",
      "shcheshnyak-n2-std-7-unsecure",
      "__INPROGRESS"
    ),
    new Recommendation(
      "PERFORMANCE",
      "30.00$",
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "leftsizer-test",
      "shcheshnyak-test-for-probers-n2-std-4-toobusy",
      "FAILED"
    ),
    new Recommendation(
      "UNSPECIFIED",
      "123$",
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "leftsizer-test",
      "timus-test-for-probers-n2-std-4-unknown",
      "ACTIVE"
    ),
    new Recommendation(
      "SECURITY",
      "Security recommendations don't have a cost",
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-prod",
      "shcheshnyak-n2-std-7-unsecure-2",
      "__INPROGRESS"
    ),
  ];

  private summary = new Summary(this.recommendations);
  private successfullyLoaded = false;
  private progressPercentage = 0;

  private async mounted() {
    // Simulate fetching recommendations at the beginning
    let i = 0;
    const progressBarSteps = 25;
    const totalWaitTimeMs = 3000;
    for (; i < progressBarSteps; i++) {
      const delay = (ms: number) => new Promise(res => setTimeout(res, ms));
      await delay((totalWaitTimeMs / progressBarSteps) * 2 * Math.random());
      this.progressPercentage += 100 / progressBarSteps;
    }
    this.successfullyLoaded = true;
  }
}
</script>
