import assert from "node:assert";
import { IncomingMessage } from "node:http";
import { KVStore } from "./kv";

const kv = new KVStore();

export const alteredResponses = [
  {
    matcher: (req: IncomingMessage) =>
      req.url === "/machines" && req.method === "POST",
    handler: async ({ body }: { body: string }): Promise<{ body: string }> => {
      assert(body, "Body is required");
      assert(typeof body === "string", "Body must be a string");

      const parsedBody = JSON.parse(body);

      if ("queueUrl" in parsedBody) {
        const queueUrl = parsedBody.queueUrl;
        assert(typeof queueUrl === "string", "queueUrl must be a string");
        kv.set("queueUrl", {
          type: "string",
          value: queueUrl,
        });

        assert(
          process.env.SENTINEL_EXTERNAL_URL,
          "SENTINEL_EXTERNAL_URL is not set"
        );
        parsedBody.queueUrl = `${process.env.SENTINEL_EXTERNAL_URL}/sqs`;

        return {
          body: JSON.stringify(parsedBody),
        };
      } else {
        throw new Error(`Unexpected body: ${JSON.stringify(parsedBody)}`);
      }
    },
  },
];
