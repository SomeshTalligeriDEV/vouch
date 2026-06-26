"use client";

import { use, useState } from "react";
import Link from "next/link";

import { useProject, useReviews } from "@/hooks";
import { useAuth } from "@/store/auth";
import { api } from "@/lib/api";
import { useQueryClient } from "@tanstack/react-query";
import { formatCurrency } from "@/lib/utils";

function ReviewForm({ projectId }: { projectId: string }) {
  const qc = useQueryClient();
  const [rating, setRating] = useState(5);
  const [body, setBody] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState(false);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await api.createReview({ project_id: projectId, rating, body });
      qc.invalidateQueries({ queryKey: ["reviews", projectId] });
      qc.invalidateQueries({ queryKey: ["project", projectId] });
      setDone(true);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  if (done) return <p className="text-sm text-emerald-400">Review submitted — thanks!</p>;

  return (
    <form onSubmit={submit} className="card space-y-3">
      <h3 className="font-semibold text-sm">Leave a review</h3>
      <label className="block">
        <span className="mb-1 block text-xs text-ink/60">Rating (1–5)</span>
        <input
          type="number"
          min={1}
          max={5}
          value={rating}
          onChange={(e) => setRating(parseInt(e.target.value))}
          className="w-20 rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
        />
      </label>
      <label className="block">
        <span className="mb-1 block text-xs text-ink/60">Your review</span>
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          rows={3}
          placeholder="How has this project helped you?"
          className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
        />
      </label>
      <button
        type="submit"
        disabled={loading || !body}
        className="btn-primary"
      >
        {loading ? "Submitting…" : "Submit review"}
      </button>
      {error && <p className="text-sm text-red-400">{error}</p>}
    </form>
  );
}

export default function ProjectDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const { data: project, isLoading, isError } = useProject(id);
  const { data: reviews } = useReviews(id);
  const user = useAuth((s) => s.user);

  if (isLoading) return <p className="text-ink/60">Loading…</p>;
  if (isError || !project) return <p className="text-red-400">Project not found.</p>;

  const isOwner = user?.id === project.builder_id;
  const statusColor: Record<string, string> = {
    live: "text-emerald-400",
    draft: "text-yellow-400",
    acquired: "text-purple-400",
    archived: "text-ink/40",
  };

  return (
    <div className="grid gap-6 lg:grid-cols-3">
      <div className="lg:col-span-2 space-y-6">
        {/* Header */}
        <div className="card">
          <div className="flex items-start gap-4">
            {project.logo_url && (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={project.logo_url}
                alt={project.title}
                className="h-14 w-14 rounded-xl border border-line object-cover"
              />
            )}
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 flex-wrap">
                <h1 className="text-2xl font-bold truncate">{project.title}</h1>
                <span className={`text-xs font-mono uppercase ${statusColor[project.status] ?? ""}`}>
                  {project.status}
                </span>
              </div>
              <p className="mt-1 text-sm text-ink/70">{project.tagline}</p>
              <div className="mt-3 flex flex-wrap gap-2">
                {project.live_url && (
                  <a href={project.live_url} target="_blank" rel="noopener noreferrer"
                    className="text-xs border border-line rounded-full px-3 py-1 hover:border-accent">
                    Live site ↗
                  </a>
                )}
                {project.repo_url && (
                  <a href={project.repo_url} target="_blank" rel="noopener noreferrer"
                    className="text-xs border border-line rounded-full px-3 py-1 hover:border-accent">
                    Repo ↗
                  </a>
                )}
                {project.payment_link && (
                  <a href={project.payment_link} target="_blank" rel="noopener noreferrer"
                    className="btn-primary text-xs">
                    Buy / Subscribe
                  </a>
                )}
              </div>
            </div>
          </div>

          {project.description && (
            <p className="mt-4 text-sm text-ink/80 leading-relaxed whitespace-pre-line">
              {project.description}
            </p>
          )}

          {project.tags?.length > 0 && (
            <div className="mt-4 flex flex-wrap gap-2">
              {project.tags.map((t) => (
                <span key={t} className="tag">{t}</span>
              ))}
            </div>
          )}
        </div>

        {/* Reviews */}
        <div className="space-y-3">
          <h2 className="font-semibold">Reviews ({project.review_count})</h2>
          {user && !isOwner && <ReviewForm projectId={id} />}
          {reviews?.items.map((r) => (
            <div key={r.id} className="card">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">@{r.reviewer_username}</span>
                <span className="text-xs text-ink/50">{"★".repeat(r.rating)}{"☆".repeat(5 - r.rating)}</span>
              </div>
              <p className="mt-2 text-sm text-ink/80">{r.body}</p>
            </div>
          ))}
          {reviews?.items.length === 0 && (
            <p className="text-sm text-ink/50">No reviews yet — be the first.</p>
          )}
        </div>
      </div>

      {/* Sidebar */}
      <div className="space-y-4">
        <div className="card space-y-3">
          <h3 className="font-semibold text-sm">Stats</h3>
          <div className="grid grid-cols-2 gap-3 text-center">
            <div className="rounded-lg bg-ink/5 p-3">
              <div className="text-xl font-bold">{project.verified_users.toLocaleString()}</div>
              <div className="text-xs text-ink/50 mt-0.5">Verified users</div>
            </div>
            <div className="rounded-lg bg-ink/5 p-3">
              <div className="text-xl font-bold">{formatCurrency(project.mrr)}</div>
              <div className="text-xs text-ink/50 mt-0.5">MRR</div>
            </div>
            <div className="rounded-lg bg-ink/5 p-3">
              <div className="text-xl font-bold">{project.review_count}</div>
              <div className="text-xs text-ink/50 mt-0.5">Reviews</div>
            </div>
            <div className="rounded-lg bg-ink/5 p-3">
              <div className="text-xl font-bold">
                {project.average_rating > 0 ? project.average_rating.toFixed(1) : "—"}
              </div>
              <div className="text-xs text-ink/50 mt-0.5">Avg rating</div>
            </div>
          </div>
        </div>

        {project.for_sale && project.ask_price > 0 && (
          <div className="card border-accent/40">
            <p className="text-xs font-mono text-accent uppercase tracking-wider">For sale</p>
            <p className="mt-1 text-2xl font-bold">{formatCurrency(project.ask_price)}</p>
            {project.payment_link && (
              <a href={project.payment_link} target="_blank" rel="noopener noreferrer"
                className="btn-primary mt-3 block text-center">
                Make an offer
              </a>
            )}
          </div>
        )}

        {isOwner && (
          <div className="card">
            <p className="text-xs text-ink/50 mb-2">You own this project</p>
            <Link href="/projects" className="btn-ghost w-full block text-center text-sm">
              Manage projects
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}
