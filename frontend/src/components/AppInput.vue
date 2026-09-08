<script setup lang="ts">
import { useId } from "vue";

const model = defineModel<string>({ required: true });

withDefaults(
  defineProps<{
    label: string;
    type?: string;
    placeholder?: string;
    required?: boolean;
    minlength?: number;
    hint?: string;
    error?: string;
    autocomplete?: string;
  }>(),
  {
    type: "text",
    required: false,
  },
);

const id = useId();
</script>

<template>
  <div class="field">
    <label :for="id">{{ label }}</label>
    <input
      :id="id"
      v-model="model"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :minlength="minlength"
      :autocomplete="autocomplete"
      class="input"
      :class="{ 'has-error': error }"
      :aria-invalid="error ? 'true' : undefined"
      :aria-describedby="error || hint ? `${id}-note` : undefined"
    />
    <span v-if="error" :id="`${id}-note`" class="error">{{ error }}</span>
    <span v-else-if="hint" :id="`${id}-note`" class="hint">{{ hint }}</span>
  </div>
</template>
