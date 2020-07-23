export function extractFromResource(
  property: string,
  resource: string
): string {
  const sliceLen = `/${property}/`.length;
  const pattern = `/${property}/[^/]*`;
  const regex = new RegExp(pattern);

  const found = regex.exec(resource);
  if (found === null) throw `couldn't parse resource identifier: ${resource}`;

  const result = found[0].slice(sliceLen);
  return result;
}
