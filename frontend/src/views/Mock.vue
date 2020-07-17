<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png" />
    <HelloWorld msg="Welcome to Your Vue.js App" />
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from "@/components/HelloWorld.vue";

export default {
  name: "Home",
  components: {
    HelloWorld
  }
};
</script>
<template>
  <div class="mock">
    <h1>Recomator</h1>

    <v-simple-table>
      <thead>
        <tr>
          <th v-for="header in headers1" v-bind:key="header.value"><span> {{header.text}} </span></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(recommendation, index) in recommendations1" v-bind:key="index">
            <td> {{recommendation.type}} </td>
            <td> {{recommendation.cost}} </td>
            <td> {{recommendation.path}} </td>
            <td> <v-btn rounded color="primary" :disabled="recommendation.status == applicable" dark x-small>Apply Recommendation</v-btn> </td>
            <td> {{recommendation.status}} </td> 
        </tr>
      </tbody>
    </v-simple-table>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";

@Component
export default class Mock extends Vue {
  private headers1 = [
    { text: "Type", align: "start", sortable: false, value: "type" },
    { text: "Cost", value: "cost" },
    { text: "VM path", value: "path" },
    { text: "Apply", value: "apply" },
    { text: "Status", value: "status" },
  ];
  private recommendations1 = [
    {
      type: "xxx",
      cost: "xxx",
      path: "xxx",
      status: "applicable",
    },
  ];
  private headers2 = [];
  private recommendations2 = [];
  private name = "Mock";
}
</script>
