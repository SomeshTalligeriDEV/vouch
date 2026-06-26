// Typed API client. ALL fetch calls to vouch-api live here.

import type {
  AdminStats,
  ApiEnvelope,
  AuthResponse,
  BuilderScore,
  Company,
  CompanyAuthResponse,
  Paginated,
  Problem,
  Project,
  Review,
  TokenPair,
  User,
} from "@/types";
import {
  clearTokens,
  clearCompanyTokens,
  getAccessToken,
  getCompanyAccessToken,
  getCompanyRefreshToken,
  getRefreshToken,
  storeCompanyTokens,
  storeTokens,
} from "@/lib/auth";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export class ApiError extends Error {
  code: string;
  status: number;
  constructor(status: number, code: string, message: string) {
    super(message);
    this.code = code;
    this.status = status;
  }
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  auth?: boolean;
  companyAuth?: boolean;
  query?: Record<string, string | number | boolean | undefined>;
  retryOn401?: boolean;
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, auth = false, companyAuth = false, query, retryOn401 = true } = opts;

  const url = new URL(BASE_URL + path);
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined && v !== "") url.searchParams.set(k, String(v));
    }
  }

  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (companyAuth) {
    const token = getCompanyAccessToken();
    if (token) headers["Authorization"] = `Bearer ${token}`;
  } else if (auth) {
    const token = getAccessToken();
    if (token) headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(url.toString(), {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });

  // Transparently refresh once on a 401, then replay the original request.
  if (res.status === 401 && retryOn401) {
    if (companyAuth) {
      const refreshed = await tryCompanyRefresh();
      if (refreshed) {
        return request<T>(path, { ...opts, retryOn401: false });
      }
      clearCompanyTokens();
    } else if (auth) {
      const refreshed = await tryRefresh();
      if (refreshed) {
        return request<T>(path, { ...opts, retryOn401: false });
      }
      clearTokens();
    }
  }

  const envelope = (await res.json()) as ApiEnvelope<T>;
  if (!res.ok || !envelope.success) {
    const err = envelope.error;
    throw new ApiError(
      res.status,
      err?.code ?? "unknown",
      err?.message ?? "request failed",
    );
  }
  return envelope.data as T;
}

async function requestList<T>(
  path: string,
  opts: RequestOptions = {},
): Promise<Paginated<T>> {
  const { method = "GET", auth = false, query } = opts;
  const url = new URL(BASE_URL + path);
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined && v !== "") url.searchParams.set(k, String(v));
    }
  }
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (auth) {
    const token = getAccessToken();
    if (token) headers["Authorization"] = `Bearer ${token}`;
  }
  const res = await fetch(url.toString(), { method, headers, cache: "no-store" });
  const envelope = (await res.json()) as ApiEnvelope<T[]>;
  if (!res.ok || !envelope.success) {
    throw new ApiError(
      res.status,
      envelope.error?.code ?? "unknown",
      envelope.error?.message ?? "request failed",
    );
  }
  return {
    items: envelope.data ?? [],
    page: envelope.meta?.page ?? 1,
    limit: envelope.meta?.limit ?? 20,
    total: envelope.meta?.total ?? 0,
  };
}

let refreshInFlight: Promise<boolean> | null = null;

