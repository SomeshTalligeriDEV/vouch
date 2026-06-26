import { cn } from "@/lib/utils";

interface ScoreRingProps {
  score: number;
  maxScore?: number;
  size?: number;
  strokeWidth?: number;
  className?: string;
}

export function ScoreRing({
  score,
  maxScore = 1000,
  size = 80,
  strokeWidth = 6,
  className,
}: ScoreRingProps) {
  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const pct = Math.min(score / maxScore, 1);
  const offset = circumference * (1 - pct);

  const color =
    score >= 900
      ? "#f59e0b" // elite — amber
      : score >= 700
        ? "#8b5cf6" // pro — violet
        : score >= 500
          ? "#3b82f6" // builder — blue
          : score >= 300
            ? "#22c55e" // rising — green
            : "#9ca3af"; // newcomer — gray

  return (
    <div
      className={cn("relative inline-flex items-center justify-center", className)}
      style={{ width: size, height: size }}
      aria-label={`Score: ${score}`}
    >
      <svg width={size} height={size} className="-rotate-90">
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          className="text-muted opacity-20"
        />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke={color}
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          style={{ transition: "stroke-dashoffset 0.6s ease" }}
        />
      </svg>
      <span className="absolute text-sm font-bold tabular-nums">{score}</span>
    </div>
  );
}
