import { create } from "zustand";

import { api, ApiError } from "@/lib/api";
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
    try {
      // Verify the token is real by fetching the actual user profile.
      // This catches expired tokens, revoked sessions, and forged JWTs.
      const user = await api.getMe();
      set({ user, status: "authenticated" });
    } catch (err) {
      // 401 means the token is expired/invalid — clear and force re-login.
      // Any other error (network down, 5xx) keeps the session alive optimistically
      // so a server hiccup doesn't log everyone out.
      if (err instanceof ApiError && err.status === 401) {
        clearTokens();
        set({ user: null, status: "unauthenticated" });
      } else {
        set({ status: "authenticated" });
      }
    }
  },

  logout: () => {
    clearTokens();
    set({ user: null, status: "unauthenticated" });
  },

  setUser: (user: User) => set({ user, status: "authenticated" }),
}));
