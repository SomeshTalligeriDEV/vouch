import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { formatRelativeTime } from "@/lib/format";

interface ReviewCardProps {
  reviewerName: string;
  reviewerAvatar?: string | null;
  reviewerUsername: string;
  rating: number;
  body: string;
  verifiedPurchase?: boolean;
  createdAt: string;
}

function Stars({ rating }: { rating: number }) {
  return (
    <span aria-label={`${rating} out of 5 stars`} className="text-yellow-500 text-sm">
      {"★".repeat(Math.min(5, Math.max(0, rating)))}
      {"☆".repeat(5 - Math.min(5, Math.max(0, rating)))}
    </span>
  );
}

export function ReviewCard({
  reviewerName,
  reviewerAvatar,
  reviewerUsername,
  rating,
  body,
  verifiedPurchase,
  createdAt,
}: ReviewCardProps) {
  return (
    <div className="rounded-xl border bg-card p-5 space-y-3">
      <div className="flex items-center gap-3">
        <Avatar src={reviewerAvatar} alt={reviewerName} size={36} />
        <div>
          <p className="text-sm font-medium">{reviewerName}</p>
          <p className="text-xs text-muted-foreground">@{reviewerUsername}</p>
        </div>
        <div className="ml-auto flex items-center gap-2">
          <Stars rating={rating} />
          {verifiedPurchase && (
            <Badge variant="secondary">Verified</Badge>
          )}
        </div>
      </div>
      <p className="text-sm text-muted-foreground leading-relaxed">{body}</p>
      <p className="text-xs text-muted-foreground">{formatRelativeTime(createdAt)}</p>
    </div>
  );
}
