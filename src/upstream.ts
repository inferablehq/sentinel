import { KVStore } from "./kv";
import assert from "node:assert";
const SQS_PATH = "/sqs";

export const upstreamUrlStore = new KVStore();

export const getUpstreamUrl = (path: string): string | null => {
  if (path.startsWith(SQS_PATH)) {
    const key = path.split(SQS_PATH)[1].replace(/\//g, "");
    console.log("key", key, path);
    const queueUrl = upstreamUrlStore.get(key);
    if (queueUrl) {
      return queueUrl.value as string;
    } else {
      throw new Error(`Queue URL not found for key: ${key}`);
    }
  }

  return null;
};
