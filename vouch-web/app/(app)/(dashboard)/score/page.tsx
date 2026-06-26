"use client";

import { useState } from "react";

import { ScoreBreakdown } from "@/components/score/score-breakdown";
import { api } from "@/lib/api";
import { useScore, useScoreSSE } from "@/hooks";

export default function ScorePage() {
  const [username, setUsername] = useState("");
  const [submitted, setSubmitted] = useState("");
  const { data: score, isLoading, isError } = useScore(submitted);
  useScoreSSE(submitted); // auto-refresh on real-time score updates
  const [recalcMsg, setRecalcMsg] = useState<string | null>(null);

  const recalc = async () => {
    try {
      await api.recalculateScore();
      setRecalcMsg("Recalculation queued — refresh in a moment.");
    } catch {
      setRecalcMsg("You must be signed in to recalculate your score.");
    }
  };

  return (
    <div className="mx-auto max-w-xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Builder score</h1>
        <p className="text-sm text-ink/60">
          Look up any builder&apos;s verified score.
        </p>
      </div>

      <form
        onSubmit={(e) => {
          e.preventDefault();
          setSubmitted(username.trim());
        }}
        className="flex gap-2"
      >
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="builder username"
          className="flex-1 rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-crimson"
        />
        <button type="submit" className="btn-primary">
          Look up
        </button>
      </form>

      {isLoading && submitted && <p className="text-ink/60">Loading…</p>}
      {isError && <p className="text-red-400">No score found for that user.</p>}
      {score && <ScoreBreakdown score={score} />}

      <div className="card flex items-center justify-between">
        <span className="text-sm text-ink/60">Recalculate your own score</span>
        <button onClick={recalc} className="btn-ghost">
          Recalculate
        </button>
      </div>
      {recalcMsg && <p className="text-sm text-ink/60">{recalcMsg}</p>}
    </div>
  );
}
