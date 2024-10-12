import { initClient, tsRestFetchApi } from "@ts-rest/core";
import { contract } from "./contract";
const { version: SDK_VERSION } = require("../package.json");

export const createApiClient = ({
  baseUrl,
  machineId,
  clientAbortController,
  apiSecret,
}: {
  baseUrl?: string;
  machineId?: string;
  clientAbortController?: AbortController;
  apiSecret?: string;
}) =>
  initClient(contract, {
    baseUrl: baseUrl ?? "https://api.inferable.ai",
    baseHeaders: {
      "x-machine-sdk-version": SDK_VERSION,
      "x-machine-sdk-language": "typescript",
      ...(apiSecret ? { authorization: apiSecret } : {}),
      ...(machineId ? { "x-machine-id": machineId } : {}),
    },
    api: clientAbortController
      ? (args) => {
          return tsRestFetchApi({
            ...args,
            signal: clientAbortController.signal,
          });
        }
      : undefined,
  });
