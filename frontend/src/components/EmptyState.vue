<script setup lang="ts">
defineProps<{ title: string; message: string }>();
</script>

<template>
  <div class="empty-state">
    <h3>{{ title }}</h3>
    <p>{{ message }}</p>
    <div v-if="$slots.default" style="margin-top: 14px">
      <slot />
    </div>
  </div>
</template>
