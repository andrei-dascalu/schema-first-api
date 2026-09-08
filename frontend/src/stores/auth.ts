import { defineStore } from "pinia";

const STORAGE_KEY = "spec-first.token";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem(STORAGE_KEY) as string | null,
  }),
  getters: {
    isAuthenticated: (state) => state.token !== null,
  },
  actions: {
    setToken(token: string) {
      this.token = token;
      localStorage.setItem(STORAGE_KEY, token);
    },
    logout() {
      this.token = null;
      localStorage.removeItem(STORAGE_KEY);
    },
  },
});
