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

import { RecommendationExtra } from "./data_model/recommendation_extra";

export interface ICoreTableStoreState {
  resourceNameSearchText: string;
  descriptionSearchText: string;
  projectsSelected: string[];
  typesSelected: string[];
  statusesSelected: string[];
  costCategoriesSelected: string[];
  selected: RecommendationExtra[];
  currentlySelectable: RecommendationExtra[];
}

export function coreTableStoreStateFactory(): ICoreTableStoreState {
  return {
    resourceNameSearchText: "",
    descriptionSearchText: "",
    projectsSelected: [],
    typesSelected: [],
    statusesSelected: [],
    costCategoriesSelected: [],
    selected: [],
    currentlySelectable: []
  };
}
