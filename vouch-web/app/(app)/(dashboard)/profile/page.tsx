"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { api } from "@/lib/api";
import { useAuth } from "@/store/auth";

export default function ProfilePage() {
  const user = useAuth((s) => s.user);
  const setUser = useAuth((s) => s.setUser);
  const router = useRouter();

  const [form, setForm] = useState({
    name: "",
    bio: "",
    website_url: "",
    twitter_handle: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (!user) { router.replace("/login"); return; }
    setForm({
      name: user.name ?? "",
      bio: user.bio ?? "",
      website_url: user.website_url ?? "",
      twitter_handle: user.twitter_handle ?? "",
    });
  }, [user, router]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSaved(false);
    try {
      const updated = await api.updateMe(form);
      setUser(updated);
      setSaved(true);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  if (!user) return null;

  const field = (label: string, key: keyof typeof form, placeholder = "") => (
    <label className="block">
      <span className="mb-1 block text-sm text-ink/60">{label}</span>
      <input
        type="text"
        value={form[key]}
        onChange={(e) => setForm({ ...form, [key]: e.target.value })}
        placeholder={placeholder}
        className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
      />
    </label>
  );

  return (
    <div className="mx-auto max-w-xl space-y-6">
      <div className="flex items-center gap-4">
        {user.avatar_url && (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={user.avatar_url} alt={user.username} className="h-16 w-16 rounded-full border border-line" />
        )}
        <div>
          <h1 className="text-2xl font-bold">{user.name || user.username}</h1>
          <p className="text-sm text-ink/60">@{user.username} · {user.email}</p>
        </div>
      </div>

      <form onSubmit={onSubmit} className="card space-y-4">
        <h2 className="font-semibold">Edit profile</h2>
        {field("Display name", "name", "Your full name")}
        <label className="block">
          <span className="mb-1 block text-sm text-ink/60">Bio</span>
          <textarea
            value={form.bio}
            onChange={(e) => setForm({ ...form, bio: e.target.value })}
            rows={3}
            placeholder="Tell the market who you are and what you build."
            className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
          />
        </label>
        {field("Website", "website_url", "https://yoursite.com")}
        {field("Twitter / X handle", "twitter_handle", "@handle")}

        <button type="submit" disabled={loading} className="btn-primary w-full">
          {loading ? "Saving…" : "Save changes"}
        </button>

        {saved && <p className="text-sm text-emerald-400">Profile updated.</p>}
        {error && <p className="text-sm text-red-400">{error}</p>}
      </form>

      <div className="card">
        <h2 className="font-semibold mb-1">GitHub</h2>
        <p className="text-sm text-ink/60">
          Connected as <span className="font-mono">@{user.github_login}</span>
        </p>
      </div>
    </div>
  );
}
