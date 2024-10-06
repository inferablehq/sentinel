import { Inferable } from "inferable";
import { z } from "zod";
import { test } from "node:test";
import { createServer } from "../server";
import http from "node:http";

const startServer = async (port: number) => {
  const server = createServer();
  await new Promise<void>((resolve) => {
    server.listen(port, () => {
      console.log(`Test server running on port ${port}`);
      resolve();
    });
  });
  return server;
};

const stopServer = async (server: http.Server) => {
  await new Promise<void>((resolve) => server.close(() => resolve()));
};

test("should be able to register a machine", async () => {
  const port = 3195;

  const server = await startServer(port);

  const client = new Inferable({
    apiSecret: process.env.INFERABLE_API_SECRET,
    endpoint: `http://localhost:${port}`,
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

  const response = await fetch(
    `http://localhost:${port}/clusters/${process.env.INFERABLE_CLUSTER_ID}/runs`,
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

  const data = await response.json();

  // TODO: assert that the response is correct

  await stopServer(server);
});
