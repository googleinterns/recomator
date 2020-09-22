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

import { delay, infiniteDurationMs } from "@/store/utils/misc";

// shows the given error on an error page
// refreshes the page to stop all current code (unless criticalError==false)
export async function showError(
  header: string,
  body: Record<string, string>,
  criticalError = true
): Promise<void> {
  if (!criticalError) {
    console.log(["Warning:", header, body]);
    return;
  }

  const headerParam = encodeURIComponent(header ? header : "No name");
  const bodyParam = encodeURIComponent(
    body ? JSON.stringify(body) : JSON.stringify({})
  );

  window.location.assign(`/errorMsg?header=${headerParam}&body=${bodyParam}`);

  // It might take a while for the browser to actually register this and kill our JS,
  // so let's just wait to make sure that no further code is executed (except the message queue)
  await delay(infiniteDurationMs);
}
