"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { api } from "@/lib/api";
import { storeCompanyTokens } from "@/lib/auth";

export default function CompanyLoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res = await api.companyLogin(email, password);
      storeCompanyTokens(res.tokens);
      localStorage.setItem("vouch_company", JSON.stringify(res.company));
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
        <p className="mt-1 text-sm text-ink/60">Sign in to your company account</p>
      </div>

      <form onSubmit={onSubmit} className="card space-y-4">
        <label className="block">
          <span className="mb-1 block text-sm text-ink/60">Work email</span>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
          />
        </label>
        <label className="block">
          <span className="mb-1 block text-sm text-ink/60">Password</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="w-full rounded-lg border border-line bg-panel px-3 py-2 text-sm outline-none focus:border-accent"
          />
        </label>

        <button type="submit" disabled={loading} className="btn-primary w-full">
          {loading ? "Signing in…" : "Sign in"}
        </button>

        {error && <p className="text-sm text-red-400">{error}</p>}
      </form>

      <p className="text-center text-sm text-ink/50">
        No account yet?{" "}
        <Link href="/company/register" className="text-accent-ink underline">
          Create one
        </Link>
      </p>
      <p className="text-center text-sm text-ink/50">
        Builder?{" "}
        <Link href="/login" className="text-accent-ink underline">
          Sign in with GitHub
        </Link>
      </p>
    </div>
  );
}
