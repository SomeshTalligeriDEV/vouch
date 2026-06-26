"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { useCompanyAuth } from "@/store/company";

export default function CompanyLoginPage() {
  const router = useRouter();
  const login = useCompanyAuth((s) => s.login);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await login(email, password);
      router.replace("/company/dashboard");
    } catch {
      setError("Invalid email or password.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-sm py-16 space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold">Company sign in</h1>
        <p className="mt-1 text-sm text-muted-foreground">Sign in to your company account</p>
      </div>

      <form onSubmit={onSubmit} className="rounded-xl border border-border bg-card p-6 space-y-4">
        <label className="block">
          <span className="mb-1 block text-sm font-medium">Work email</span>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoComplete="email"
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </label>
        <label className="block">
          <span className="mb-1 block text-sm font-medium">Password</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="current-password"
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </label>

        <button
          type="submit"
          disabled={loading || !email || !password}
          className="w-full rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 transition-colors"
        >
          {loading ? "Signing in…" : "Sign in"}
        </button>

        {error && <p className="text-sm text-destructive">{error}</p>}
      </form>

      <p className="text-center text-sm text-muted-foreground">
        No account yet?{" "}
        <Link href="/company/register" className="text-primary hover:underline">
          Create one
        </Link>
      </p>
      <p className="text-center text-sm text-muted-foreground">
        Builder?{" "}
        <Link href="/login" className="text-primary hover:underline">
          Sign in with GitHub
        </Link>
      </p>
    </div>
  );
}
