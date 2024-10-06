import { randomBytes } from "crypto";
import http, { IncomingMessage, ServerResponse } from "http";
import assert from "node:assert";
import { tokenizeJson } from "./tokenizer";

export function createServer() {
  return http.createServer(
    async (req: IncomingMessage, res: ServerResponse) => {
      try {
        await handleRequest(req, res);
      } catch (error) {
        console.error("Unhandled error:", error);
        res.statusCode = 500;
        res.end("Internal Server Error");
      }
    }
  );
}

async function handleRequest(req: IncomingMessage, res: ServerResponse) {
  const requestId = randomBytes(4).toString("hex");

  console.log(`[${requestId}] -> ${req.method} \t ${req.url}`);

  assert(req.url, "Request URL is required");
  const url = new URL(`http://host${req.url}`);

  // Handle the /live endpoint
  if (url.pathname === "/live") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ status: "ok" }));
    return;
  }

  const destination = new URL(
    process.env.DESTINATION_URL || "https://api.inferable.ai"
  );

  const body = await getRequestBody(req);
  let processedBody = body;

  const headers = new Headers();

  for (const [key, value] of Object.entries(req.headers)) {
    headers.set(key, value as string);
  }

  const noMask = headers.get("x-sentinel-no-mask");

  if (noMask) {
    console.log(`[${requestId}] -> No mask`);
  }

  if (body && !noMask) {
    const unmaskable = headers.get("x-sentinel-unmask-keys");

    if (unmaskable) {
      console.log(`[${requestId}] -> Unmaskable: ${unmaskable}`);
    }

    try {
      const contentType = headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        console.log(`[${requestId}] -> Tokenizing JSON`);
        processedBody = tokenizeJson(body, unmaskable);
        headers.set("Content-Type", "application/json");
      } else {
        console.log(`[${requestId}] -> Content-Type: ${contentType}`);
        headers.set("Content-Type", contentType || "text/plain");
      }
      headers.set(
        "Content-Length",
        Buffer.byteLength(processedBody).toString()
      );
    } catch (error) {
      console.error(`[${requestId}] -> Error tokenizing JSON: ${error}`);
      res.statusCode = 400;
      res.end(error instanceof Error ? error.message : String(error));
      return;
    }
  }

  const targetUrl = new URL(url.pathname + url.search, destination);

  try {
    console.log(`[${requestId}] -> Fetching ${targetUrl}`);

    const response = await fetch(targetUrl, {
      method: req.method,
      headers: headers,
      body: processedBody || undefined,
    });

    console.log(`[${requestId}] <- Response from upstream: ${response.status}`);

    // Forward the response headers
    for (const [key, value] of response.headers.entries()) {
      res.setHeader(key, value);
    }

    // Forward the response status
    res.writeHead(response.status);

    // Stream the response body
    if (response.body) {
      const reader = response.body.getReader();
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        res.write(value);
      }
    }

    console.log(`[${requestId}] <- Response sent`);

    res.end();
  } catch (error) {
    console.error(`[${requestId}] <- Forwarded request error: ${error}`);
    res.statusCode = 500;
    res.end(`Error: ${error instanceof Error ? error.message : String(error)}`);
  }

  console.log(`[${requestId}] <- End`);
}

async function getRequestBody(req: IncomingMessage): Promise<string> {
  const chunks: Buffer[] = [];
  for await (const chunk of req) {
    chunks.push(chunk);
  }
  return Buffer.concat(chunks).toString();
}
