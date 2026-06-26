"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { api } from "@/lib/api";
import { useAuth } from "@/store/auth";
import type { AdminStats, Company } from "@/types";

export default function AdminDashboard() {
  const user = useAuth((s) => s.user);
  const router = useRouter();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) { router.replace("/login"); return; }
    if (user.role !== "admin") { router.replace("/dashboard"); return; }

    Promise.all([api.adminStats(), api.adminListCompanies()])
      .then(([s, c]) => {
        setStats(s);
        setCompanies(c.items);
      })
      .catch((err) => setError((err as Error).message));
  }, [user, router]);

  if (!user || user.role !== "admin") return null;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold">Admin dashboard</h1>
        <p className="text-sm text-ink/60">Platform-wide overview</p>
      </div>

      {error && <p className="text-red-400 text-sm">{error}</p>}

      {/* Stats grid */}
      {stats && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[
            { label: "Builders", value: stats.total_users },
            { label: "Companies", value: stats.total_companies },
            { label: "Projects", value: stats.total_projects },
            { label: "Open problems", value: stats.open_problems },
          ].map(({ label, value }) => (
            <div key={label} className="card text-center">
              <div className="text-3xl font-bold">{value ?? "—"}</div>
              <div className="text-xs text-ink/60 mt-1 uppercase tracking-wider font-mono">{label}</div>
            </div>
          ))}
        </div>
      )}

      {/* Quick links */}
      <div className="grid gap-3 sm:grid-cols-3">
        <Link href="/discover" className="card hover:border-accent/40">
          <h3 className="font-semibold">All projects</h3>
          <p className="text-sm text-ink/60 mt-1">Browse every shipped project on the platform.</p>
        </Link>
        <Link href="/problems" className="card hover:border-accent/40">
          <h3 className="font-semibold">Demand board</h3>
          <p className="text-sm text-ink/60 mt-1">All problems posted by companies and users.</p>
        </Link>
        <Link href="/score" className="card hover:border-accent/40">
          <h3 className="font-semibold">Leaderboard</h3>
          <p className="text-sm text-ink/60 mt-1">Top builders by verified score.</p>
        </Link>
      </div>

      {/* Companies */}
      <div>
        <h2 className="font-semibold mb-3">Registered companies ({companies.length})</h2>
        {companies.length === 0 && <p className="text-sm text-ink/50">No companies yet.</p>}
        <div className="space-y-2">
          {companies.map((c) => (
            <div key={c.id} className="card flex items-center justify-between">
              <div className="flex items-center gap-3">
                {c.logo_url && (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={c.logo_url} alt={c.name} className="h-8 w-8 rounded-lg border border-line" />
                )}
                <div>
                  <p className="font-medium text-sm">{c.name}</p>
                  <p className="text-xs text-ink/50">{c.email} · {c.size} people</p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                {c.is_verified && (
                  <span className="text-xs font-mono text-emerald-400">Verified</span>
                )}
                {c.website && (
                  <a href={c.website} target="_blank" rel="noopener noreferrer"
                    className="text-xs text-ink/50 hover:text-ink">
                    Website ↗
                  </a>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
