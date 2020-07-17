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

  <v-container>
    <h1>{{ msg }}</h1>
    <v-row class="text-center">
      <v-col class="mb-4">
        <h1 class="display-2 font-weight-bold mb-3">
          Welcome to the... Recomator! {{ msg }}
        </h1>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vuetify from "vuetify";

import { Component, Prop, Vue } from "vue-property-decorator";
Vue.use(Vuetify);

@Component
export default class HelloWorld extends Vue {
  @Prop() private msg!: string;
}
</script>
