"use client";

import { useState } from "react";

import { ProjectCard } from "@/components/project/project-card";
import { useProjects } from "@/hooks";

export default function DiscoverPage() {
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState("recent");
  const { data, isLoading, isError } = useProjects({ search, sort });

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">Discover</h1>
          <p className="text-sm text-ink/60">
            Real products with real users and revenue.
          </p>
        </div>
        <div className="flex gap-2">
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search projects…"
            className="rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-crimson"
          />
          <select
            value={sort}
            onChange={(e) => setSort(e.target.value)}
            className="rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-crimson"
          >
            <option value="recent">Recent</option>
            <option value="mrr">Top MRR</option>
            <option value="users">Most users</option>
            <option value="rating">Best rated</option>
          </select>
        </div>
      </div>

      {isLoading && <p className="text-ink/60">Loading…</p>}
      {isError && <p className="text-red-400">Failed to load projects.</p>}
      {data && data.items.length === 0 && (
        <p className="text-ink/60">No projects yet — be the first to ship.</p>
      )}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {data?.items.map((p) => (
          <ProjectCard key={p.id} project={p} />
        ))}
      </div>
    </div>
  );
}
