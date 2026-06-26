import { cn } from "@/lib/utils";

interface StatCardProps {
  label: string;
  value: string | number;
  sub?: string;
  icon?: React.ReactNode;
  trend?: "up" | "down" | "neutral";
  className?: string;
}

export function StatCard({ label, value, sub, icon, trend, className }: StatCardProps) {
  return (
    <div className={cn("rounded-xl border bg-card p-6 space-y-2", className)}>
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground font-medium">{label}</p>
        {icon && <span className="text-muted-foreground opacity-70">{icon}</span>}
      </div>
      <p className="text-3xl font-bold tracking-tight tabular-nums">{value}</p>
      {(sub || trend) && (
        <p
          className={cn(
            "text-xs",
            trend === "up" && "text-green-600",
            trend === "down" && "text-red-500",
            !trend && "text-muted-foreground",
          )}
        >
          {trend === "up" && "↑ "}
          {trend === "down" && "↓ "}
          {sub}
        </p>
      )}
    </div>
  );
}
