import { IRootStoreState } from "./../root_state";
// captures the authStore and router, so that authFetch can be used like a normal fetch
export function getAuthFetch(rootState: IRootStoreState) {
  // adds the idToken to the request, redirects to google sign-in if auth failed
  // assumes init, inits.headers are in Record<string, *> format
  return async function authFetch(
    input: string,
    init?: RequestInit
  ): Promise<Response> {
    // first, make sure that init, init.headers exist
    if (init == undefined) init = {};
    if (init.headers == undefined) init.headers = {};

    init.headers = Object.assign({}, init.headers as Record<string, string>, {
      Authorization: `Bearer ${rootState.authStore!.idToken}`
    });

    const response = await fetch(input, init);
    const responseCode = response.status;

    // access denied
    if (responseCode === 401) {
      // redirect to GoogleSignIn: this will shut down everyting,
      // including fetching recommendations and check status
      rootState.router!.push({ name: "GoogleSignIn" });
    }

    return response;
  };
}
