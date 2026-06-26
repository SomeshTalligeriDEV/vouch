import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-6 px-4 text-center">
      <div className="space-y-2">
        <p className="font-mono text-xs uppercase tracking-widest text-ink/40">404</p>
        <h1 className="font-display text-3xl font-bold text-ink">Page not found</h1>
        <p className="text-sm text-ink/60 max-w-sm">
          This page doesn&apos;t exist or was moved.
        </p>
      </div>
      <Link
        href="/"
        className="rounded-full bg-accent px-5 py-2 font-bold text-ink shadow-hard transition hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-hard-lg"
      >
        Go home
      </Link>
    </div>
  );
}
