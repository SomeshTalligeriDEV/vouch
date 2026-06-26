export const TIER_LABELS: Record<number, string> = {
  1: "Tier 1 — Elite",
  2: "Tier 2 — Pro",
  3: "Tier 3 — Builder",
  4: "Tier 4 — Rising",
  5: "Tier 5 — Newcomer",
};

export const TIER_COLORS: Record<number, string> = {
  1: "bg-yellow-100 text-yellow-800",
  2: "bg-violet-100 text-violet-800",
  3: "bg-blue-100 text-blue-800",
  4: "bg-green-100 text-green-800",
  5: "bg-gray-100 text-gray-700",
};

export const COMPANY_SIZE_LABELS: Record<string, string> = {
  solo: "Solo / 1 person",
  small: "2–10 employees",
  medium: "11–50 employees",
  large: "51–200 employees",
  enterprise: "200+ employees",
};

export const SCORE_RANGES: Record<string, { min: number; max: number; label: string }> = {
  elite: { min: 900, max: 1000, label: "Elite" },
  pro: { min: 700, max: 899, label: "Pro" },
  builder: { min: 500, max: 699, label: "Builder" },
  rising: { min: 300, max: 499, label: "Rising" },
  newcomer: { min: 0, max: 299, label: "Newcomer" },
};

export const MAX_BIO_LENGTH = 300;
export const MAX_PROJECT_TITLE_LENGTH = 80;
export const MAX_PROJECT_DESCRIPTION_LENGTH = 2000;
export const MAX_REVIEW_LENGTH = 1000;

export const PAGINATION_DEFAULT_LIMIT = 20;
export const LEADERBOARD_PAGE_SIZE = 50;
