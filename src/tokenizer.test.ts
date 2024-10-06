import { test } from "node:test";
import assert from "node:assert";
import { tokenizeJson } from "./tokenizer";

test("tokenizeJson", async (t) => {
  await t.test("should tokenize simple string", () => {
    const input = "Hello, world!";
    const tokenized = tokenizeJson(JSON.stringify(input));
    const parsed: any = JSON.parse(tokenized);

    assert.strictEqual(typeof parsed, "string");
    assert.strictEqual(parsed.length, 16); // Hex token length
  });

  await t.test("should tokenize number", () => {
    const input = 42;
    const tokenized = tokenizeJson(JSON.stringify(input));
    const parsed: any = JSON.parse(tokenized);

    assert.strictEqual(typeof parsed, "string");
    assert.strictEqual(parsed.length, 16);
  });

  await t.test("should tokenize array", () => {
    const input = ["a", 1, "b"];
    const tokenized = tokenizeJson(JSON.stringify(input));
    const parsed: any = JSON.parse(tokenized);

    assert.strictEqual(Array.isArray(parsed), true);
    assert.strictEqual(parsed.length, 3);
    parsed.forEach((item: any) => {
      assert.strictEqual(typeof item, "string");
      assert.strictEqual(item.length, 16);
    });
  });

  await t.test("should tokenize object", () => {
    const input = { name: "Alice", age: 30 };
    const tokenized = tokenizeJson(JSON.stringify(input));
    const parsed: any = JSON.parse(tokenized);

    assert.strictEqual(typeof parsed, "object");
    assert.strictEqual(parsed !== null, true);
    assert.strictEqual(Object.keys(parsed).length, 2);
    assert.strictEqual(typeof parsed.name, "string");
    assert.strictEqual(parsed.name.length, 16);
    assert.strictEqual(typeof parsed.age, "string");
    assert.strictEqual(parsed.age.length, 16);
  });

  await t.test("should handle nested structures", () => {
    const input = {
      user: {
        name: "Bob",
        scores: [85, 90, 95],
      },
      active: true,
    };
    const tokenized = tokenizeJson(JSON.stringify(input));
    const parsed: any = JSON.parse(tokenized);

    assert.strictEqual(typeof parsed, "object");
    assert.strictEqual(parsed !== null, true);
    assert.strictEqual(Object.keys(parsed).length, 2);
    assert.strictEqual(typeof parsed.user, "object");
    assert.strictEqual(parsed.user !== null, true);
    assert.strictEqual(typeof parsed.user.name, "string");
    assert.strictEqual(parsed.user.name.length, 16);
    assert.strictEqual(Array.isArray(parsed.user.scores), true);
    assert.strictEqual(parsed.user.scores.length, 3);
    parsed.user.scores.forEach((score: any) => {
      assert.strictEqual(typeof score, "string");
      assert.strictEqual(score.length, 16);
    });
    assert.strictEqual(typeof parsed.active, "boolean");
    assert.strictEqual(parsed.active, true);
  });

  await t.test("should handle very deep JSON structure", () => {
    const createDeepObject = (depth: number): any => {
      if (depth === 0) {
        return {
          string: "deepest level",
          number: 42,
          boolean: true,
          null: null,
          array: [134, "two", false],
        };
      }
      return {
        nested: createDeepObject(depth - 1),
      };
    };

    const input = createDeepObject(100); // Create a 100-level deep object
    const tokenized = tokenizeJson(JSON.stringify(input));

    assert.equal(tokenized.includes("deepest level"), false);
    assert.equal(tokenized.includes("42,"), false);
    assert.equal(tokenized.includes("134,"), false);
    assert.equal(tokenized.includes("two"), false);

    assert.equal(tokenized.includes("false"), true);
    assert.equal(tokenized.includes("true"), true);
    assert.equal(tokenized.includes("null"), true);
  });
});
