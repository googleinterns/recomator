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
  <v-app mt-10 mb-10 ml-10 mr-10>
    <v-app-bar app color="primary" dark>
      <h1>Recomator</h1>
    </v-app-bar>

    <v-main>
      <v-progress-linear v-if="!successfullyLoaded" :value="progressPercentage">
      </v-progress-linear>
      <v-container fluid v-if="successfullyLoaded">
        <v-row>
          <v-col :cols="3">
            <v-card>
              <h2>Filters</h2>
            </v-card>
            <v-card class="ma-2 pa-2">
              <h3>Cost</h3>
              <v-text-field
                type="number"
                label="Minimal cost"
                v-on:input="filterRecommendation.setMinimalPrice($event)"
              ></v-text-field>
              <v-text-field
                type="number"
                label="Maximal cost"
                v-on:input="filterRecommendation.setMaximalPrice($event)"
              ></v-text-field>
            </v-card>

            <v-card class="ma-2 pa-2">
              <h3>Project</h3>
              <v-select
                :items="summary.getProjectList()"
                label="Select project"
                v-on:input="filterRecommendation.setProject($event)"
              >
              </v-select>

              <v-btn
                rounded
                color="primary"
                dark
                small
                v-on:click="filterRecommendation.setProject(undefined)"
                >Clear</v-btn
              >
            </v-card>
          </v-col>
          <v-col>
            <v-card class="pa-5">
              <!-- <h2>{{ summary.toString() }}</h2> -->
              <h3>
                Spend 1234$ more and save 2314$ per week by applying all
                recommendations
              </h3>
              <v-btn rounded color="primary" dark small
                v-on:click="applyAllRecommendations"
                >Apply All Recomendations</v-btn
              >
              <v-text-field
                v-model="search"
                append-icon="mdi-magnify"
                label="Search"
                single-line
                hide-details
              ></v-text-field>
            </v-card>

            <v-data-table
              :headers="headers"
              :items="
                recommendations.filter(recommendation =>
                  filterRecommendation.predicate(recommendation)
                )
              "
              show-group-by
              :search="search"
            >
              <template v-slot:item="recommendation">
                <tr>
                  <td class="text-left">{{ recommendation.item.project }}</td>
                  <td class="text-left">
                    <a :href="recommendation.item.path">
                      {{ recommendation.item.name }}
                    </a>
                  </td>
                  <td class="text-left">
                    {{ recommendation.item.description }}
                  </td>
                  <td>
                    <v-chip :color="recommendation.item.getCostColour()" dark>
                      {{ Math.abs(recommendation.item.cost) }}$</v-chip>
                  </td>

                  <td class="text-left">
                    <v-btn
                      rounded
                      color="primary"
                      v-if="recommendation.item.applicable()"
                      v-on:click="runRecommendation(recommendation.item)"
                      dark
                      x-small
                      >Apply</v-btn
                    >
                    <v-progress-circular
                      color="primary"
                      :indeterminate="true"
                      v-if="recommendation.item.inProgress()"
                    ></v-progress-circular>

                    <v-icon
                      color="green darken-2"
                      v-if="recommendation.item.succeded()"
                      >mdi-check-circle</v-icon
                    >
                    <v-icon
                      color="red darken-2"
                      v-if="recommendation.item.failed()"
                      >mdi-close-circle</v-icon
                    >
                    <div
                      v-if="
                        !recommendation.item.applicable() &&
                          !recommendation.item.inProgress() &&
                          !recommendation.item.failed() &&
                          !recommendation.item.succeded()
                      "
                    >
                      {{ recommendation.item.status }}
                    </div>
                  </td>
                </tr>
              </template>
            </v-data-table>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

class Recommendation {
  type: string;
  cost: number;
  path: string;
  project: string;
  name: string;
  status: string;
  description: string;

