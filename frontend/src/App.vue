<script setup lang="ts">
import { useRoute } from "vue-router";
import { useAuthStore } from "./stores/auth";

const auth = useAuthStore();
const route = useRoute();
</script>

<template>
  <header class="site-header">
    <div class="site-header-inner">
      <RouterLink :to="auth.isAuthenticated ? '/events' : '/login'" class="brand">
        Gather
      </RouterLink>
      <nav v-if="auth.isAuthenticated" class="site-nav">
        <RouterLink to="/events" :class="{ 'is-current': route.name === 'events' }">
          Events
        </RouterLink>
        <button type="button" class="btn btn-ghost btn-sm" @click="auth.logout()">
          Log out
        </button>
      </nav>
    </div>
  </header>
  <RouterView />
</template>

<style scoped>
.site-header {
  border-bottom: 1px solid var(--border);
}

.site-header-inner {
  max-width: 720px;
  margin: 0 auto;
  padding: 18px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.brand {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--ink);
  text-decoration: none;
  letter-spacing: -0.2px;
}

.site-nav {
  display: flex;
  align-items: center;
  gap: 18px;
}

.site-nav a {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-soft);
  text-decoration: none;
  transition: color 0.15s ease;
}

.site-nav a:hover,
.site-nav a.is-current {
  color: var(--ink);
}
</style>
