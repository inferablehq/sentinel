import { randomBytes } from "crypto";
import http, { IncomingMessage, ServerResponse } from "http";
import assert from "node:assert";
import { tokenizeJson } from "./tokenizer";
import { alteredResponses } from "./altered-responses";

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

assert(process.env.DESTINATION_URL, "DESTINATION_URL is required");

const destination = new URL(process.env.DESTINATION_URL);

async function handleRequest(req: IncomingMessage, res: ServerResponse) {
  const requestId = randomBytes(4).toString("hex");

  console.log(`[${requestId}] -> ${req.method} \t ${req.url}`);

  assert(req.url, "Request URL is required");
  const url = new URL(`http://host${req.url}`);

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

    if (!response.body) {
      res.end();
      console.log(`[${requestId}] <- Response sent`);
    } else {
      // do we need to alter the response?
      const alteredResponse = alteredResponses.find((alteredResponse) =>
        alteredResponse.matcher(req)
      );

      if (alteredResponse) {
        const altered = await alteredResponse.handler({
          body: await response.text(),
        });

        res.setHeader("Content-Length", Buffer.byteLength(altered.body));
        res.writeHead(response.status);

        res.write(altered.body);
        res.end();
        console.log(`[${requestId}] <- Altered response sent`);
      } else {
        res.writeHead(response.status);
        const reader = response.body.getReader();
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;
          res.write(value);
        }
        res.end();
        console.log(`[${requestId}] <- Response sent`);
      }
    }
  } catch (error) {
    console.error(`[${requestId}] <- Forwarded response error: ${error}`);
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
