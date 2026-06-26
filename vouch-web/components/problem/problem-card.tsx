"use client";

import type { Problem } from "@/types";
import { formatCurrency } from "@/lib/utils";
import { useUpvoteProblem } from "@/hooks";

const STATUS_STYLES: Record<Problem["status"], string> = {
  open: "border-emerald-500/40 text-emerald-400",
  claimed: "border-amber-500/40 text-amber-400",
  shipped: "border-crimson/40 text-crimson",
  cancelled: "border-line text-ink/50",
};

export function ProblemCard({ problem }: { problem: Problem }) {
  const upvote = useUpvoteProblem();

  return (
    <div className="card flex gap-4">
      <button
        onClick={() => upvote.mutate(problem.id)}
        disabled={upvote.isPending}
        className="flex h-14 w-12 shrink-0 flex-col items-center justify-center rounded-lg border border-line text-sm hover:border-crimson/40"
        aria-label="Upvote"
      >
        <span className="text-crimson">▲</span>
        <span className="font-semibold">{problem.upvotes}</span>
      </button>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <h3 className="truncate font-semibold">{problem.title}</h3>
          <span className={`tag ${STATUS_STYLES[problem.status]}`}>
            {problem.status}
          </span>
        </div>
        <p className="mt-1 line-clamp-2 text-sm text-ink/60">
          {problem.description}
        </p>
        <div className="mt-3 text-sm text-ink/80">
          Budget: {formatCurrency(problem.budget_min)} –{" "}
          {formatCurrency(problem.budget_max)}
        </div>
      </div>
    </div>
  );
}
