<p align="center">
  <img src="https://a.inferable.ai/logo.png" width="200" style="border-radius: 10px" />
</p>

# Typescript SDK

This is the official Inferable AI SDK for Typescript.

## Installation

### npm

```bash
npm install inferable
```

### yarn

```bash
yarn add inferable
```

### pnpm

```bash
pnpm add inferable
```

## Quick Start

### 1. Initializing Inferable

Create a file named i.ts which will be used to initialize Inferable. This file will export the Inferable instance.

```typescript
// d.ts

import { Inferable } from "inferable";

// Initialize the Inferable client with your API secret.
// Get yours at https://console.inferable.ai.
export const d = new Inferable({
  apiSecret: "YOUR_API_SECRET",
});
```

### 2. Hello World Service

In a separate file, create the "Hello World" service. This file will import the Inferable instance from i.ts and define the service.

```typescript
// service.ts

import { i } from "./i";

// Define a simple function that returns "Hello, World!"
const sayHello = async ({ to }: { to: string }) => {
  return `Hello, ${to}!`;
};

// Create the service
export const helloWorldService = d.service({
  name: "helloWorld",
});

helloWorldService.register({
  name: "sayHello",
  func: sayHello,
  schema: {
    input: z.object({
      to: z.string(),
    }),
  },
});
```

### 3. Running the Service

To run the service, simply run the file with the service definition. This will start the service and make it available to your Inferable agent.

```bash
tsx service.ts
```

## Documentation

- [Inferable documentation](https://docs.inferable.ai/) contains all the information you need to get started with Inferable.
