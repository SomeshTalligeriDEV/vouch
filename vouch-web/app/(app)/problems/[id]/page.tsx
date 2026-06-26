"use client";

import { use, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuth } from "@/store/auth";
import { formatRelativeTime } from "@/lib/format";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";

interface Props {
  params: Promise<{ id: string }>;
}

const STATUS_COLOR: Record<string, string> = {
  open: "bg-emerald-500/10 text-emerald-400 border-emerald-500/20",
  claimed: "bg-yellow-500/10 text-yellow-400 border-yellow-500/20",
  shipped: "bg-purple-500/10 text-purple-400 border-purple-500/20",
  cancelled: "bg-muted text-muted-foreground",
};

export default function ProblemDetailPage({ params }: Props) {
  const { id } = use(params);
  const user = useAuth((s) => s.user);
  const qc = useQueryClient();
  const [claiming, setClaiming] = useState(false);
  const [upvoting, setUpvoting] = useState(false);
  const [claimError, setClaimError] = useState<string | null>(null);

  const { data: problem, isLoading } = useQuery({
    queryKey: ["problem", id],
    queryFn: () => api.getProblem(id),
    enabled: !!id,
  });

  const handleClaim = async () => {
    setClaiming(true);
    setClaimError(null);
    try {
      await api.claimProblem(id);
      qc.invalidateQueries({ queryKey: ["problem", id] });
    } catch (err) {
      setClaimError((err as Error).message);
    } finally {
      setClaiming(false);
    }
  };

  const handleUpvote = async () => {
    setUpvoting(true);
    try {
      await api.upvoteProblem(id);
      qc.invalidateQueries({ queryKey: ["problem", id] });
    } finally {
      setUpvoting(false);
    }
  };

  if (isLoading) {
    return (
      <main className="max-w-2xl mx-auto px-4 py-10 space-y-4">
        <Skeleton className="h-8 w-3/4" />
        <Skeleton className="h-4 w-40" />
        <Skeleton className="h-32 rounded-xl" />
      </main>
    );
  }

  if (!problem) {
    return (
      <main className="max-w-2xl mx-auto px-4 py-20 text-center text-muted-foreground">
        Problem not found.
      </main>
    );
  }

  return (
    <main className="max-w-2xl mx-auto px-4 py-10 space-y-6">
      {/* Back */}
      <Link href="/problems" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
        ← Demand board
      </Link>

      {/* Header */}
      <div className="rounded-xl border border-border bg-card p-6 space-y-4">
        <div className="flex items-start gap-3">
          <div className="flex-1">
            <div className="flex flex-wrap items-center gap-2 mb-2">
              <span
                className={`text-xs font-medium rounded-full border px-2.5 py-0.5 ${STATUS_COLOR[problem.status] ?? ""}`}
              >
                {problem.status}
              </span>
              <span className="text-xs text-muted-foreground">
                {formatRelativeTime(problem.created_at)}
              </span>
            </div>
            <h1 className="text-2xl font-bold">{problem.title}</h1>
          </div>
        </div>

        <p className="text-sm leading-relaxed whitespace-pre-line">{problem.description}</p>

        {/* Tags */}
        {problem.tags?.length > 0 && (
          <div className="flex flex-wrap gap-1.5">
            {problem.tags.map((tag: string) => (
              <Badge key={tag} variant="secondary">{tag}</Badge>
            ))}
          </div>
        )}

        {/* Budget */}
        {problem.budget_max > 0 && (
          <div className="rounded-lg bg-muted/50 p-3 text-sm">
            <span className="font-medium">Budget: </span>
            <span className="text-foreground">
              ${problem.budget_min.toLocaleString()} – ${problem.budget_max.toLocaleString()} USD
            </span>
          </div>
        )}

        {/* Actions */}
        <div className="flex flex-wrap gap-3 pt-2">
          <button
            onClick={handleUpvote}
            disabled={upvoting}
            className="flex items-center gap-1.5 rounded-lg border border-border bg-background px-3 py-1.5 text-sm hover:border-primary/40 hover:bg-accent/30 transition-colors"
          >
            ▲ {problem.upvotes} upvotes
          </button>

          {user && problem.status === "open" && (
            <button
              onClick={handleClaim}
              disabled={claiming}
              className="rounded-lg bg-primary px-4 py-1.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
            >
              {claiming ? "Claiming…" : "Claim this problem"}
            </button>
          )}

          {!user && problem.status === "open" && (
            <Link
              href="/login"
              className="rounded-lg bg-primary px-4 py-1.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
            >
              Sign in to claim
            </Link>
          )}

          {problem.status === "shipped" && problem.shipped_project_id && (
            <Link
              href={`/projects/${problem.shipped_project_id}`}
              className="rounded-lg border border-purple-500/30 bg-purple-500/10 px-4 py-1.5 text-sm font-medium text-purple-400 hover:bg-purple-500/20 transition-colors"
            >
              View solution →
            </Link>
          )}
        </div>

        {claimError && (
          <p className="text-sm text-destructive">{claimError}</p>
        )}

        {problem.status === "claimed" && (
          <p className="text-sm text-yellow-400 bg-yellow-500/10 rounded-lg px-3 py-2">
            A builder has claimed this problem and is working on a solution.
          </p>
        )}
      </div>

      {/* CTA for companies */}
      {!user && (
        <div className="rounded-xl border border-dashed border-border p-6 text-center space-y-3">
          <p className="text-sm font-medium">Are you a company with a similar problem?</p>
          <p className="text-xs text-muted-foreground">
            Post your own version and let builders compete to solve it.
          </p>
          <Link
            href="/company/register"
            className="inline-block rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
          >
            Post a problem
          </Link>
        </div>
      )}
    </main>
  );
}
