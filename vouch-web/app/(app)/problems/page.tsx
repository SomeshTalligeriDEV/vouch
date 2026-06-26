"use client";

import { useState } from "react";

import { ProblemCard } from "@/components/problem/problem-card";
import { useProblems } from "@/hooks";

export default function ProblemsPage() {
  const [status, setStatus] = useState("");
  const { data, isLoading, isError } = useProblems({ status, sort: "upvotes" });

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">Demand board</h1>
          <p className="text-sm text-ink/60">
            Real problems with budgets. Claim one and become the first paid user.
          </p>
        </div>
        <select
          value={status}
          onChange={(e) => setStatus(e.target.value)}
          className="rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-crimson"
        >
          <option value="">All</option>
          <option value="open">Open</option>
          <option value="claimed">Claimed</option>
          <option value="shipped">Shipped</option>
        </select>
      </div>

      {isLoading && <p className="text-ink/60">Loading…</p>}
      {isError && <p className="text-red-400">Failed to load problems.</p>}
      {data && data.items.length === 0 && (
        <p className="text-ink/60">No problems posted yet.</p>
      )}

      <div className="grid gap-4">
        {data?.items.map((p) => (
          <ProblemCard key={p.id} problem={p} />
        ))}
      </div>
    </div>
  );
}
