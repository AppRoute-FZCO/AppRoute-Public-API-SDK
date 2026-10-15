import { AppRouteError } from "./base.js";
import type { FieldError } from "../models/common.js";

/**
 * Raised when the API returns an error response.
 */
export class ApiError extends AppRouteError {
  public readonly code: string;
  public readonly traceId: string;
  public readonly statusCode: number;
  public readonly errors: FieldError[];
  public readonly apiStatusCode?: number;
  public readonly apiErrorCode?: string;

  constructor(
    code: string,
    message: string,
    traceId: string,
    statusCode: number,
    errors: FieldError[] = [],
    apiStatusCode?: number,
    apiErrorCode?: string,
  ) {
    const parts = [`[${code}] ${message} (traceId=${traceId})`];
    for (const err of errors) {
      parts.push(`  - ${err.field}: ${err.code} — ${err.message}`);
    }
    super(parts.join("\n"));
    this.name = "ApiError";
    this.code = code;
    this.traceId = traceId;
    this.statusCode = statusCode;
    this.errors = errors;
    this.apiStatusCode = apiStatusCode;
    this.apiErrorCode = apiErrorCode;
  }
}
