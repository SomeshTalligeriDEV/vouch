import type { Tier } from "@/types";
import { TIER_COLORS, cn, formatScore } from "@/lib/utils";

export function ScoreBadge({
  tier,
  score,
}: {
  tier: Tier;
  score: number;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-2 rounded-full border px-3 py-1 text-sm font-semibold",
        TIER_COLORS[tier],
      )}
    >
      <span>{tier}</span>
      <span className="opacity-70">·</span>
      <span>{formatScore(score)}</span>
    </span>
  );
}
