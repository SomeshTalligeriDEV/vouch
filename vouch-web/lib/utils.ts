import type { Tier } from "@/types";

export function cn(...classes: (string | false | null | undefined)[]): string {
  return classes.filter(Boolean).join(" ");
}

export function formatCurrency(n: number): string {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 0,
  }).format(n);
}

export function formatScore(n: number): string {
  return new Intl.NumberFormat("en-US", { maximumFractionDigits: 0 }).format(n);
}

export const TIER_COLORS: Record<Tier, string> = {
  Bronze: "text-amber-800 border-amber-700/30 bg-amber-700/10",
  Silver: "text-ink/70 border-ink/20 bg-ink/5",
  Gold: "text-crimson-dark border-gold bg-gold/40",
  Platinum: "text-teal-800 border-teal-700/30 bg-teal-600/15",
  "24 Karat": "text-crimson-dark border-gold bg-gold/60",
};

export function githubOAuthURL(): string {
  const clientId = process.env.NEXT_PUBLIC_GITHUB_CLIENT_ID ?? "";
  const redirect =
    process.env.NEXT_PUBLIC_GITHUB_REDIRECT_URL ??
    "http://localhost:3000/login";
  const params = new URLSearchParams({
    client_id: clientId,
    redirect_uri: redirect,
    scope: "read:user user:email",
  });
  return `https://github.com/login/oauth/authorize?${params.toString()}`;
}
