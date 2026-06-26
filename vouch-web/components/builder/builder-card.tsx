import Link from "next/link";
import { Avatar } from "@/components/ui/avatar";
import { TierBadge } from "@/components/ui/badge";
import { ScoreRing } from "@/components/ui/score-ring";
import { formatScore } from "@/lib/format";

interface BuilderCardProps {
  username: string;
  name?: string;
  avatarUrl?: string | null;
  bio?: string;
  totalScore: number;
  tier: number;
  rank?: number;
}

export function BuilderCard({ username, name, avatarUrl, bio, totalScore, tier, rank }: BuilderCardProps) {
  return (
    <Link
      href={`/u/${username}`}
      className="group block rounded-xl border bg-card p-5 hover:shadow-md transition-shadow"
    >
      <div className="flex items-start gap-4">
        {rank != null && (
          <span className="shrink-0 w-8 text-center text-sm font-bold text-muted-foreground tabular-nums">
            #{rank}
          </span>
        )}
        <Avatar src={avatarUrl} alt={name ?? username} size={48} />
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap">
            <p className="font-semibold text-foreground group-hover:text-primary transition-colors truncate">
              {name ?? username}
            </p>
            <TierBadge tier={tier} />
          </div>
          <p className="text-xs text-muted-foreground mt-0.5">@{username}</p>
          {bio && <p className="text-sm text-muted-foreground mt-1 line-clamp-2">{bio}</p>}
        </div>
        <ScoreRing score={totalScore} size={56} strokeWidth={5} />
      </div>
      <div className="mt-3 flex items-center justify-between text-xs text-muted-foreground">
        <span>Score: <strong className="text-foreground">{formatScore(totalScore)}</strong></span>
      </div>
    </Link>
  );
}
