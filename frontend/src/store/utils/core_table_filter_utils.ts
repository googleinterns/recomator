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

// The primary motivation for having these in a separate file is so that
//  they can be mocked by Jest:
//  https://stackoverflow.com/questions/51269431/jest-mock-inner-function

import { RecommendationExtra } from "../recommendation_extra";
import { ICoreTableStoreState, costCategoriesNames } from "../core_table";

// case-insensitive search (searching for 'myvm' should match 'MyVm')
export function isSearchTextInCell(
  searchText: string,
  cellText: string
): boolean {
  return cellText.toLowerCase().indexOf(searchText.toLowerCase()) !== -1;
}

// Individual filters: (each one returns true for empty input)

export function projectFilterAccepted(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
) {
  return (
    tableState.projectsSelected.length === 0 ||
    tableState.projectsSelected.includes(recExtra.projectCol)
  );
}

export function typeFilterAccepted(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
) {
  return (
    tableState.typesSelected.length === 0 ||
    tableState.typesSelected.includes(recExtra.typeCol)
  );
}

export function statusFilterAccepted(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
) {
  return (
    tableState.statusesSelected.length === 0 ||
    tableState.statusesSelected.includes(recExtra.statusCol)
  );
}

export function costFilterAccepted(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
) {
  return (
    tableState.costCategoriesSelected.length == 0 ||
    (recExtra.costCol >= 0 &&
      tableState.costCategoriesSelected.includes(costCategoriesNames.costs)) ||
    (recExtra.costCol <= 0 &&
      tableState.costCategoriesSelected.includes(costCategoriesNames.gains))
  );
}

export function resourceFilterAccepted(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
) {
  return (
    tableState.resourceNameSearchText.length === 0 ||
    isSearchTextInCell(tableState.resourceNameSearchText, recExtra.resourceCol)
  );
}

export function descriptionFilterAccepted(
  tableState: ICoreTableStoreState,
  recExtra: RecommendationExtra
) {
  return (
    tableState.descriptionSearchText.length === 0 ||
    isSearchTextInCell(tableState.descriptionSearchText, recExtra.description)
  );
}
