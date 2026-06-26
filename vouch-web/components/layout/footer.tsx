import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t bg-background">
      <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
        <div className="grid grid-cols-2 gap-8 sm:grid-cols-4">
          <div>
            <p className="text-sm font-semibold text-foreground">Platform</p>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li><Link href="/discover" className="hover:text-foreground transition-colors">Discover</Link></li>
              <li><Link href="/leaderboard" className="hover:text-foreground transition-colors">Leaderboard</Link></li>
              <li><Link href="/problems" className="hover:text-foreground transition-colors">Problems</Link></li>
            </ul>
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">Builders</p>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li><Link href="/dashboard" className="hover:text-foreground transition-colors">Dashboard</Link></li>
              <li><Link href="/auth/github" className="hover:text-foreground transition-colors">Sign in with GitHub</Link></li>
            </ul>
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">Companies</p>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li><Link href="/company/dashboard" className="hover:text-foreground transition-colors">Company Dashboard</Link></li>
              <li><Link href="/company/register" className="hover:text-foreground transition-colors">Post a Problem</Link></li>
            </ul>
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">Legal</p>
            <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
              <li><Link href="/privacy" className="hover:text-foreground transition-colors">Privacy</Link></li>
              <li><Link href="/terms" className="hover:text-foreground transition-colors">Terms</Link></li>
            </ul>
          </div>
        </div>
        <div className="mt-10 border-t pt-6 flex items-center justify-between">
          <p className="text-xs text-muted-foreground">
            © {new Date().getFullYear()} Vouch. Demand-first builder reputation.
          </p>
          <p className="text-xs text-muted-foreground">v0.1.0</p>
        </div>
      </div>
    </footer>
  );
}
