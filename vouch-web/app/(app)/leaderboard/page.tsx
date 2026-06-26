"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { formatScore } from "@/lib/format";
import { TIER_COLORS, TIER_LABELS } from "@/lib/constants";
import { ScoreRing } from "@/components/ui/score-ring";
import { TierBadge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar } from "@/components/ui/avatar";
import Link from "next/link";
import type { BuilderScore } from "@/types";

export default function LeaderboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["leaderboard"],
    queryFn: () => api.leaderboard(50),
  });

  return (
    <main className="max-w-3xl mx-auto px-4 py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Leaderboard</h1>
        <p className="text-muted-foreground mt-1">Top builders ranked by Vouch score</p>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 10 }).map((_, i) => (
            <Skeleton key={i} className="h-16 rounded-xl" />
          ))}
        </div>
      ) : !data?.length ? (
        <div className="rounded-xl border border-dashed border-border p-16 text-center text-muted-foreground">
          No builders on the leaderboard yet.
        </div>
      ) : (
        <ol className="space-y-3">
          {data.map((entry: BuilderScore, idx: number) => (
            <LeaderboardRow key={entry.user_id} entry={entry} rank={idx + 1} />
          ))}
        </ol>
      )}
    </main>
  );
}

function LeaderboardRow({ entry, rank }: { entry: BuilderScore; rank: number }) {
  const tier = entry.tier ?? "bronze";
  const color = TIER_COLORS[tier as keyof typeof TIER_COLORS] ?? "#CD7F32";

  return (
    <li>
      <Link
        href={`/u/${entry.username}`}
        className="flex items-center gap-4 rounded-xl border border-border bg-card px-4 py-3 hover:border-primary/40 hover:bg-accent/30 transition-colors"
      >
        <span
          className="w-7 text-right text-sm font-bold tabular-nums"
          style={{ color: rank <= 3 ? color : undefined }}
        >
          {rank <= 3 ? ["🥇", "🥈", "🥉"][rank - 1] : `#${rank}`}
        </span>
        <ScoreRing score={entry.total_score} tier={tier} size={44} />
        <div className="flex-1 min-w-0">
          <p className="font-semibold truncate">{entry.username}</p>
          <TierBadge tier={tier} className="mt-0.5" />
        </div>
        <div className="text-right">
          <p className="text-lg font-bold tabular-nums" style={{ color }}>
            {formatScore(entry.total_score)}
          </p>
          <p className="text-xs text-muted-foreground">{TIER_LABELS[tier as keyof typeof TIER_LABELS]}</p>
        </div>
      </Link>
    </li>
  );
}
