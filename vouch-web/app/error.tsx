"use client";

import { useEffect } from "react";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Forward to Sentry / any error reporting you have client-side
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-6 px-4 text-center">
      <div className="space-y-2">
        <p className="font-mono text-xs uppercase tracking-widest text-ink/40">
          something went wrong
        </p>
        <h1 className="font-display text-3xl font-bold text-ink">
          Unexpected error
        </h1>
        <p className="text-sm text-ink/60 max-w-sm">
          The page crashed. This has been noted. You can try again or go back home.
        </p>
        {error.digest && (
          <p className="font-mono text-xs text-ink/30">ref: {error.digest}</p>
        )}
      </div>
      <div className="flex gap-3">
        <button
          onClick={reset}
          className="rounded-full bg-accent px-5 py-2 font-bold text-ink shadow-hard transition hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-hard-lg"
        >
          Try again
        </button>
        <a
          href="/"
          className="rounded-full border border-line px-5 py-2 text-sm font-medium text-ink/70 hover:bg-ink/5"
        >
          Go home
        </a>
      </div>
    </div>
  );
}
