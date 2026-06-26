import { create } from "zustand";

import { api } from "@/lib/api";
import { clearTokens, getAccessToken, storeTokens } from "@/lib/auth";
import type { User } from "@/types";

interface AuthState {
  user: User | null;
  status: "idle" | "loading" | "authenticated" | "unauthenticated";
  loginWithGitHub: (code: string) => Promise<User>;
  hydrate: () => Promise<void>;
  logout: () => void;
  setUser: (user: User) => void;
}

export const useAuth = create<AuthState>((set) => ({
  user: null,
  status: "idle",

  loginWithGitHub: async (code: string) => {
    set({ status: "loading" });
    const res = await api.loginWithGitHub(code);
    storeTokens(res.tokens);
    set({ user: res.user, status: "authenticated" });
    return res.user;
  },

  hydrate: async () => {
    const token = getAccessToken();
    if (!token) {
      set({ status: "unauthenticated" });
      return;
    }
    // We have a token; trust it until a request proves otherwise. The username
    // is unknown here, so callers fetch fresh profile data as needed.
    set({ status: "authenticated" });
  },

  logout: () => {
    clearTokens();
    set({ user: null, status: "unauthenticated" });
  },

  setUser: (user: User) => set({ user, status: "authenticated" }),
}));
