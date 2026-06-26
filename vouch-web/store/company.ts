import { create } from "zustand";

import { api, ApiError } from "@/lib/api";
import {
  clearCompanyTokens,
  getCompanyAccessToken,
  getCompanyRefreshToken,
  storeCompanyTokens,
} from "@/lib/auth";
import type { Company } from "@/types";

interface CompanyAuthState {
  company: Company | null;
  status: "idle" | "loading" | "authenticated" | "unauthenticated";
  login: (email: string, password: string) => Promise<Company>;
  register: (input: { name: string; email: string; password: string; size?: string; website?: string }) => Promise<Company>;
  hydrate: () => Promise<void>;
  logout: () => void;
  setCompany: (c: Company) => void;
}

export const useCompanyAuth = create<CompanyAuthState>((set) => ({
  company: null,
  status: "idle",

  login: async (email, password) => {
    set({ status: "loading" });
    const res = await api.companyLogin(email, password);
    storeCompanyTokens(res.tokens);
    localStorage.setItem("vouch_company", JSON.stringify(res.company));
    set({ company: res.company, status: "authenticated" });
    return res.company;
  },

  register: async (input) => {
    set({ status: "loading" });
    const res = await api.companyRegister(input);
    storeCompanyTokens(res.tokens);
    localStorage.setItem("vouch_company", JSON.stringify(res.company));
    set({ company: res.company, status: "authenticated" });
    return res.company;
  },

  hydrate: async () => {
    const token = getCompanyAccessToken();
    if (!token) {
      set({ status: "unauthenticated" });
      return;
    }
    try {
      const company = await api.getCompanyMe();
      localStorage.setItem("vouch_company", JSON.stringify(company));
      set({ company, status: "authenticated" });
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        clearCompanyTokens();
        localStorage.removeItem("vouch_company");
        set({ company: null, status: "unauthenticated" });
      } else {
        // Optimistic: keep session on network errors.
        const stored = localStorage.getItem("vouch_company");
        if (stored) {
          try {
            set({ company: JSON.parse(stored), status: "authenticated" });
          } catch {
            set({ status: "unauthenticated" });
          }
        } else {
          set({ status: "unauthenticated" });
        }
      }
    }
  },

  logout: () => {
    const refreshToken = getCompanyRefreshToken();
    if (refreshToken) {
      api.companyRefreshLogout(refreshToken).catch(() => {/* best effort */});
    }
    clearCompanyTokens();
    localStorage.removeItem("vouch_company");
    set({ company: null, status: "unauthenticated" });
  },

  setCompany: (company) => set({ company, status: "authenticated" }),
}));