  constructor(
    type: string,
    cost: number,
    path: string,
    project: string,
    name: string,
    status: string,
    description: string
  ) {
    this.type = type; // UPPERCASED, contained in the set {RESIZE, REMOVE, PERFORMANCE} (doing this should be possible)
    this.cost = cost; // monthly cost, parsed to be positive (I assume, that the number would always be negative for resize and delete, and positive for performance)
    this.path = path;
    this.project = project;
    this.name = name;
    this.status = status;
    this.description = description;
  }

  copy(): Recommendation {
    return new Recommendation(
      this.type,
      this.cost,
      this.path,
      this.project,
      this.name,
      this.status,
      this.description
    );
  }

  private costToNumber(cost: string): number {
    return parseFloat(cost.slice(0, -1));
  }

  smallerCost(cost: number) {
    return this.cost < cost;
  }

  biggerCost(cost: number) {
    return this.cost > cost;
  }

  applicable(): boolean {
    return this.status == "ACTIVE";
  }

  inProgress(): boolean {
    return this.status == "__INPROGRESS";
  }

  succeded(): boolean {
    return this.status == "SUCCEEDED";
  }

  failed(): boolean {
    return this.status == "FAILED";
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
  
  private getCostColour(): string {
    return this.cost > 0 ? "green" : "orange";
  }
}

class FilterRecommendation {
  private minimalPrice = -Infinity;
  private maximalPrice = Infinity;
  private project: string | undefined = undefined;

  public setMaximalPrice(cost: string) {
    if (cost === "") {
      this.maximalPrice = Infinity;
      return;
    }

    this.maximalPrice = parseFloat(cost);
  }

  public setMinimalPrice(cost: string) {
    if (cost === "") {
      this.minimalPrice = -Infinity;
      return;
    }

    this.minimalPrice = parseFloat(cost);
  }

  public setProject(project: string | undefined) {
    if (project === "") {
      this.project = undefined;
      return;
    }

    this.project = project;
  }

  private pricePredicate(recommendation: Recommendation): boolean {
    return (
      recommendation.smallerCost(this.maximalPrice) &&
      recommendation.biggerCost(this.minimalPrice)
    );
  }

  private projectPredicate(recommendation: Recommendation): boolean {
    return (
      this.project === undefined || recommendation.project === this.project
    );
  }

  public predicate(recommendation: Recommendation): boolean {
    return (
      this.pricePredicate(recommendation) &&
      this.projectPredicate(recommendation)
    );
  }
}

class Summary {
  private recommendationCount: number;
  private moneySaved: string;
  private projectList = new Array<string>();

  constructor(recommendationList: Recommendation[]) {
    this.recommendationCount = recommendationList.length;
    this.moneySaved = Summary.calculateSavings(recommendationList);

    for (const recommendation of recommendationList) {
      if (!this.projectList.includes(recommendation.project)) {
        this.projectList.push(recommendation.project);
      }
    }
  }

  private static getCurrency(cost: string): string {
    return cost.slice(-1); //TODO take all the not numeric characters from the end? It would be good to write some regex.
  }

  public static costToNumber(cost: string): number {
    return parseFloat(cost.slice(0, -1));
  }

  public getProjectList(): string[] {
    return this.projectList;
  }

