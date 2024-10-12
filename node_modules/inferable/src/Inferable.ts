import debug from "debug";
import path from "path";
import { z } from "zod";
import zodToJsonSchema from "zod-to-json-schema";
import { InferableError } from "./errors";
import { FunctionRegistration } from "./types";
import { machineId } from "./machine-id";
import { Service } from "./service";
import {
  FunctionConfig,
  FunctionInput,
  FunctionRegistrationInput,
  JsonSchemaInput,
  RegisteredService,
} from "./types";
import {
  isZodType,
  validateDescription,
  validateFunctionName,
  validateFunctionSchema,
  validateServiceName,
} from "./util";
import * as links from "./links";

// Custom json formatter
debug.formatters.J = (json) => {
  return JSON.stringify(json, null, 2);
};

export const log = debug("inferable:client");

/**
 * The Inferable client. This is the main entry point for using Inferable.
 *
 * Inferable client exposes two main methods:
 * * `service` - Registers a service with Inferable.
 * * `workflow` - Starts a workflow with Inferable
 *
 * @example Basic usage
 * ```ts
 * // src/service.ts
 *
 * // create a new Inferable instance
 * const d = new Inferable({
 *  apiSecret: "API_SECRET",
 * });
 *
 * const myService = d.service({
 *   name: "my-service",
 * });
 *
 * myService.register("hello", z.object({name: z.string()}), async ({name}: {name: string}) => {
 *  return `Hello ${name}`;
 * })
 *
 * await myService.start();
 *
 * // stop the service on shutdown
 * process.on("beforeExit", async () => {
 *   await myService.stop();
 * });
 *
 * ```
 */
export class Inferable {
  static getMachineId(): string {
    return machineId();
  }

  static getVersion(): string {
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    return require(path.join(__dirname, "..", "package.json")).version;
  }

  private apiSecret: string;
  private endpoint: string;
  private machineId: string;

  private services: Service[] = [];

  private functionRegistry: { [key: string]: FunctionRegistration } = {};

  /**
   * Initializes a new Inferable instance.
   * @param apiSecret The API Secret for your Inferable cluster. If not provided, it will be read from the `INFERABLE_API_SECRET` environment variable.
   * @param options Additional options for the Inferable client.
   * @param options.endpoint The endpoint for the Inferable cluster. Defaults to https://api.inferable.ai.
   *
   * @example
   * ```ts
   * // Basic usage
   * const d = new Inferable({
   *  apiSecret: "API_SECRET",
   * });
   *
   * // OR
   *
   * process.env.INFERABLE_API_SECRET = "API_SECRET";
   * const d = new Inferable();
   *
   *
   * // With encryption
   * const d = new Inferable({
   *  encryptionKeys: [
   *    Buffer.from("abcdefghijklmnopqrstuvwxzy123456"), // current key
   *    Buffer.from("abcdefghijklmnopqrstuvwxzy123old"), // previous key
   *  ],
   * });
   * ```
   */
  constructor(options?: {
    apiSecret?: string;
    endpoint?: string;
    jobPollWaitTime?: number;
  }) {
    if (options?.apiSecret && process.env.INFERABLE_API_SECRET) {
      log(
        "API Secret was provided as an option and environment variable. Constructor argument will be used.",
      );
    }

    const apiSecret = options?.apiSecret || process.env.INFERABLE_API_SECRET;

    if (!apiSecret) {
      throw new InferableError(
        `No API Secret provided. Please see ${links.DOCS_AUTH}`,
      );
    }

    if (!apiSecret.startsWith("sk_cluster_machine")) {
      if (apiSecret.startsWith("sk_")) {
        throw new InferableError(
          `Provided non-Machine API Secret. Please see ${links.DOCS_AUTH}`,
        );
      }

      throw new InferableError(
        `Invalid API Secret. Please see ${links.DOCS_AUTH}`,
      );
    }

    this.apiSecret = apiSecret;

    this.endpoint =
      options?.endpoint ||
      process.env.INFERABLE_API_ENDPOINT ||
      "https://api.inferable.ai";
    this.machineId = machineId();
  }

  /**
   * An array containing the name of all services currently polling.
   */
  public get activeServices() {
    return this.services.filter((s) => s.polling).map((s) => s.name);
  }

  /**
   * An array containing the name of all services not currently polling.
   *
   * Note that this will only include services which have been started (`.start()` called).
   */
  public get inactiveServices() {
    return this.services.filter((s) => !s.polling).map((s) => s.name);
  }

  /**
   * An array containing the name of all functions which have been registered.
   */
  public get registeredFunctions() {
    return Object.values(this.functionRegistry).map((f) => f.name);
  }

