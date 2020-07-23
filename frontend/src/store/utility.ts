export function extractFromResource(
  property: string,
  resource: string
): string {
  const sliceLen = `/${property}/`.length;
  const pattern = `/${property}/[^/]*`;
  const regex = new RegExp(pattern);

  const found = regex.exec(resource);
  if (found === null) {
    return "";
  }

  const result = found[0].slice(sliceLen);
  console.log(result);

  return result;
}
