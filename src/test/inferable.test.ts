import { Inferable } from "inferable";
import { z } from "zod";
import { test } from "node:test";
import { createServer } from "../server";
import http from "node:http";
import assert from "node:assert";

test("should be able to register a machine", async (t) => {
  assert(
    process.env.INFERABLE_MACHINE_API_SECRET,
    "Machine API secret is not set"
  );

  const client = new Inferable({
    apiSecret: process.env.INFERABLE_MACHINE_API_SECRET,
    endpoint: process.env.SENTINEL_EXTERNAL_URL,
  });

  const address = {
    line1: "123 Main St",
    line2: "Apt 4B",
    city: "New York",
    state: "NY",
    zip: "10001",
  };

  client.default.register({
    name: "getUserDetails",
    description: "Gets user details",
    schema: {
      input: z.object({
        userId: z.string(),
      }),
    },
    func: async (input: { userId: string }) => {
      return {
        id: input.userId,
        name: "Not John Smith",
        age: 30,
        email: "notjohnsmith@example.com",
        addresses: [address],
      };
    },
  });

  client.default.register({
    name: "updateUserZip",
    description: "Updates user zip code",
    schema: {
      input: z.object({
        userId: z.string(),
        zip: z.string(),
      }),
    },
    func: async (input: { userId: string; zip: string }) => {
      address.zip = input.zip;

      return address;
    },
  });

  await client.default.start();

  assert(client.clusterId, "Cluster ID is not set");
  assert(
    process.env.INFERABLE_CLUSTER_CONSUME_SECRET,
    "Cluster consume secret is not set"
  );

  const response = await fetch(
    `${process.env.SENTINEL_EXTERNAL_URL}/clusters/${client.clusterId}/runs`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${process.env.INFERABLE_CLUSTER_CONSUME_SECRET}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        message: JSON.stringify({
          prompts:
            "Check the user's zip code, if it matches the old value, update it to the new value.",
          data: {
            userId: "123",
            oldZip: "10001",
            newZip: "10002",
          },
        }),
      }),
    }
  );

  const data: any = await response.json();

  assert(data.id, "Run ID is not set");

  // await client.default.stop();
  // await stopServer(server);
});
