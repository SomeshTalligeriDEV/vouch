"use client";

import { useEffect } from "react";
import { useCompanyAuth } from "@/store/company";

/**
 * Hydrates the company auth store from localStorage + API on mount.
 * Call once in a root layout or providers component.
 */
export function useHydrateCompanyAuth() {
  const hydrate = useCompanyAuth((s) => s.hydrate);
  const status = useCompanyAuth((s) => s.status);

  useEffect(() => {
    if (status === "idle") {
      hydrate();
    }
  }, [status, hydrate]);
}

/**
 * Returns true if a company session is authenticated and ready.
 */
export function useIsCompanyAuthenticated(): boolean {
  return useCompanyAuth((s) => s.status === "authenticated" && s.company !== null);
}

/**
 * Returns the authenticated company or null.
 */
export function useCurrentCompany() {
  return useCompanyAuth((s) => s.company);
}
