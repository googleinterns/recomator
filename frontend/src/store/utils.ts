const DEVELOPMENT_SERVER_ADDRESS = "http://localhost:8000";
const PRODUCTION_SERVER_ADDRESS = "";

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
    return DEVELOPMENT_SERVER_ADDRESS;
  }

  return PRODUCTION_SERVER_ADDRESS;
}
