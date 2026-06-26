// Token persistence helpers. Builder tokens live in Zustand-backed localStorage.
// Company tokens are stored separately so both sessions can coexist.

import type { TokenPair } from "@/types";

const ACCESS_KEY = "vouch_access_token";
const REFRESH_KEY = "vouch_refresh_token";
const COMPANY_ACCESS_KEY = "vouch_company_access_token";
const COMPANY_REFRESH_KEY = "vouch_company_refresh_token";

// ── Builder (GitHub OAuth) ────────────────────────────────────────────────────

export function getAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(ACCESS_KEY);
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(REFRESH_KEY);
}

export function storeTokens(tokens: TokenPair): void {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(ACCESS_KEY, tokens.access_token);
  window.localStorage.setItem(REFRESH_KEY, tokens.refresh_token);
}

export function clearTokens(): void {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(ACCESS_KEY);
  window.localStorage.removeItem(REFRESH_KEY);
}

// ── Company (email + password) ────────────────────────────────────────────────

export function getCompanyAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(COMPANY_ACCESS_KEY);
}

export function getCompanyRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(COMPANY_REFRESH_KEY);
}

export function storeCompanyTokens(tokens: TokenPair): void {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(COMPANY_ACCESS_KEY, tokens.access_token);
  window.localStorage.setItem(COMPANY_REFRESH_KEY, tokens.refresh_token);
}

export function clearCompanyTokens(): void {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(COMPANY_ACCESS_KEY);
  window.localStorage.removeItem(COMPANY_REFRESH_KEY);
  window.localStorage.removeItem("vouch_company");
}
