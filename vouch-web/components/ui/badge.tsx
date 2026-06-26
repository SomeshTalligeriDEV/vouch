import { cn } from "@/lib/utils";
import { TIER_COLORS, TIER_LABELS } from "@/lib/constants";

interface BadgeProps {
  children: React.ReactNode;
  variant?: "default" | "secondary" | "outline" | "destructive";
  className?: string;
}

export function Badge({ children, variant = "default", className }: BadgeProps) {
  const variants = {
    default: "bg-primary text-primary-foreground",
    secondary: "bg-secondary text-secondary-foreground",
    outline: "border border-border text-foreground",
    destructive: "bg-destructive text-destructive-foreground",
  };

  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold",
        variants[variant],
        className,
      )}
    >
      {children}
    </span>
  );
}

interface TierBadgeProps {
  tier: number;
  className?: string;
}

export function TierBadge({ tier, className }: TierBadgeProps) {
  const label = TIER_LABELS[tier] ?? `T${tier}`;
  const color = TIER_COLORS[tier] ?? "bg-gray-100 text-gray-700";
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold",
        color,
        className,
      )}
    >
      {label}
    </span>
  );
}
