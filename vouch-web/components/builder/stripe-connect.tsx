"use client";

function stripeConnectURL(): string {
  const clientId = process.env.NEXT_PUBLIC_STRIPE_CLIENT_ID ?? "";
  const redirect =
    process.env.NEXT_PUBLIC_STRIPE_REDIRECT_URL ??
    "http://localhost:3000/dashboard";
  const params = new URLSearchParams({
    response_type: "code",
    client_id: clientId,
    scope: "read_only",
    redirect_uri: redirect,
  });
  return `https://connect.stripe.com/oauth/authorize?${params.toString()}`;
}

/**
 * StripeConnect starts the read-only Stripe OAuth flow. Connecting verifies a
 * builder's revenue and lifts their score multiplier from 0.6 → 1.0.
 */
export function StripeConnect({ connected }: { connected?: boolean }) {
  if (connected) {
    return (
      <span className="inline-flex items-center gap-2 rounded-lg border border-emerald-500/40 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-400">
        ✓ Stripe verified — full score multiplier
      </span>
    );
  }
  return (
    <a href={stripeConnectURL()} className="btn-primary">
      Connect Stripe (read-only)
    </a>
  );
}