  private static calculateSavings(
    recommendationList: Recommendation[]
  ): string {
    let result = 0;

    for (const recommendation of recommendationList) {
      if (["RESIZE", "REMOVE", "PERFORMANCE"].includes(recommendation.type)) {
        const cost: number = recommendation.cost;
        result += cost;
      }
    }

    return result.toFixed(2) + " $";
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
  private filterRecommendation = new FilterRecommendation();

  private search = "";
  private headers = [
    { text: "Project", value: "project", filterable: true },
    { text: "Name", value: "name", groupable: false },
    {
      text: "Description",
      align: "start",
      sortable: false,
      value: "type",
      groupable: false
    },
    { text: "Savings/cost per week", value: "cost", groupable: false },
    { text: "", value: "apply", groupable: false }
  ];

  private recommendations_core = [
    new Recommendation(
      "RESIZE",
      6.5,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-test",
      "timus-test-for-probers-n2-std-4-bored",
      "ACTIVE",
      "Save cost by changing machine type from n1-highcpu-8 to n1-highcpu-4"
    ),
    new Recommendation(
      "REMOVE",
      10.0,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-prod",
      "timus-test-for-probers-n2-std-4-very-bored",
      "ACTIVE",
      "Save cost by removing machine"
    ),
    new Recommendation(
      "PERFORMANCE",
      -30.0,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "leftsizer-test",
      "shcheshnyak-test-for-probers-n2-std-4-toobusy",
      "ACTIVE",
      "Increase performance by changing machine type from n1-highcpu-4 to n1-highcpu-8"
    ),
    new Recommendation(
      "UNSPECIFIED",
      123,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "leftsizer-test",
      "timus-test-for-probers-n2-std-4-unknown",
      "ACTIVE",
      "Save cost by changing machine type from n1-highcpu-8 to n1-highcpu-4"
    ),
    new Recommendation(
      "RESIZE",
      3.5,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "middlesizer-test",
      "timus-test-for-probers-n2-std-4-bored-2",
      "ACTIVE",
      "Save cost by changing machine type from n1-highcpu-16 to n1-highcpu-8"
    ),
    new Recommendation(
      "REMOVE",
      10.0,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "rightsizer-prod",
      "timus-test-for-probers-n2-std-4-very-bored",
      "SUCCEEDED",
      "Save cost by removing machine"
    ),
    new Recommendation(
      "RESIZE",
      11.75,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "middlesizer-test",
      "timus-test-for-probers-n2-std-4-bored-3",
      "ACTIVE",
      "Save cost by changing machine type from n1-highcpu-4 to n1-highcpu-2"
    ),
    new Recommendation(
      "PERFORMANCE",
      -70.0,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "leftsizer-test",
      "shcheshnyak-test-for-probers-n2-std-4-toobusy",
      "ACTIVE",
      "Improve performance by changing machine type from n1-highcpu-4 to n1-highcpu-8"
    ),
    new Recommendation(
      "RESIZE",
      -11.5,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "search",
      "shcheshnyak-test-for-probers-n2-std-4-toobusy",
      "ACTIVE",
      "Save cost by changing machine type from n1-highcpu-16 to n1-highcpu-2"
    ),
    new Recommendation(
      "RESIZE",
      -7.8,
      "https://pantheon.corp.google.com/compute/instancesDetail/zones/us-central1-c/\
      instances/timus-test-for-probers-n2-std-4-idling?project=rightsizer-test&supportedpurview=project",
      "search",
      "shcheshnyak-test-for-probers-n2-std-4-toobusy",
      "ACTIVE",
      "Save cost by changing machine type from n1-highcpu-6 to n1-highcpu-2"
    )
  ];

  private generateRecommendations(): void {
    this.recommendations = Array(10)
      .fill(this.recommendations_core)
      .reduce((agg, arr) => agg.concat(arr), [])
      .sort(function() {
        return 0.5 - Math.random();
      })
      .map((rec: Recommendation) => rec.copy());
    this.recommendations.forEach((rec: Recommendation, index) => {
      rec.name += index;
    });
  }

  private recommendations = [];

  private shuffleRecommendations(): void {
    this.recommendations.sort(function() {
      return 0.5 - Math.random();
    });
  }

  private successfullyLoaded = false;
  private progressPercentage = 0;
  private summary = new Summary(this.recommendations_core);

  private async runRecommendation(recommendation: Recommendation) {
    recommendation.status = "__INPROGRESS";
    await new Promise(res => setTimeout(res, 2500));
    if (Math.random() < 0.33) recommendation.status = "FAILED";
    else {
      await new Promise(res => setTimeout(res, 5000));
      recommendation.status = "SUCCEEDED";
    }
  }

  private applyAllRecommendations() {
    this.recommendations.forEach((rec: Recommendation) => {
      if(rec.applicable())
        this.runRecommendation(rec);
    });
  }

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
    this.generateRecommendations();
  }
}
</script>
