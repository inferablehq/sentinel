import { z } from 'zod';
import { AppRoute } from './dsl';
export declare class ResponseValidationError extends Error {
    appRoute: AppRoute;
    cause: z.ZodError;
    constructor(appRoute: AppRoute, cause: z.ZodError);
}
