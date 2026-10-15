import { readFileSync } from "node:fs";
import { describe, it, expect, vi, afterEach } from "vitest";
import { HttpTransport } from "../src/transport/http-transport.js";
import * as errors from "../src/errors/index.js";

const cases = JSON.parse(
  readFileSync(new URL("../../testdata/envelopes.json", import.meta.url), "utf8"),
);

describe("Public API wire envelope", () => {
  afterEach(() => vi.unstubAllGlobals());

  for (const fixture of cases) {
    it(fixture.name, async () => {
      // Keep the real Response JSON decoder and transport; replace only the network.
      vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify(fixture.body), {
        status: fixture.httpStatus,
        headers: { "Content-Type": "application/json" },
      })));
      const transport = new HttpTransport("https://api.test/api/v1", "key", 1000, 0);
      const request = transport.request("GET", "/services");

      if (fixture.invalid) {
        await expect(request).rejects.toThrow(/Invalid.*envelope/);
      } else if ("data" in fixture) {
        await expect(request).resolves.toEqual(fixture.data);
      } else {
        const ErrorClass = errors[(fixture.errorType + "Error") as keyof typeof errors];
        await expect(request).rejects.toBeInstanceOf(ErrorClass);
        await expect(request).rejects.toMatchObject({
          code: fixture.code,
          traceId: fixture.body.traceId,
          statusCode: fixture.httpStatus,
          apiStatusCode: fixture.apiStatusCode,
          apiErrorCode: fixture.apiErrorCode,
          errors: fixture.errors ?? [],
        });
        await expect(request).rejects.toThrow(fixture.message);
      }
    });
  }
});
