import { Skeleton, SkeletonProfile } from "@/components/ui/skeleton";

export default function LeaderboardLoading() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-10">
      <div className="h-8 w-48 rounded-md bg-muted animate-pulse mb-8" />
      <div className="space-y-3">
        {Array.from({ length: 20 }).map((_, i) => (
          <div key={i} className="flex items-center gap-4 rounded-xl border px-4 py-3">
            <Skeleton className="h-5 w-6 shrink-0" />
            <SkeletonProfile />
            <Skeleton className="ml-auto h-5 w-16" />
          </div>
        ))}
      </div>
    </div>
  );
}
