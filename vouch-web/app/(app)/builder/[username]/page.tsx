"use client";

import { use } from "react";

import { ScoreBadge } from "@/components/score/score-badge";
import { ScoreBreakdown } from "@/components/score/score-breakdown";
import { useBuilder, useScore, useScoreSSE } from "@/hooks";

export default function BuilderProfilePage({
  params,
}: {
  params: Promise<{ username: string }>;
}) {
  const { username } = use(params);
  const { data: builder, isLoading, isError } = useBuilder(username);
  const { data: score } = useScore(username);
  useScoreSSE(username);

  if (isLoading) return <p className="text-ink/60">Loading builder…</p>;
  if (isError || !builder)
    return <p className="text-red-400">Builder not found.</p>;

  return (
    <div className="grid gap-6 lg:grid-cols-3">
      <div className="lg:col-span-2 space-y-6">
        <div className="card">
          <div className="flex items-center gap-4">
            {builder.avatar_url && (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={builder.avatar_url}
                alt={builder.username}
                className="h-16 w-16 rounded-full border border-line"
              />
            )}
            <div>
              <h1 className="text-2xl font-bold">
                {builder.name || builder.username}
              </h1>
              <p className="text-sm text-ink/60">@{builder.username}</p>
            </div>
            {score && (
              <div className="ml-auto">
                <ScoreBadge tier={score.tier} score={score.total_score} />
              </div>
            )}
          </div>
          {builder.bio && (
            <p className="mt-4 text-ink/80">{builder.bio}</p>
          )}
          <div className="mt-4 flex gap-4 text-sm text-ink/60">
            {builder.website_url && (
              <a href={builder.website_url} className="hover:text-crimson">
                Website
              </a>
            )}
            {builder.github_login && (
              <a
                href={`https://github.com/${builder.github_login}`}
                className="hover:text-crimson"
              >
                GitHub
              </a>
            )}
          </div>
        </div>
      </div>

      <div>{score && <ScoreBreakdown score={score} />}</div>
    </div>
  );
}
