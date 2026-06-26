import { Skeleton } from "@/components/ui/skeleton";

export default function ProblemLoading() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-10 space-y-6">
      <Skeleton className="h-8 w-72" />
      <Skeleton className="h-32 w-full" />
      <div className="flex gap-2">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={i} className="h-6 w-16 rounded-full" />
        ))}
      </div>
      <Skeleton className="h-10 w-36 rounded-md" />
    </div>
  );
}