  /**
   * Convenience reference to a service with name 'default'.
   * @returns A registered service instance.
   * @see {@link service}
   * @example
   * ```ts
   * const d = new Inferable({apiSecret: "API_SECRET"});
   *
   * d.default.register("hello", z.object({name: z.string()}), async ({name}: {name: string}) => {
   *   return `Hello ${name}`;
   * });
   *
   * // start the service
   * await d.default.start();
   *
   * // stop the service on shutdown
   * process.on("beforeExit", async () => {
   *   await d.default.stop();
   * });
   *
   */
  public get default() {
    return this.service({
      name: "default",
    });
  }

  /**
   * Registers a service with Inferable. This will register all functions on the service.
   * @param input The service definition.
   * @returns A registered service instance.
   * @example
   * ```ts
   * const d = new Inferable({apiSecret: "API_SECRET"});
   *
   * const service = d.service({
   *   name: "my-service",
   * });
   *
   * service.register("hello", z.object({name: z.string()}), async ({name}: {name: string}) => {
   *   return `Hello ${name}`;
   * });
   *
   * // start the service
   * await service.start();
   *
   * // stop the service on shutdown
   * process.on("beforeExit", async () => {
   *   await service.stop();
   * });
   * ```
   */
  public service<T extends z.ZodTypeAny | JsonSchemaInput>(input: {
    name: string;
    functions?:
      | FunctionRegistrationInput<T>[]
      | Promise<FunctionRegistrationInput<T>[]>;
  }): RegisteredService {
    validateServiceName(input.name);

    const register: RegisteredService["register"] = ({
      name,
      func,
      schema,
      config,
      description,
      authenticate,
    }) => {
      this.registerFunction({
        name,
        authenticate,
        serviceName: input.name,
        func,
        inputSchema: schema.input,
        config,
        description,
      });
    };

    return {
      definition: input,
      register,
      start: async () => {
        const functions = await input.functions;
        functions?.forEach(register);

        const existing = this.services.find(
          (service) => service.name == input.name,
        );

        if (existing) {
          throw new InferableError(`Service is already started`, {
            serviceName: input.name,
          });
        }

        const serivce = new Service({
          endpoint: this.endpoint,
          machineId: this.machineId,
          apiSecret: this.apiSecret,
          service: input.name,
          functions: Object.values(this.functionRegistry).filter(
            (f) => f.serviceName == input.name,
          ),
        });

        this.services.push(serivce);
        await serivce.start();
      },
      stop: async () => {
        const existing = this.services.find(
          (service) => service.name == input.name,
        );

        if (!existing) {
          throw new InferableError(`Service is not started`, {
            serviceName: input.name,
          });
        }

        await existing.stop();
      },
    };
  }

  /**
   * The cluster ID for this Inferable instance.
   */
  get clusterId(): string | null {
    return this.services[0]?.clusterId || null;
  }

  private registerFunction<T extends z.ZodTypeAny | JsonSchemaInput>({
    name,
    authenticate,
    serviceName,
    func,
    inputSchema,
    config,
    description,
  }: {
    authenticate?: (
      authContext: string,
      args: FunctionInput<T>,
    ) => Promise<void>;
    name: string;
    serviceName: string;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    func: (input: FunctionInput<T>) => any;
    inputSchema: T;
    config?: FunctionConfig;
    description?: string;
  }) {
    if (this.functionRegistry[name]) {
      throw new InferableError(
        `Function name '${name}' is already registered by another service.`,
      );
    }

    // We accept both Zod types and JSON schema as an input, convert to JSON schema if the input is a Zod type
    const inputJson = (
      isZodType(inputSchema) ? zodToJsonSchema(inputSchema) : inputSchema
    ) as JsonSchemaInput;

    validateFunctionName(name);
    validateDescription(description);

    const schemaErrors = validateFunctionSchema(inputJson);

    if (schemaErrors.length > 0) {
      log(
        `Schema ${serviceName}${name} failed validation: %J with failures %O`,
        inputSchema,
        schemaErrors,
      );
      throw new InferableError(
        `JSON schema was not valid for service '${serviceName}.${name}'. Run with debug logging (DEBUG=inferable:*) for more details.`,
        {
          failures: schemaErrors,
        },
      );
    }

    const registration: FunctionRegistration<T> = {
      name,
      authenticate,
      serviceName,
      func,
      schema: {
        input: inputSchema,
        inputJson: JSON.stringify(inputJson),
      },
      config,
      description,
    };

    const existing = this.services.find(
      (service) => service.name == serviceName,
    );

    if (existing) {
      throw new InferableError(
        `Functions must be registered before starting the service. Please see ${links.DOCS_FUNCTIONS}`,
        {
          serviceName: registration.serviceName,
        },
      );
    }

    if (typeof registration.func !== "function") {
      throw new InferableError(
        `func must be a function. Please see ${links.DOCS_FUNCTIONS}`,
      );
    }

    log(`Registering function`, {
      name: registration.name,
    });

    this.functionRegistry[registration.name] = registration;
  }
}
