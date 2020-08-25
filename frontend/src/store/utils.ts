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

import { ProjectConfig } from "@/config";
import { internalStatusMap } from "@/store/recommendation_extra";

export function extractFromResource(
  property: string,
  resource: string
): string {
  const sliceLen = `/${property}/`.length;
  const pattern = `/${property}/[^/]*`;
  const regex = new RegExp(pattern);

  const found = regex.exec(resource);
  if (found === null) {
    // TODO: this function doesn't support snapshots, so I have temporarily disabled errors
    return "NOT IMPLEMENTED";
    //throw `couldn't parse resource identifier: ${resource}`;
  }

  const result = found[0].slice(sliceLen);
  return result;
}

export function delay(miliseconds: number) {
  return new Promise(resolve => setTimeout(resolve, miliseconds));
}

export function getServerAddress(): string {
  if (process.env.NODE_ENV === "development") {
    return ProjectConfig.DEVELOPMENT_SERVER_ADDRESS;
  }

  return ProjectConfig.PRODUCTION_SERVER_ADDRESS;
}

export function throwIfInvalidStatus(statusName: string): void {
  if (!(statusName in internalStatusMap))
    throw `invalid status name passed: ${statusName}`;
}

// We don't want to display API status names directly
//  For details on how this is stored, see RecommendationExtra definition above
export function getInternalStatusMapping(statusName: string): string {
  throwIfInvalidStatus(statusName);
  return internalStatusMap[statusName];
}
