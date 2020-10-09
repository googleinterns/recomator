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

export class ProjectConfig {
  public static DEVELOPMENT_BACKEND_ADDRESS = "http://localhost:8000/api";
  public static PRODUCTION_BACKEND_ADDRESS =
    process.env.VUE_APP_BACKEND_ADDRESS;
}

export function getBackendAddress(): string {
  return process.env.NODE_ENV === "development"
    ? ProjectConfig.DEVELOPMENT_BACKEND_ADDRESS
    : ProjectConfig.PRODUCTION_BACKEND_ADDRESS === undefined
    ? "https://recomator-282910.ey.r.appspot.com"
    : ProjectConfig.PRODUCTION_BACKEND_ADDRESS;
}
