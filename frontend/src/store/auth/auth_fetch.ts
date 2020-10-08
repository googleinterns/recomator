import { IRootStoreState } from "./../root_state";
import { showError } from "@/router/show_error";
import { delay, infiniteDurationMs } from "../utils/misc";

function addAuth(
  init: RequestInit | undefined,
  rootState: IRootStoreState
): RequestInit {
  // first, make sure that init, init.headers exist
  if (init == undefined) init = {};
  if (init.headers == undefined) init.headers = {};

  init.headers = Object.assign({}, init.headers as Record<string, string>, {
    Authorization: `Bearer ${rootState.authStore!.idToken}`
  });
  return init;
}

// captures the authStore and router, so that authFetch can be used like a normal window.fetch
export function getAuthFetch(rootState: IRootStoreState, maxRetries = 0) {
  // - mimics window.fetch, but also adds the idToken to the request
  // - redirects to google sign-in if auth failed
  // - assumes init, inits.headers are in Record<string, *> format
  // - if the request fails (but not because of auth) and no retries left,
  // it will kill the app and redirect to an error page
  return async function authFetch(
    input: string,
    init?: RequestInit
  ): Promise<Response> {
    const initWithAuth = addAuth(init, rootState);

    let response: Response;
    let requestsLeft = maxRetries + 1;
    let retryWaitMs = 2000; // 2s
    while (requestsLeft > 0) {
      requestsLeft--;
      // wait before the second or later attempt
      if (requestsLeft < maxRetries) {
        // typically gives us 2+4+8=14s for temporary internet problems
        await delay(retryWaitMs);
        retryWaitMs *= 2;
      }

      try {
        response = await fetch(input, initWithAuth);
      } catch (error) {
        // server unreachable or another serious problem with connection
        const isCriticalError = requestsLeft === 0;
        await showError(
          `Failed to reach the Recomator backend`,
          {
            Message: error.Message,
            URL: input,
            Init: JSON.stringify(init),
            "Stack trace": error.stack
          },
          isCriticalError
        );
        continue;
      }

      if (response.ok) break;

      // not successful (outside of 200-299) but server is responsive

      const responseCode = response!.status;
      // access denied
      if (responseCode === 401 || responseCode === 403) {
        // redirect to GoogleSignIn: this will shut down everyting,
        // including fetching recommendations and check status
        await rootState.router!.push({ name: "GoogleSignIn" });
        // make sure we don't ever return from here
        await delay(infiniteDurationMs);
      }

      // server responsive, failed not because of auth
      const responseJSON = await response.json();
      await showError(
        `Network request failed: ${response!.status}(${response!.statusText})`,
        { URL: input, Init: JSON.stringify(init),
           ErrorMessage: responseJSON == null ? response.body : responseJSON.errorMessage},
        true
      );
    }
    return response!;
  };
}
