"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.machineId = void 0;
const debug_1 = __importDefault(require("debug"));
const node_machine_id_1 = require("node-machine-id");
const crypto_1 = require("crypto");
const log = (0, debug_1.default)("inferable:machine-id");
let cachedId;
const machineId = () => {
    if (cachedId) {
        return cachedId;
    }
    try {
        const id = (0, node_machine_id_1.machineIdSync)();
        const hash = (0, crypto_1.createHash)("shake256", {
            outputLength: 8,
        })
            .update(id)
            .digest("hex");
        cachedId = hash;
        return hash;
    }
    catch (e) {
        log("Failed to get machine ID. Defaulting to random string.", { error: e });
        const str = Math.random().toString(36).substring(8);
        cachedId = str;
        return str;
    }
};
exports.machineId = machineId;
//# sourceMappingURL=machine-id.js.map