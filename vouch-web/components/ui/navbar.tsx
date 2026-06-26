"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useAuth } from "@/store/auth";
import { clearCompanyTokens } from "@/lib/auth";
import { githubOAuthURL } from "@/lib/utils";
import type { Company } from "@/types";

function Logo() {
  return (
    <svg width="30" height="30" viewBox="0 0 28 28" aria-hidden="true">
      <rect x="2" y="2" width="24" height="24" rx="8" fill="#C8F24C" />
      <rect x="9.1" y="3.4" width="3.4" height="9.2" rx="1.7" fill="#14181B" transform="rotate(-11 10.8 8)" />
      <rect x="15.5" y="3.4" width="3.4" height="9.2" rx="1.7" fill="#14181B" transform="rotate(11 17.2 8)" />
      <circle cx="14" cy="17" r="7.3" fill="#14181B" />
      <circle cx="11.1" cy="16.1" r="2.7" fill="#C8F24C" />
      <circle cx="16.9" cy="16.1" r="2.7" fill="#C8F24C" />
      <circle cx="14" cy="20.6" r="1" fill="#C8F24C" />
    </svg>
  );
}

export function Navbar() {
  const user = useAuth((s) => s.user);
  const logout = useAuth((s) => s.logout);
  const [company, setCompany] = useState<Company | null>(null);

  useEffect(() => {
    const stored = localStorage.getItem("vouch_company");
    if (stored) { try { setCompany(JSON.parse(stored)); } catch { /* ignore */ } }
  }, []);

  return (
    <header className="sticky top-0 z-20 border-b border-line bg-paper/85 backdrop-blur">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link href="/" className="flex items-center gap-2">
          <Logo />
          <span className="font-display text-xl font-bold tracking-tight text-ink">
            vouch
          </span>
        </Link>

        <div className="flex items-center gap-1 font-mono text-xs font-bold uppercase tracking-wide text-ink/70">
          <Link href="/discover" className="rounded-full px-3 py-2 hover:bg-ink/5">
            Discover
          </Link>
          <Link href="/problems" className="rounded-full px-3 py-2 hover:bg-ink/5">
            Demand
          </Link>

          {company && !user ? (
            <>
              <Link href="/company/dashboard" className="rounded-full px-3 py-2 hover:bg-ink/5">
                Dashboard
              </Link>
              <span className="rounded-full px-3 py-2 text-ink/70">{company.name}</span>
              <button
                onClick={() => { clearCompanyTokens(); setCompany(null); }}
                className="ml-1 rounded-full px-3 py-2 hover:bg-ink/5 text-ink/50"
              >
                Sign out
              </button>
            </>
          ) : user ? (
            <>
              <Link
                href={user.role === "admin" ? "/admin" : "/dashboard"}
                className="rounded-full px-3 py-2 hover:bg-ink/5"
              >
                {user.role === "admin" ? "Admin" : "Dashboard"}
              </Link>
              <Link href="/profile" className="flex items-center gap-2 rounded-full px-3 py-2 hover:bg-ink/5">
                {user.avatar_url && (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={user.avatar_url} alt={user.username} className="h-5 w-5 rounded-full" />
                )}
                {user.username}
              </Link>
              <button
                onClick={logout}
                className="ml-1 rounded-full px-3 py-2 hover:bg-ink/5 text-ink/50"
              >
                Sign out
              </button>
            </>
          ) : (
            <div className="flex items-center gap-2">
              <Link href="/company/login" className="rounded-full px-3 py-2 hover:bg-ink/5 text-ink/60">
                Company
              </Link>
              <a
                href={githubOAuthURL()}
                className="ml-1 rounded-full bg-accent px-4 py-2 font-bold text-ink shadow-hard transition hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-hard-lg"
              >
                Sign in
              </a>
            </div>
          )}
        </div>
      </nav>
    </header>
  );
}
