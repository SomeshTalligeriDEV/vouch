"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";

export function useRequireAuth(redirectTo = "/") {
  const router = useRouter();
  const { status } = useAuthStore();

  useEffect(() => {
    if (status === "unauthenticated") {
      router.replace(redirectTo);
    }
  }, [status, router, redirectTo]);

  return { isLoading: status === "loading", isAuthenticated: status === "authenticated" };
}

export function useRedirectIfAuthenticated(redirectTo = "/dashboard") {
  const router = useRouter();
  const { status } = useAuthStore();

  useEffect(() => {
    if (status === "authenticated") {
      router.replace(redirectTo);
    }
  }, [status, router, redirectTo]);

  return { isLoading: status === "loading" };
}
