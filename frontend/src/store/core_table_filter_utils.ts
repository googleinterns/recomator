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

import { RecommendationExtra } from "./model";
import { ICoreTableStoreState } from "./core_table";

// The primary motivation for having these in a separate file is so that
//  they can be mocked by Jest:
//  https://stackoverflow.com/questions/51269431/jest-mock-inner-function

// We want 'save' in a search field to match 'Save' in a cell
export function isSearchTextInCell(
  searchText: string,
  cellText: string
): boolean {
  return cellText.toLowerCase().indexOf(searchText.toLowerCase()) !== -1;
}

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
