"use client";

import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { api, type ProblemFilters, type ProjectFilters } from "@/lib/api";

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
