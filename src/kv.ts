import fs from "fs";
import path from "path";

export type KVValue =
  | {
      type: "string";
      value: string;
    }
  | {
      type: "number";
      value: number;
    };

interface IKVStore {
  set(key: string, value: KVValue): void;
  get(key: string): KVValue | null;
}

export class KVStore implements IKVStore {
  private dataDir: string;

  constructor() {
    this.dataDir = path.join(process.cwd(), "data");
    this.ensureDataDirExists();
  }

  private ensureDataDirExists(): void {
    if (!fs.existsSync(this.dataDir)) {
      fs.mkdirSync(this.dataDir, { recursive: true });
    }
  }

  private getFilePath(key: string): string {
    return path.join(this.dataDir, `${key}`);
  }

  set(key: string, value: KVValue): void {
    const filePath = this.getFilePath(key);
    fs.writeFileSync(filePath, JSON.stringify(value), "utf8");
  }

  get(key: string): KVValue | null {
    const filePath = this.getFilePath(key);
    if (fs.existsSync(filePath)) {
      const data = fs.readFileSync(filePath, "utf8");
      const parsed = JSON.parse(data) as KVValue;
      return parsed;
    }
    return null;
  }
}
