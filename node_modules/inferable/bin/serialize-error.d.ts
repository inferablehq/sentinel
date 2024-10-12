export function serializeError(value: any, options?: {}): any;
export function isErrorLike(value: any): boolean;
export class NonError extends Error {
    static _prepareSuperMessage(message: any): string;
    constructor(message: any);
}
