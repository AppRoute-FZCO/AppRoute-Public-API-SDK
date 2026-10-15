/**
 * Field-level error detail returned in API validation responses.
 */
export interface FieldError {
  field: string;
  code: string;
  message: string;
}

/**
 * Standard API response envelope wrapping all responses.
 */
export interface Envelope<T = unknown> {
  status: string;
  statusCode?: number;
  statusMessage?: string;
  /** Legacy envelope fields, present only on older API deployments. */
  code?: string;
  message?: string;
  traceId?: string;
  data?: T;
  errors?: FieldError[] | null;
  errorCode?: string | null;
}
