"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { formatScore, formatRelativeTime } from "@/lib/format";
import { TIER_COLORS, TIER_LABELS } from "@/lib/constants";
import { ScoreRing } from "@/components/ui/score-ring";
import { TierBadge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar } from "@/components/ui/avatar";
import { use } from "react";
import Link from "next/link";
import type { Project } from "@/types";

interface Props {
  params: Promise<{ username: string }>;
}

export default function BuilderProfilePage({ params }: Props) {
  const { username } = use(params);

  const { data: user, isLoading: userLoading } = useQuery({
    queryKey: ["user", username],
    queryFn: () => api.getUser(username),
  });

  const { data: score, isLoading: scoreLoading } = useQuery({
    queryKey: ["score", username],
    queryFn: () => api.getScore(username),
    enabled: !!username,
  });

  const { data: projects } = useQuery({
    queryKey: ["projects", username],
    queryFn: () => api.listProjects({ limit: 6 }),
    enabled: !!username,
  });

  const isLoading = userLoading || scoreLoading;

  if (isLoading) {
    return (
      <main className="max-w-3xl mx-auto px-4 py-10 space-y-6">
        <div className="flex items-center gap-6">
          <Skeleton className="w-20 h-20 rounded-full" />
          <div className="space-y-2 flex-1">
            <Skeleton className="h-6 w-40" />
            <Skeleton className="h-4 w-64" />
          </div>
        </div>
        <Skeleton className="h-32 rounded-xl" />
      </main>
    );
  }

  if (!user) {
    return (
      <main className="max-w-3xl mx-auto px-4 py-20 text-center text-muted-foreground">
        Builder not found.
      </main>
    );
  }

  const tier = score?.tier ?? "bronze";
  const tierColor = TIER_COLORS[tier as keyof typeof TIER_COLORS] ?? "#CD7F32";

  return (
    <main className="max-w-3xl mx-auto px-4 py-10">
      {/* Header */}
      <div className="flex items-start gap-6 mb-8">
        <Avatar src={user.avatar_url} name={user.display_name || user.username} size={80} />
        <div className="flex-1 min-w-0">
          <h1 className="text-2xl font-bold">{user.display_name || user.username}</h1>
          <p className="text-muted-foreground">@{user.username}</p>
          {user.bio && <p className="mt-2 text-sm">{user.bio}</p>}
          <div className="flex flex-wrap gap-2 mt-3">
            {user.location && (
              <span className="text-xs text-muted-foreground">📍 {user.location}</span>
            )}
            {user.website && (
              <a
                href={user.website}
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs text-primary hover:underline"
              >
                🔗 Website
              </a>
            )}
            {user.github_username && (
              <a
                href={`https://github.com/${user.github_username}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs text-muted-foreground hover:text-foreground"
              >
                GitHub →
              </a>
            )}
          </div>
        </div>
        {score && (
          <div className="text-right">
            <ScoreRing score={score.total_score} tier={tier} size={72} />
            <TierBadge tier={tier} className="mt-1" />
          </div>
        )}
      </div>

      {/* Score breakdown */}
      {score && (
        <section className="rounded-xl border border-border bg-card p-5 mb-8">
          <h2 className="text-sm font-semibold uppercase tracking-widest text-muted-foreground mb-4">
            Score Breakdown
          </h2>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <ScoreStat label="Projects" value={score.breakdown?.projects_score ?? 0} color={tierColor} />
            <ScoreStat label="Reviews" value={score.breakdown?.reviews_score ?? 0} color={tierColor} />
            <ScoreStat label="Vouches" value={score.breakdown?.vouches_score ?? 0} color={tierColor} />
            <ScoreStat label="Activity" value={score.breakdown?.activity_score ?? 0} color={tierColor} />
          </div>
          <div className="mt-4 pt-4 border-t border-border flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Total Vouch Score</span>
            <span className="text-2xl font-bold tabular-nums" style={{ color: tierColor }}>
              {formatScore(score.total_score)}
            </span>
          </div>
        </section>
      )}

      {/* Projects */}
      {projects?.items?.length ? (
        <section>
          <h2 className="text-lg font-semibold mb-4">Projects</h2>
          <div className="grid gap-4 sm:grid-cols-2">
            {projects.items.slice(0, 6).map((p: Project) => (
              <Link
                key={p.id}
                href={`/projects/${p.id}`}
                className="rounded-xl border border-border bg-card p-4 hover:border-primary/40 hover:bg-accent/30 transition-colors"
              >
                <p className="font-semibold truncate">{p.name}</p>
                <p className="text-xs text-muted-foreground mt-1 line-clamp-2">{p.description}</p>
                <div className="flex flex-wrap gap-1 mt-2">
                  {p.tags?.slice(0, 3).map((tag: string) => (
                    <span
                      key={tag}
                      className="rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
                <p className="text-xs text-muted-foreground mt-2">
                  {formatRelativeTime(p.created_at)}
                </p>
              </Link>
            ))}
          </div>
        </section>
      ) : null}
    </main>
  );
}

function ScoreStat({
  label,
  value,
  color,
}: {
  label: string;
  value: number;
  color: string;
}) {
  return (
    <div className="text-center">
      <p className="text-2xl font-bold tabular-nums" style={{ color }}>
        {formatScore(value)}
      </p>
      <p className="text-xs text-muted-foreground mt-0.5">{label}</p>
    </div>
  );
}
