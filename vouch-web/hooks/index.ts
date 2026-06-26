"use client";

import { useEffect } from "react";
import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { api, type ProblemFilters, type ProjectFilters } from "@/lib/api";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

// Projects
export function useProjects(filters: ProjectFilters = {}) {
  return useQuery({
    queryKey: ["projects", filters],
    queryFn: () => api.listProjects(filters),
  });
}

export function useProject(id: string) {
  return useQuery({
    queryKey: ["project", id],
    queryFn: () => api.getProject(id),
    enabled: !!id,
  });
}

export function useCreateProject() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.createProject,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["projects"] }),
  });
}

// Problems
export function useProblems(filters: ProblemFilters = {}) {
  return useQuery({
    queryKey: ["problems", filters],
    queryFn: () => api.listProblems(filters),
  });
}

export function useProblem(id: string) {
  return useQuery({
    queryKey: ["problem", id],
    queryFn: () => api.getProblem(id),
    enabled: !!id,
  });
}

export function useClaimProblem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.claimProblem(id),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ["problem", id] });
      qc.invalidateQueries({ queryKey: ["problems"] });
    },
  });
}

export function useUpvoteProblem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.upvoteProblem(id),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ["problem", id] });
      qc.invalidateQueries({ queryKey: ["problems"] });
    },
  });
}

// Builder + score
export function useBuilder(username: string) {
  return useQuery({
    queryKey: ["builder", username],
    queryFn: () => api.getUser(username),
    enabled: !!username,
  });
}

export function useScore(username: string) {
  return useQuery({
    queryKey: ["score", username],
    queryFn: () => api.getScore(username),
    enabled: !!username,
  });
}

export function useLeaderboard(limit = 25) {
  return useQuery({
    queryKey: ["leaderboard", limit],
    queryFn: () => api.leaderboard(limit),
  });
}

export function useReviews(projectId: string) {
  return useQuery({
    queryKey: ["reviews", projectId],
    queryFn: () => api.listReviews(projectId),
    enabled: !!projectId,
  });
}

// useLeaderboardSSE subscribes to real-time leaderboard updates via SSE and
// invalidates the React Query cache whenever the server publishes a new event.
export function useLeaderboardSSE() {
  const qc = useQueryClient();
  useEffect(() => {
    const url = BASE_URL.replace("/api/v1", "") + "/api/v1/sse/leaderboard";
    const es = new EventSource(url);
    es.addEventListener("leaderboard", () => {
      qc.invalidateQueries({ queryKey: ["leaderboard"] });
    });
    es.onerror = () => es.close();
    return () => es.close();
  }, [qc]);
}

// useScoreSSE subscribes to real-time score updates for a specific builder.
export function useScoreSSE(username: string) {
  const qc = useQueryClient();
  useEffect(() => {
    if (!username) return;
    const url = `${BASE_URL.replace("/api/v1", "")}/api/v1/sse/scores/${username}`;
    const es = new EventSource(url);
    es.addEventListener("score", () => {
      qc.invalidateQueries({ queryKey: ["score", username] });
      qc.invalidateQueries({ queryKey: ["leaderboard"] });
    });
    es.onerror = () => es.close();
    return () => es.close();
  }, [username, qc]);
}
