import { setActivePinia, createPinia } from "pinia";
import { beforeEach, describe, expect, it } from "vitest";
import { useAuthStore } from "./auth";

describe("auth store", () => {
  beforeEach(() => {
    localStorage.clear();
    setActivePinia(createPinia());
  });

  it("starts unauthenticated", () => {
    const auth = useAuthStore();
    expect(auth.isAuthenticated).toBe(false);
  });

  it("persists the token on setToken and clears it on logout", () => {
    const auth = useAuthStore();
    auth.setToken("abc.def.ghi");
    expect(auth.isAuthenticated).toBe(true);
    expect(localStorage.getItem("spec-first.token")).toBe("abc.def.ghi");

    auth.logout();
    expect(auth.isAuthenticated).toBe(false);
    expect(localStorage.getItem("spec-first.token")).toBeNull();
  });
});
