import { createServer } from "./server";

const port = 3194;

const server = createServer();

process.on("uncaughtException", (err) => {
  console.error(err);
});

server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}/`);
});

// Export the server
export { server };
