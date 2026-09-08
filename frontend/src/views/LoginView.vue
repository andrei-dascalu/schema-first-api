<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "../api/apiClient";
import { useAuthStore } from "../stores/auth";
import AppInput from "../components/AppInput.vue";
import AppButton from "../components/AppButton.vue";
import AlertBanner from "../components/AlertBanner.vue";

const email = ref("");
const password = ref("");
const error = ref("");
const submitting = ref(false);
const router = useRouter();
const auth = useAuthStore();

async function onSubmit() {
  error.value = "";
  submitting.value = true;
  const { data, error: apiError } = await login({
    body: { email: email.value, password: password.value },
  });
  submitting.value = false;
  if (apiError) {
    error.value = "message" in apiError ? apiError.message : "That email and password don't match.";
    return;
  }
  auth.setToken(data.token);
  router.push({ name: "events" });
}
</script>

<template>
  <div class="auth-shell">
    <aside class="auth-brand">
      <p class="auth-brand-mark">Gather</p>
      <p>
        Keep every event's guest list in one place — who's invited, who's coming,
        and who you're still waiting on.
      </p>
    </aside>
    <div class="auth-form-wrap">
      <form class="auth-form" @submit.prevent="onSubmit">
        <h1>Log in</h1>
        <AppInput
          v-model="email"
          label="Email"
          type="email"
          autocomplete="email"
          required
        />
        <AppInput
          v-model="password"
          label="Password"
          type="password"
          autocomplete="current-password"
          required
        />
        <AlertBanner v-if="error" :message="error" />
        <AppButton type="submit" :loading="submitting">
          {{ submitting ? "Logging in…" : "Log in" }}
        </AppButton>
        <RouterLink to="/register" class="switch-link">Need an account? Create one</RouterLink>
      </form>
    </div>
  </div>
</template>
