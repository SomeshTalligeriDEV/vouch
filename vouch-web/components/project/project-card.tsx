import Link from "next/link";

import type { Project } from "@/types";
import { formatCurrency } from "@/lib/utils";

export function ProjectCard({ project }: { project: Project }) {
  return (
    <Link
      href={`/projects/${project.id}`}
      className="card block transition hover:border-crimson/40"
    >
      <div className="flex items-start justify-between gap-3">
        <div>
          <h3 className="font-semibold">{project.title}</h3>
          <p className="mt-1 text-sm text-ink/60">{project.tagline}</p>
        </div>
        {project.for_sale && (
          <span className="tag border-gold bg-gold/30 text-crimson-dark">For sale</span>
        )}
      </div>
      <div className="mt-4 flex flex-wrap gap-4 text-sm text-ink/60">
        <span>{formatCurrency(project.mrr)}/mo MRR</span>
        <span>{project.verified_users} users</span>
        <span>
          ★ {project.average_rating.toFixed(1)} ({project.review_count})
        </span>
      </div>
      {project.tags.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1.5">
          {project.tags.slice(0, 4).map((t) => (
            <span key={t} className="tag">
              {t}
            </span>
          ))}
        </div>
      )}
    </Link>
  );
}
