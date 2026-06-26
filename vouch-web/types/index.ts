// TypeScript types mirroring the Go domain structs in vouch-api.

export type Role = "builder" | "user" | "admin";

export interface User {
  id: string;
  email: string;
  username: string;
  name: string;
  bio: string;
  avatar_url: string;
  github_id: number;
  github_login: string;
  role: Role;
  is_verified: boolean;
  website_url: string;
  twitter_handle: string;
  created_at: string;
  updated_at: string;
}

export type ProjectStatus = "draft" | "live" | "acquired" | "archived";

export interface Project {
  id: string;
  builder_id: string;
  title: string;
  slug: string;
  tagline: string;
  description: string;
  logo_url: string;
  live_url: string;
  repo_url: string;
  payment_link: string;
  tags: string[];
  status: ProjectStatus;
  for_sale: boolean;
  ask_price: number;
  verified_users: number;
  mrr: number;
  review_count: number;
  average_rating: number;
  created_at: string;
  updated_at: string;
}

export type Tier = "Bronze" | "Silver" | "Gold" | "Platinum" | "24 Karat";

export interface ScoreBreakdown {
  user: number;
  revenue: number;
  impact: number;
  velocity: number;
}

export interface BuilderScore {
  id: string;
  builder_id: string;
  username: string;
  name: string;
  avatar_url: string;
  total_score: number;
  tier: Tier;
  breakdown: ScoreBreakdown;
  stripe_verified: boolean;
  stripe_multiplier: number;
  calculated_at: string;
  updated_at: string;
}

export type ProblemStatus = "open" | "claimed" | "shipped" | "cancelled";

export interface Problem {
  id: string;
  poster_id: string;
  claimed_by?: string;
  shipped_project_id?: string;
  title: string;
  slug: string;
  description: string;
  tags: string[];
  budget_min: number;
  budget_max: number;
  status: ProblemStatus;
  upvotes: number;
  created_at: string;
  updated_at: string;
}

export interface Review {
  id: string;
  project_id: string;
  reviewer_id: string;
  reviewer_username: string;
  rating: number;
  body: string;
  verified_purchase: boolean;
  created_at: string;
  updated_at: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface AuthResponse {
  user: User;
  tokens: TokenPair;
}

export interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
  error?: { code: string; message: string };
  meta?: { page: number; limit: number; total: number };
}

export interface Paginated<T> {
  items: T[];
  page: number;
  limit: number;
  total: number;
}
