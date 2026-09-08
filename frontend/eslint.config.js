import js from "@eslint/js";
import tseslint from "typescript-eslint";
import pluginVue from "eslint-plugin-vue";
import vueParser from "vue-eslint-parser";
import globals from "globals";

export default tseslint.config(
  {
    ignores: ["dist/**", "src/api/generated/**"],
  },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  // "essential" only (not "recommended"/"strongly-recommended"): those tiers
  // add stylistic/whitespace rules that fight Prettier, which owns formatting here.
  ...pluginVue.configs["flat/essential"],
  {
    files: ["**/*.vue"],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tseslint.parser,
      },
    },
  },
  {
    languageOptions: {
      globals: globals.browser,
    },
  },
);
