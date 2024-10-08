import { randomBytes } from "crypto";
import assert from "node:assert";
import { KVStore, KVValue } from "./kv";

const KV = new KVStore();

function getToken(value: KVValue) {
  const t = randomBytes(8).toString("hex");
  KV.set(t, value);
  return t;
}

function replaceValues(obj: any): any {
  if (typeof obj === "string") {
    return getToken({
      type: "string",
      value: obj,
    });
  } else if (typeof obj === "number") {
    return getToken({
      type: "number",
      value: obj,
    });
  } else if (Array.isArray(obj)) {
    return obj.map(replaceValues);
  } else if (typeof obj === "object" && obj !== null) {
    const newObj: { [key: string]: any } = {};
    for (const key in obj) {
      newObj[key] = replaceValues(obj[key]);
    }
    return newObj;
  }
  return obj;
}

export const tokenizeJson = (
  json: string,
  unmaskable?: string | null
): string => {
  assert(typeof json === "string", "JSON must be a string");

  const parsed = JSON.parse(json);

  let unmasked: { [key: string]: unknown } = {};

  if (unmaskable) {
    const unmaskableKeys = unmaskable.split(",");

    for (const key of unmaskableKeys) {
      unmasked[key] = parsed[key];
      delete parsed[key];
    }
  }

  const masked = replaceValues(parsed);

  return JSON.stringify(Object.assign(unmasked, masked));
};
