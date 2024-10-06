import { createServer } from "./server";
import http from "http";
import assert from "assert";
import { test } from "node:test";

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

const tests = [
  {
    name: "should respond with 200 OK for /live endpoint",
    port: 3195,
    path: "/live",
    expectedStatus: 200,
    expectedBody: { status: "ok" },
  },
  {
    name: "should respond with 404 for /nonexistent endpoint",
    port: 3195,
    path: "/nonexistent",
    expectedStatus: 404,
  },
];

test("should respond with 200 OK for /live endpoint", async () => {
  const port = 3195;
  const server = await startServer(port);
  try {
    for (const test of tests) {
      const response = await fetch(`http://localhost:${test.port}${test.path}`);
      assert.strictEqual(response.status, test.expectedStatus);
      const body = await response.json();
      if (test.expectedBody) {
        assert.deepStrictEqual(body, test.expectedBody);
      }
    }
  } finally {
    await stopServer(server);
  }
});
