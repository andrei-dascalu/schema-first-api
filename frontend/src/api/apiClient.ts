import { client } from "./generated/client.gen";
import { useAuthStore } from "../stores/auth";

client.setConfig({
  baseUrl: "http://localhost:8080",
  auth: () => useAuthStore().token ?? undefined,
});

export * from "./generated/sdk.gen";
export * from "./generated/types.gen";
