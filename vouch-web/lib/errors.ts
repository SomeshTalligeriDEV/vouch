export class ApiError extends Error {
  constructor(
    public readonly code: string,
    message: string,
    public readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }

  isUnauthorized() {
    return this.status === 401;
  }

  isForbidden() {
    return this.status === 403;
  }

  isNotFound() {
    return this.status === 404;
  }

  isRateLimited() {
    return this.status === 429;
  }

  isServerError() {
    return this.status >= 500;
  }
}

export function isApiError(err: unknown): err is ApiError {
  return err instanceof ApiError;
}

export function friendlyMessage(err: unknown): string {
  if (isApiError(err)) {
    if (err.isUnauthorized()) return "Please sign in to continue.";
    if (err.isForbidden()) return "You don't have permission to do that.";
    if (err.isNotFound()) return "That resource doesn't exist.";
    if (err.isRateLimited()) return "Too many requests — please slow down.";
    if (err.isServerError()) return "Something went wrong on our end. Try again shortly.";
    return err.message;
  }
  if (err instanceof Error) return err.message;
  return "An unexpected error occurred.";
}
