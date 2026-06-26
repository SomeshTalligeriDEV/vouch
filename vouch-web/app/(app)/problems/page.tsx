"use client";

import { useState } from "react";

import { ProblemCard } from "@/components/problem/problem-card";
import { useProblems } from "@/hooks";
import { api } from "@/lib/api";
import { useAuth } from "@/store/auth";
import { useQueryClient } from "@tanstack/react-query";

const INITIAL_FORM = {
  title: "",
  description: "",
  tags: "",
  budget_min: "",
  budget_max: "",
};

function PostProblemForm({ onClose }: { onClose: () => void }) {
  const [form, setForm] = useState(INITIAL_FORM);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const qc = useQueryClient();

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title || !form.description) return;
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
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const field = (label: string, key: keyof typeof form, type = "text", placeholder = "") => (
    <label className="block">
      <span className="mb-1 block text-sm text-ink/60">{label}</span>
      <input
        type={type}
        value={form[key]}
        onChange={(e) => setForm({ ...form, [key]: e.target.value })}
        placeholder={placeholder}
        className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
      />
    </label>
  );

  return (
    <form onSubmit={onSubmit} className="card space-y-4">
      <h2 className="font-semibold">Post a problem</h2>
      <p className="text-xs text-ink/50">
        Describe a real problem you need solved. Builders will claim and ship
        a solution — your budget becomes their first paying customer.
      </p>

      {field("Title *", "title", "text", "e.g. I need a tool that tracks my Notion deadlines in Slack")}

      <label className="block">
        <span className="mb-1 block text-sm text-ink/60">Description *</span>
        <textarea
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          rows={4}
          placeholder="Describe the problem in detail. What have you tried? What outcome do you need?"
          className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
        />
      </label>

      {field("Tags (comma-separated)", "tags", "text", "productivity, slack, notion")}

      <div className="grid grid-cols-2 gap-3">
        {field("Budget min (USD)", "budget_min", "number", "100")}
        {field("Budget max (USD)", "budget_max", "number", "500")}
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={loading || !form.title || !form.description}
          className="btn-primary flex-1"
        >
          {loading ? "Posting…" : "Post problem"}
        </button>
        <button type="button" className="btn-ghost" onClick={onClose}>
          Cancel
        </button>
      </div>

      {error && <p className="text-sm text-red-400">{error}</p>}
    </form>
  );
}

export default function ProblemsPage() {
  const [status, setStatus] = useState("");
  const [showForm, setShowForm] = useState(false);
  const user = useAuth((s) => s.user);
  const { data, isLoading, isError } = useProblems({ status, sort: "upvotes" });

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">Demand board</h1>
          <p className="text-sm text-ink/60">
            Real problems with budgets. Builders claim them, ship a solution,
            and earn the first paying customer.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
          >
            <option value="">All</option>
            <option value="open">Open</option>
            <option value="claimed">Claimed</option>
            <option value="shipped">Shipped</option>
          </select>
          {user ? (
            <button
              className="btn-primary"
              onClick={() => setShowForm((v) => !v)}
            >
              {showForm ? "Cancel" : "+ Post a problem"}
            </button>
          ) : (
            <a href="/login" className="btn-primary">
              Sign in to post
            </a>
          )}
        </div>
      </div>

      {showForm && <PostProblemForm onClose={() => setShowForm(false)} />}

      {isLoading && <p className="text-ink/60">Loading…</p>}
      {isError && <p className="text-red-400">Failed to load problems.</p>}
      {data && data.items.length === 0 && (
        <div className="card py-12 text-center">
          <p className="text-sm text-ink/60">No problems posted yet.</p>
          {user && (
            <button className="btn-primary mt-4" onClick={() => setShowForm(true)}>
              Post the first problem
            </button>
          )}
        </div>
      )}

      <div className="grid gap-4">
        {data?.items.map((p) => (
          <ProblemCard key={p.id} problem={p} />
        ))}
      </div>
    </div>
  );
}
