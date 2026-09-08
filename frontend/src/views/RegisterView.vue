<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { registerUser } from "../api/apiClient";
import AppInput from "../components/AppInput.vue";
import AppButton from "../components/AppButton.vue";
import AlertBanner from "../components/AlertBanner.vue";

const email = ref("");
const password = ref("");
const name = ref("");
const error = ref("");
const submitting = ref(false);
const router = useRouter();

async function onSubmit() {
  error.value = "";
  submitting.value = true;
  const { error: apiError } = await registerUser({
    body: { email: email.value, password: password.value, name: name.value },
  });
  submitting.value = false;
  if (apiError) {
    error.value =
      "message" in apiError ? apiError.message : "Couldn't create that account.";
    return;
  }
  router.push({ name: "login" });
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
        <h1>Create your account</h1>
        <AppInput v-model="name" label="Name" autocomplete="name" required />
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
          autocomplete="new-password"
          :minlength="8"
          hint="At least 8 characters"
          required
        />
        <AlertBanner v-if="error" :message="error" />
        <AppButton type="submit" :loading="submitting">
          {{ submitting ? "Creating account…" : "Create account" }}
        </AppButton>
        <RouterLink to="/login" class="switch-link">Already have an account? Log in</RouterLink>
      </form>
    </div>
  </div>
</template>
