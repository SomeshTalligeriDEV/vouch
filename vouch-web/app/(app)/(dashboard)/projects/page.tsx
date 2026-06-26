"use client";

import { useState } from "react";

import { ImageUpload } from "@/components/ui/image-upload";
import { useCreateProject } from "@/hooks";

export default function ProjectsPage() {
  const create = useCreateProject();
  const [form, setForm] = useState({
    title: "",
    tagline: "",
    logo_url: "",
    live_url: "",
    repo_url: "",
    payment_link: "",
    for_sale: false,
  });

  const reset = () =>
    setForm({
      title: "",
      tagline: "",
      logo_url: "",
      live_url: "",
      repo_url: "",
      payment_link: "",
      for_sale: false,
    });

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    create.mutate({ ...form, tags: [] }, { onSuccess: reset });
  };

  const field = (label: string, key: keyof typeof form, type = "text") => (
    <label className="block">
      <span className="mb-1 block text-sm text-ink/60">{label}</span>
      <input
        type={type}
        value={form[key] as string}
        onChange={(e) => setForm({ ...form, [key]: e.target.value })}
        className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-crimson"
      />
    </label>
  );

  return (
    <div className="mx-auto max-w-xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Add a project</h1>
        <p className="text-sm text-ink/60">
          Every shipped product feeds your verified Builder Score.
        </p>
      </div>

      <form onSubmit={onSubmit} className="card space-y-4">
        <ImageUpload
          value={form.logo_url}
          onChange={(url) => setForm({ ...form, logo_url: url })}
        />
        {field("Title", "title")}
        {field("Tagline", "tagline")}
        {field("Live URL", "live_url", "url")}
        {field("Repo URL", "repo_url", "url")}
        {field("Payment link", "payment_link", "url")}
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={form.for_sale}
            onChange={(e) => setForm({ ...form, for_sale: e.target.checked })}
          />
          List this project for sale
        </label>

        <button
          type="submit"
          disabled={create.isPending || !form.title}
          className="btn-primary w-full"
        >
          {create.isPending ? "Creating…" : "Create project"}
        </button>

        {create.isError && (
          <p className="text-sm text-red-400">
            {(create.error as Error).message}
          </p>
        )}
        {create.isSuccess && (
          <p className="text-sm text-emerald-400">Project created.</p>
        )}
      </form>
    </div>
  );
}
