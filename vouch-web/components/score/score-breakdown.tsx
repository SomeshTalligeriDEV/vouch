import type { BuilderScore } from "@/types";
import { formatScore } from "@/lib/utils";

const ROWS: { key: keyof BuilderScore["breakdown"]; label: string; cap: number }[] = [
  { key: "user", label: "Users", cap: 30000 },
  { key: "revenue", label: "Revenue", cap: 20000 },
  { key: "impact", label: "Impact", cap: 15000 },
  { key: "velocity", label: "Velocity", cap: 5000 },
];

export function ScoreBreakdown({ score }: { score: BuilderScore }) {
  return (
    <div className="card space-y-4">
      <div className="flex items-baseline justify-between">
        <h3 className="font-semibold">Score breakdown</h3>
        <span className="text-sm text-ink/60">
          ×{score.stripe_multiplier.toFixed(1)}{" "}
          {score.stripe_verified ? "Stripe verified" : "unverified"}
        </span>
      </div>
      <div className="space-y-3">
        {ROWS.map((row) => {
          const value = score.breakdown[row.key];
          const pct = Math.min((value / row.cap) * 100, 100);
          return (
            <div key={row.key}>
              <div className="mb-1 flex justify-between text-sm">
                <span className="text-ink/80">{row.label}</span>
                <span className="text-ink/60">{formatScore(value)}</span>
              </div>
              <div className="h-2 overflow-hidden rounded-full bg-ink/10">
                <div
                  className="h-full rounded-full bg-crimson"
                  style={{ width: `${pct}%` }}
                />
              </div>
            </div>
          );
        })}
      </div>
      <div className="flex items-baseline justify-between border-t border-line pt-3">
        <span className="font-semibold">Total</span>
        <span className="text-xl font-bold text-crimson">
          {formatScore(score.total_score)}
        </span>
      </div>
    </div>
  );
}
