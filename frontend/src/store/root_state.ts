/* Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

import { IRecommendationsStoreState } from "./recommendations_state";
import { ICoreTableStoreState } from "./core_table_state";
import { IAuthStoreState } from "./auth_state";

export interface IRootStoreState {
  // Static type checking needs to know of these properties (added dynamically)
  //  so they are added as optional as a workaround:
  //  related issue: https://forum.vuejs.org/t/vuex-submodules-with-typescript/40903
  // Therefore, the ! operator needs to be used whenever the state of any module
  //  is accessed from outside.
  recommendationsStore?: IRecommendationsStoreState;
  coreTableStore?: ICoreTableStoreState;
  authStore?: IAuthStoreState;
}
