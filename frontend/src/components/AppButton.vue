<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: "primary" | "secondary" | "ghost" | "danger" | "danger-solid";
    type?: "button" | "submit";
    loading?: boolean;
    disabled?: boolean;
    size?: "md" | "sm";
  }>(),
  {
    variant: "primary",
    type: "button",
    loading: false,
    disabled: false,
    size: "md",
  },
);
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="btn"
    :class="[
      variant === 'secondary' && 'btn-secondary',
      variant === 'ghost' && 'btn-ghost',
      variant === 'danger' && 'btn-danger',
      variant === 'danger-solid' && 'btn-danger btn-solid',
      size === 'sm' && 'btn-sm',
    ]"
  >
    <span v-if="loading" class="spinner" aria-hidden="true" />
    <slot />
  </button>
</template>
