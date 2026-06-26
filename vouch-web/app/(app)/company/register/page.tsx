"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { useCompanyAuth } from "@/store/company";

const SIZES = [
  { value: "1", label: "Solo founder" },
  { value: "2-10", label: "2–10 people" },
  { value: "11-50", label: "11–50 people" },
  { value: "51-200", label: "51–200 people" },
  { value: "200+", label: "200+ people" },
];

export default function CompanyRegisterPage() {
  const router = useRouter();
  const register = useCompanyAuth((s) => s.register);
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    website: "",
    size: "2-10",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (form.password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      await register(form);
      router.replace("/company/dashboard");
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const field = (label: string, key: keyof typeof form, type = "text", placeholder = "", required = true) => (
    <label className="block">
      <span className="mb-1 block text-sm font-medium">{label}</span>
      <input
        type={type}
        value={form[key]}
        onChange={(e) => setForm({ ...form, [key]: e.target.value })}
        placeholder={placeholder}
        required={required}
        className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
      />
    </label>
  );

  return (
    <div className="mx-auto max-w-md py-12 space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Create a company account</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Post real problems with budgets. Builders ship solutions and you get
          your first paying customer.
        </p>
      </div>

      <form onSubmit={onSubmit} className="rounded-xl border border-border bg-card p-6 space-y-4">
        {field("Company name *", "name", "text", "Acme Inc.")}
        {field("Work email *", "email", "email", "you@company.com")}
        {field("Password *", "password", "password", "Min 8 characters")}
        {field("Website", "website", "url", "https://company.com", false)}

        <label className="block">
          <span className="mb-1 block text-sm font-medium">Company size</span>
          <select
            value={form.size}
            onChange={(e) => setForm({ ...form, size: e.target.value })}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          >
            {SIZES.map((s) => (
              <option key={s.value} value={s.value}>{s.label}</option>
            ))}
          </select>
        </label>

        <button
          type="submit"
          disabled={loading || !form.name || !form.email || !form.password}
          className="w-full rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 transition-colors"
        >
          {loading ? "Creating account…" : "Create account"}
        </button>

        {error && <p className="text-sm text-destructive">{error}</p>}
      </form>

      <p className="text-center text-sm text-muted-foreground">
        Already have an account?{" "}
        <Link href="/company/login" className="text-primary hover:underline">
          Sign in
        </Link>
      </p>
      <p className="text-center text-sm text-muted-foreground">
        Are you a builder?{" "}
        <Link href="/login" className="text-primary hover:underline">
          Sign in with GitHub
        </Link>
      </p>
    </div>
  );
}