async function tryRefresh(): Promise<boolean> {
  if (refreshInFlight) return refreshInFlight;
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;

  refreshInFlight = (async () => {
    try {
      const res = await fetch(`${BASE_URL}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
      const env = (await res.json()) as ApiEnvelope<TokenPair>;
      if (res.ok && env.success && env.data) {
        storeTokens(env.data);
        return true;
      }
      return false;
    } catch {
      return false;
    } finally {
      refreshInFlight = null;
    }
  })();
  return refreshInFlight;
}

let companyRefreshInFlight: Promise<boolean> | null = null;

async function tryCompanyRefresh(): Promise<boolean> {
  if (companyRefreshInFlight) return companyRefreshInFlight;
  const refreshToken = getCompanyRefreshToken();
  if (!refreshToken) return false;

  companyRefreshInFlight = (async () => {
    try {
      const res = await fetch(`${BASE_URL}/companies/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
      const env = (await res.json()) as ApiEnvelope<TokenPair>;
      if (res.ok && env.success && env.data) {
        storeCompanyTokens(env.data);
        return true;
      }
      return false;
    } catch {
      return false;
    } finally {
      companyRefreshInFlight = null;
    }
  })();
  return companyRefreshInFlight;
}

export interface PresignResult {
  upload_url: string;
  public_url: string;
  key: string;
  expires_in: number;
}

/**
 * uploadFile presigns a direct-to-R2 upload, PUTs the file, and returns the
 * resulting public URL.
 */
export async function uploadFile(file: File): Promise<string> {
  const { upload_url, public_url } = await api.presignUpload(file.type);
  const res = await fetch(upload_url, {
    method: "PUT",
    headers: { "Content-Type": file.type },
    body: file,
  });
  if (!res.ok) {
    throw new ApiError(res.status, "upload_failed", "failed to upload file");
  }
  return public_url;
}

export interface ProjectFilters {
  status?: string;
  tag?: string;
  search?: string;
  sort?: string;
  for_sale?: boolean;
  page?: number;
  limit?: number;
}

export interface ProblemFilters {
  status?: string;
  tag?: string;
  search?: string;
  sort?: string;
  page?: number;
  limit?: number;
}

export const api = {
  // Auth
  loginWithGitHub: (code: string) =>
    request<AuthResponse>("/auth/github", { method: "POST", body: { code } }),
  logout: (refresh_token: string) =>
    request<unknown>("/auth/logout", { method: "POST", body: { refresh_token } }),

  // Users
  getMe: () => request<User>("/users/me", { auth: true }),
  getUser: (username: string) => request<User>(`/users/${username}`),
  updateMe: (input: Partial<User>) =>
    request<User>("/users/me", { method: "PATCH", body: input, auth: true }),
  connectStripe: (code: string) =>
    request<unknown>("/users/me/stripe", { method: "POST", body: { code }, auth: true }),

  // Projects
  listProjects: (filters: ProjectFilters = {}) =>
    requestList<Project>("/projects", { query: filters as Record<string, string | number | boolean | undefined> }),
  getProject: (id: string) => request<Project>(`/projects/${id}`),
  createProject: (input: Partial<Project>) =>
    request<Project>("/projects", { method: "POST", body: input, auth: true }),
  updateProject: (id: string, input: Partial<Project>) =>
    request<Project>(`/projects/${id}`, { method: "PATCH", body: input, auth: true }),
  archiveProject: (id: string) =>
    request<unknown>(`/projects/${id}`, { method: "DELETE", auth: true }),

  // Scores
  getScore: (username: string) => request<BuilderScore>(`/scores/${username}`),
  leaderboard: (limit = 25) =>
    request<BuilderScore[]>("/scores", { query: { limit } }),
  recalculateScore: () =>
    request<unknown>("/scores/recalculate", { method: "POST", auth: true }),

  // Problems
  listProblems: (filters: ProblemFilters = {}) =>
    requestList<Problem>("/problems", { query: filters as Record<string, string | number | boolean | undefined> }),
  getProblem: (id: string) => request<Problem>(`/problems/${id}`),
  createProblem: (input: Partial<Problem>) =>
    request<Problem>("/problems", { method: "POST", body: input, auth: true }),
  claimProblem: (id: string) =>
    request<Problem>(`/problems/${id}/claim`, { method: "POST", auth: true }),
  upvoteProblem: (id: string) =>
    request<Problem>(`/problems/${id}/upvote`, { method: "POST", auth: true }),

  // Companies
  companyRegister: (input: { name: string; email: string; password: string; website?: string; size?: string }) =>
    request<CompanyAuthResponse>("/companies/register", { method: "POST", body: input }),
  companyLogin: (email: string, password: string) =>
    request<CompanyAuthResponse>("/companies/login", { method: "POST", body: { email, password } }),
  companyRefresh: (refresh_token: string) =>
    request<TokenPair>("/companies/refresh", { method: "POST", body: { refresh_token } }),
  companyRefreshLogout: (refresh_token: string) =>
    request<unknown>("/companies/logout", { method: "POST", body: { refresh_token } }),
  getCompanyMe: () => request<Company>("/companies/me", { companyAuth: true }),
  updateCompanyMe: (input: Partial<Company>) =>
    request<Company>("/companies/me", { method: "PATCH", body: input, companyAuth: true }),
  getCompanyBySlug: (slug: string) => request<Company>(`/companies/${slug}`),

  // Admin
  adminStats: () => request<AdminStats>("/admin/stats", { auth: true }),
  adminListCompanies: (page = 1) =>
    requestList<Company>("/admin/companies", { auth: true, query: { page } }),

  // Uploads
  presignUpload: (contentType: string) =>
    request<PresignResult>("/uploads/presign", {
      method: "POST",
      body: { content_type: contentType },
      auth: true,
    }),

  // Reviews
  createReview: (input: { project_id: string; rating: number; body: string }) =>
    request<Review>("/reviews", { method: "POST", body: input, auth: true }),
  listReviews: (projectId: string, page = 1) =>
    requestList<Review>(`/reviews/project/${projectId}`, { query: { page } }),
};
