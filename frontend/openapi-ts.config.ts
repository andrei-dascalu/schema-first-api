import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "../docs/api_schema.yaml",
  output: "src/api/generated",
  plugins: ["@hey-api/client-fetch", "@hey-api/typescript", "@hey-api/sdk"],
});
