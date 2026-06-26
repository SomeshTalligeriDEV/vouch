"use client";

import { useState } from "react";
import { PAGINATION_DEFAULT_LIMIT } from "@/lib/constants";

interface UsePaginationOptions {
  initialPage?: number;
  limit?: number;
}

export function usePagination({ initialPage = 1, limit = PAGINATION_DEFAULT_LIMIT }: UsePaginationOptions = {}) {
  const [page, setPage] = useState(initialPage);

  const goTo = (p: number) => setPage(Math.max(1, p));
  const next = () => setPage((p) => p + 1);
  const prev = () => setPage((p) => Math.max(1, p - 1));
  const reset = () => setPage(1);

  return { page, limit, goTo, next, prev, reset };
}
