"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { api } from "@/lib/api";
import { useProblems } from "@/hooks";
import { useQueryClient } from "@tanstack/react-query";
import type { Company } from "@/types";

const STATUS_LABEL: Record<string, string> = {
  open: "Open",
  claimed: "Claimed",
  shipped: "Shipped ✓",
  cancelled: "Cancelled",
};
const STATUS_COLOR: Record<string, string> = {
  open: "text-emerald-400",
  claimed: "text-yellow-400",
  shipped: "text-accent-ink",
  cancelled: "text-ink/40",
};

function PostForm({ onClose }: { onClose: () => void }) {
  const qc = useQueryClient();
  const [form, setForm] = useState({ title: "", description: "", tags: "", budget_min: "", budget_max: "" });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await api.createProblem({
        title: form.title,
        description: form.description,
        tags: form.tags.split(",").map((t) => t.trim()).filter(Boolean),
        budget_min: parseFloat(form.budget_min) || 0,
        budget_max: parseFloat(form.budget_max) || 0,
      });
      qc.invalidateQueries({ queryKey: ["problems"] });
      onClose();
    } catch (err) { setError((err as Error).message); }
    finally { setLoading(false); }
  };

  const f = (label: string, key: keyof typeof form, type = "text", ph = "") => (
    <label className="block">
      <span className="mb-1 block text-xs text-ink/60">{label}</span>
      <input type={type} value={form[key]} onChange={(e) => setForm({ ...form, [key]: e.target.value })}
        placeholder={ph}
        className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent" />
    </label>
  );

  return (
    <form onSubmit={submit} className="card space-y-3">
      <h3 className="font-semibold">Post a new problem</h3>
      {f("Title *", "title", "text", "e.g. Need a Slack bot that tracks my Notion tasks")}
      <label className="block">
        <span className="mb-1 block text-xs text-ink/60">Description *</span>
        <textarea value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })}
          rows={3} placeholder="Describe the problem in detail..."
          className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent" />
      </label>
      {f("Tags (comma-separated)", "tags", "text", "saas, productivity")}
      <div className="grid grid-cols-2 gap-2">
        {f("Budget min (USD)", "budget_min", "number", "100")}
        {f("Budget max (USD)", "budget_max", "number", "500")}
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={loading || !form.title || !form.description} className="btn-primary flex-1">
          {loading ? "Posting…" : "Post problem"}
        </button>
        <button type="button" className="btn-ghost" onClick={onClose}>Cancel</button>
      </div>
      {error && <p className="text-sm text-red-400">{error}</p>}
    </form>
  );
}

export default function CompanyDashboard() {
  const router = useRouter();
  const [company, setCompany] = useState<Company | null>(null);
  const [showForm, setShowForm] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem("vouch_company");
    if (!stored) { router.replace("/company/login"); return; }
    try { setCompany(JSON.parse(stored)); } catch { router.replace("/company/login"); }
  }, [router]);

  // Problems posted by this company (filter by poster_id not supported yet — show all for now)
  const { data: problems, isLoading } = useProblems({ sort: "recent" });

  const logout = () => {
    localStorage.removeItem("vouch_company");
    router.replace("/company/login");
  };

  if (!company) return null;

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-4">
          {company.logo_url && (
            // eslint-disable-next-line @next/next/no-img-element
            <img src={company.logo_url} alt={company.name} className="h-12 w-12 rounded-xl border border-line" />
          )}
          <div>
            <h1 className="text-2xl font-bold">{company.name}</h1>
            <p className="text-sm text-ink/60">{company.email}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Link href="/company/settings" className="btn-ghost text-sm">Settings</Link>
          <button onClick={logout} className="btn-ghost text-sm text-ink/50">Sign out</button>
        </div>
      </div>

      {/* Quick stats */}
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="card text-center">
          <div className="text-2xl font-bold">{problems?.items.length ?? "—"}</div>
          <div className="text-xs text-ink/60 mt-1">Problems posted</div>
        </div>
        <div className="card text-center">
          <div className="text-2xl font-bold">
            {problems?.items.filter((p) => p.status === "claimed").length ?? "—"}
          </div>
          <div className="text-xs text-ink/60 mt-1">Claimed by builders</div>
        </div>
        <div className="card text-center">
          <div className="text-2xl font-bold">
            {problems?.items.filter((p) => p.status === "shipped").length ?? "—"}
          </div>
          <div className="text-xs text-ink/60 mt-1">Solutions shipped</div>
        </div>
      </div>

      {/* Post new problem */}
      <div className="flex items-center justify-between">
        <h2 className="font-semibold">Your problems</h2>
        <button className="btn-primary" onClick={() => setShowForm((v) => !v)}>
          {showForm ? "Cancel" : "+ Post a problem"}
        </button>
      </div>

      {showForm && <PostForm onClose={() => setShowForm(false)} />}

      {/* Problem list */}
      {isLoading && <p className="text-ink/60 text-sm">Loading…</p>}
      {problems?.items.length === 0 && !showForm && (
        <div className="card py-12 text-center">
          <p className="text-sm text-ink/60">No problems posted yet.</p>
          <p className="text-xs text-ink/40 mt-1">
            Post a problem to let builders compete to solve it.
          </p>
          <button className="btn-primary mt-4" onClick={() => setShowForm(true)}>
            Post your first problem
          </button>
        </div>
      )}
      <div className="space-y-3">
        {problems?.items.map((p) => (
          <div key={p.id} className="card">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <h3 className="font-semibold truncate">{p.title}</h3>
                <p className="mt-1 text-sm text-ink/60 line-clamp-2">{p.description}</p>
              </div>
              <span className={`shrink-0 text-xs font-mono uppercase ${STATUS_COLOR[p.status] ?? ""}`}>
                {STATUS_LABEL[p.status] ?? p.status}
              </span>
            </div>
            <div className="mt-3 flex gap-4 text-xs text-ink/50">
              <span>▲ {p.upvotes} upvotes</span>
              {p.budget_max > 0 && <span>Budget: ${p.budget_min}–${p.budget_max}</span>}
              {p.tags?.map((t) => <span key={t} className="tag">{t}</span>)}
            </div>
            {p.status === "claimed" && (
              <p className="mt-2 text-xs text-yellow-400">
                A builder has claimed this — solution incoming.
              </p>
            )}
            {p.status === "shipped" && p.shipped_project_id && (
              <Link href={`/projects/${p.shipped_project_id}`}
                className="mt-2 inline-block text-xs text-accent-ink underline">
                View shipped solution →
              </Link>
            )}
          </div>
        ))}
      </div>

      {/* Discover builders */}
      <div className="card">
        <h2 className="font-semibold mb-2">Top builders</h2>
        <p className="text-sm text-ink/60 mb-3">
          Browse verified builders by score — find who to work with.
        </p>
        <Link href="/discover" className="btn-primary inline-block">
          Browse builders & projects
        </Link>
      </div>
    </div>
  );
}
