"use client";

import { useState } from "react";

import { ImageUpload } from "@/components/ui/image-upload";
import { useCreateProject, useProjects } from "@/hooks";
import { useAuth } from "@/store/auth";
import { ProjectCard } from "@/components/project/project-card";
import type { ProjectStatus } from "@/types";

const STATUSES: { value: ProjectStatus; label: string }[] = [
  { value: "draft", label: "Draft" },
  { value: "live", label: "Live" },
  { value: "acquired", label: "Acquired" },
];

const INITIAL = {
  title: "",
  tagline: "",
  description: "",
  logo_url: "",
  live_url: "",
  repo_url: "",
  payment_link: "",
  tags: "",
  status: "live" as ProjectStatus,
  for_sale: false,
  ask_price: "",
};

export default function ProjectsPage() {
  const user = useAuth((s) => s.user);
  const create = useCreateProject();
  const [form, setForm] = useState(INITIAL);
  const [open, setOpen] = useState(false);

  const { data: myProjects, isLoading } = useProjects(
    user ? { status: undefined } : {},
  );

  const reset = () => { setForm(INITIAL); setOpen(false); };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    create.mutate(
      {
        ...form,
        tags: form.tags.split(",").map((t) => t.trim()).filter(Boolean),
        ask_price: form.for_sale ? parseFloat(form.ask_price) || 0 : 0,
      },
      { onSuccess: reset },
    );
  };

  const field = (
    label: string,
    key: keyof typeof form,
    type = "text",
    placeholder = "",
  ) => (
    <label className="block">
      <span className="mb-1 block text-sm text-ink/60">{label}</span>
      <input
        type={type}
        value={form[key] as string}
        onChange={(e) => setForm({ ...form, [key]: e.target.value })}
        placeholder={placeholder}
        className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
      />
    </label>
  );

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Your projects</h1>
          <p className="text-sm text-ink/60">
            Every shipped product feeds your verified Builder Score.
          </p>
        </div>
        <button className="btn-primary" onClick={() => setOpen((v) => !v)}>
          {open ? "Cancel" : "+ Add project"}
        </button>
      </div>

      {open && (
        <form onSubmit={onSubmit} className="card space-y-4">
          <ImageUpload
            value={form.logo_url}
            onChange={(url) => setForm({ ...form, logo_url: url })}
          />

          {field("Title *", "title", "text", "e.g. ShipFast")}
          {field("Tagline *", "tagline", "text", "One line pitch")}

          <label className="block">
            <span className="mb-1 block text-sm text-ink/60">Description</span>
            <textarea
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              rows={3}
              placeholder="What does it do, who is it for, what problem does it solve?"
              className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
            />
          </label>

          {field("Live URL", "live_url", "url", "https://")}
          {field("Repo URL", "repo_url", "url", "https://github.com/...")}
          {field("Payment link", "payment_link", "url", "https://buy.stripe.com/...")}
          {field("Tags (comma-separated)", "tags", "text", "saas, ai, productivity")}

          <label className="block">
            <span className="mb-1 block text-sm text-ink/60">Status</span>
            <select
              value={form.status}
              onChange={(e) => setForm({ ...form, status: e.target.value as ProjectStatus })}
              className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
            >
              {STATUSES.map((s) => (
                <option key={s.value} value={s.value}>{s.label}</option>
              ))}
            </select>
          </label>

          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={form.for_sale}
              onChange={(e) => setForm({ ...form, for_sale: e.target.checked })}
            />
            List this project for sale
          </label>

          {form.for_sale && (
            <label className="block">
              <span className="mb-1 block text-sm text-ink/60">Ask price (USD)</span>
              <input
                type="number"
                min={0}
                step={100}
                value={form.ask_price}
                onChange={(e) => setForm({ ...form, ask_price: e.target.value })}
                placeholder="5000"
                className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
              />
            </label>
          )}

          <button
            type="submit"
            disabled={create.isPending || !form.title || !form.tagline}
            className="btn-primary w-full"
          >
            {create.isPending ? "Creating…" : "Create project"}
          </button>

          {create.isError && (
            <p className="text-sm text-red-400">{(create.error as Error).message}</p>
          )}
          {create.isSuccess && (
            <p className="text-sm text-emerald-400">Project created — score recalculation queued.</p>
          )}
        </form>
      )}

      {isLoading && <p className="text-ink/60">Loading your projects…</p>}

      {!isLoading && myProjects?.items.length === 0 && !open && (
        <div className="card text-center py-12">
          <p className="text-ink/60 text-sm">No projects yet.</p>
          <button className="btn-primary mt-4" onClick={() => setOpen(true)}>
            Add your first project
          </button>
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {myProjects?.items.map((p) => (
          <ProjectCard key={p.id} project={p} />
        ))}
      </div>
    </div>
  );
}
