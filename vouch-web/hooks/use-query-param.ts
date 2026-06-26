"use client";

import { useCallback } from "react";
import { useRouter, useSearchParams } from "next/navigation";

export function useQueryParam(key: string) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const value = searchParams.get(key) ?? "";

  const setValue = useCallback(
    (next: string) => {
      const params = new URLSearchParams(searchParams.toString());
      if (next) {
        params.set(key, next);
      } else {
        params.delete(key);
      }
      router.replace(`?${params.toString()}`, { scroll: false });
    },
    [key, router, searchParams],
  );

  return [value, setValue] as const;
}
