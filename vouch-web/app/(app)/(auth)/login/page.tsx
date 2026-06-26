"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useRef, useState } from "react";

import { githubOAuthURL } from "@/lib/utils";
import { useAuth } from "@/store/auth";

function LoginInner() {
  const params = useSearchParams();
  const router = useRouter();
  const loginWithGitHub = useAuth((s) => s.loginWithGitHub);
  const [error, setError] = useState<string | null>(null);
  const handled = useRef(false);

  useEffect(() => {
    const code = params.get("code");
    if (!code || handled.current) return;
    handled.current = true;
    loginWithGitHub(code)
      .then((user) => router.replace(`/builder/${user.username}`))
      .catch(() => setError("GitHub sign-in failed. Please try again."));
  }, [params, loginWithGitHub, router]);

  const exchanging = !!params.get("code") && !error;

  return (
    <div className="mx-auto max-w-sm py-20 text-center">
      <h1 className="text-2xl font-bold">Sign in to Vouch</h1>
      <p className="mt-2 text-sm text-ink/60">
        Your GitHub account is the foundation of your verified score.
      </p>

      {exchanging ? (
        <p className="mt-8 text-ink/60">Completing sign-in…</p>
      ) : (
        <a href={githubOAuthURL()} className="btn-primary mt-8 w-full">
          Continue with GitHub
        </a>
      )}

      {error && <p className="mt-4 text-sm text-red-400">{error}</p>}
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={<p className="py-20 text-center text-ink/60">Loading…</p>}>
      <LoginInner />
    </Suspense>
  );
}
