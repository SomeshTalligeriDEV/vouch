"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import type { Company } from "@/types";
import { COMPANY_SIZE_LABELS } from "@/lib/constants";

export default function CompanySettingsPage() {
  const router = useRouter();
  const [company, setCompany] = useState<Company | null>(null);
  const [form, setForm] = useState({
    name: "",
    website: "",
    description: "",
    size: "",
    logo_url: "",
  });
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const stored = localStorage.getItem("vouch_company");
    if (!stored) { router.replace("/company/login"); return; }
    try {
      const c: Company = JSON.parse(stored);
      setCompany(c);
      setForm({
        name: c.name ?? "",
        website: c.website ?? "",
        description: c.description ?? "",
        size: c.size ?? "",
        logo_url: c.logo_url ?? "",
      });
    } catch { router.replace("/company/login"); }
  }, [router]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    setSaved(false);
    try {
      const updated = await api.updateCompanyMe(form);
      localStorage.setItem("vouch_company", JSON.stringify(updated));
      setCompany(updated);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setSaving(false);
    }
  };

  if (!company) return null;

  return (
    <main className="max-w-xl mx-auto px-4 py-10">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">Company Settings</h1>
        <p className="text-muted-foreground text-sm mt-1">{company.email}</p>
      </div>

      <form onSubmit={handleSave} className="space-y-5">
        <div>
          <label className="block text-sm font-medium mb-1">Company Name</label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Website</label>
          <input
            type="url"
            value={form.website}
            onChange={(e) => setForm({ ...form, website: e.target.value })}
            placeholder="https://example.com"
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Description</label>
          <textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            rows={3}
            placeholder="What does your company do?"
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary resize-none"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Company Size</label>
          <select
            value={form.size}
            onChange={(e) => setForm({ ...form, size: e.target.value })}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          >
            <option value="">Select size</option>
            {Object.entries(COMPANY_SIZE_LABELS).map(([k, v]) => (
              <option key={k} value={k}>{v}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Logo URL</label>
          <input
            type="url"
            value={form.logo_url}
            onChange={(e) => setForm({ ...form, logo_url: e.target.value })}
            placeholder="https://cdn.example.com/logo.png"
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </div>

        <div className="flex items-center gap-3 pt-2">
          <button
            type="submit"
            disabled={saving}
            className="rounded-lg bg-primary px-5 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 transition-colors"
          >
            {saving ? "Saving…" : "Save changes"}
          </button>
          {saved && <span className="text-sm text-emerald-500">Saved!</span>}
        </div>

        {error && <p className="text-sm text-destructive">{error}</p>}
      </form>
    </main>
  );
}
