import { SkeletonCard } from "@/components/ui/skeleton";

export default function ProblemsLoading() {
  return (
    <div className="mx-auto max-w-5xl px-4 py-10">
      <div className="h-8 w-40 rounded-md bg-muted animate-pulse mb-6" />
      <div className="grid gap-4 sm:grid-cols-2">
        {Array.from({ length: 6 }).map((_, i) => (
          <SkeletonCard key={i} />
        ))}
      </div>
    </div>
  );
}
