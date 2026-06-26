import Link from "next/link";
import { cn } from "@/lib/utils";

interface Crumb {
  label: string;
  href?: string;
}

interface BreadcrumbsProps {
  crumbs: Crumb[];
  className?: string;
}

export function Breadcrumbs({ crumbs, className }: BreadcrumbsProps) {
  return (
    <nav aria-label="Breadcrumb" className={cn("flex items-center gap-1.5 text-sm", className)}>
      <ol className="flex items-center gap-1.5">
        {crumbs.map((crumb, i) => (
          <li key={i} className="flex items-center gap-1.5">
            {i > 0 && (
              <span aria-hidden className="text-muted-foreground">/</span>
            )}
            {crumb.href ? (
              <Link
                href={crumb.href}
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                {crumb.label}
              </Link>
            ) : (
              <span className="text-foreground font-medium" aria-current="page">
                {crumb.label}
              </span>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
}
