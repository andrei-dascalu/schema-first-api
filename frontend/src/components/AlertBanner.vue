<script setup lang="ts">
defineProps<{ message: string }>();
</script>

<template>
  <p class="alert" role="alert">{{ message }}</p>
</template>
