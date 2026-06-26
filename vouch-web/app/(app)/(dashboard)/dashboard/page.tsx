"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useRef, useState } from "react";

import { StripeConnect } from "@/components/builder/stripe-connect";
import { ScoreBadge } from "@/components/score/score-badge";
import { api } from "@/lib/api";
import { useLeaderboard, useLeaderboardSSE } from "@/hooks";

function StripeReturn() {
  const params = useSearchParams();
  const router = useRouter();
  const handled = useRef(false);
  const [state, setState] = useState<"idle" | "connecting" | "done" | "error">(
    "idle",
  );

  useEffect(() => {
    const code = params.get("code");
    if (!code || handled.current) return;
    handled.current = true;
    setState("connecting");
    api
      .connectStripe(code)
      .then(() => {
        setState("done");
        router.replace("/dashboard");
      })
      .catch(() => setState("error"));
  }, [params, router]);

  return (
    <div className="card flex items-center justify-between">
      <div>
        <h3 className="font-semibold">Verify your revenue</h3>
        <p className="text-sm text-ink/60">
          Connect Stripe (read-only) to lift your score multiplier 0.6 → 1.0.
        </p>
        {state === "connecting" && (
          <p className="mt-2 text-sm text-ink/60">Connecting Stripe…</p>
        )}
        {state === "done" && (
          <p className="mt-2 text-sm text-emerald-400">Stripe connected.</p>
        )}
        {state === "error" && (
          <p className="mt-2 text-sm text-red-400">
            Stripe connection failed — please retry.
          </p>
        )}
      </div>
      <StripeConnect connected={state === "done"} />
    </div>
  );
}

export default function DashboardPage() {
  const { data: leaders, isLoading } = useLeaderboard(10);
  useLeaderboardSSE(); // subscribe to real-time updates

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-sm text-ink/60">
          Track your score and the top builders shipping right now.
        </p>
      </div>

      <Suspense fallback={null}>
        <StripeReturn />
      </Suspense>

      <div className="grid gap-4 sm:grid-cols-3">
        <Link href="/projects" className="card hover:border-crimson/40">
          <h3 className="font-semibold">Your projects</h3>
          <p className="mt-1 text-sm text-ink/60">
            Add a shipped product to grow your score.
          </p>
        </Link>
        <Link href="/score" className="card hover:border-crimson/40">
          <h3 className="font-semibold">Your score</h3>
          <p className="mt-1 text-sm text-ink/60">
            See your breakdown and recalculate.
          </p>
        </Link>
        <Link href="/problems" className="card hover:border-crimson/40">
          <h3 className="font-semibold">Claim demand</h3>
          <p className="mt-1 text-sm text-ink/60">
            Ship a solution, win the first paying user.
          </p>
        </Link>
      </div>

      <div className="card">
        <h2 className="mb-4 font-semibold">Leaderboard</h2>
        {isLoading && <p className="text-ink/60">Loading…</p>}
        <ol className="space-y-2">
          {leaders?.map((s, i) => (
            <li
              key={s.id}
              className="flex items-center justify-between rounded-lg border border-line px-3 py-2"
            >
              <Link href={`/builder/${s.username || s.builder_id}`} className="flex items-center gap-3 hover:opacity-80">
                <span className="w-6 text-center text-xs text-ink/50">#{i + 1}</span>
                {s.avatar_url && (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={s.avatar_url} alt={s.username} className="h-7 w-7 rounded-full border border-line" />
                )}
                <div>
                  <span className="block text-sm font-medium text-ink">{s.name || s.username || "—"}</span>
                  {s.username && <span className="block text-xs text-ink/50">@{s.username}</span>}
                </div>
              </Link>
              <ScoreBadge tier={s.tier} score={s.total_score} />
            </li>
          ))}
        </ol>
      </div>
    </div>
  );
}
