<script setup lang="ts">
import { Loader2Icon } from 'lucide-vue-next'
import { defineProps, type FunctionalComponent } from 'vue';

defineProps({
  isLoading: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  type: {
    type: String as PropType<"button" | "submit" | "reset">,
    default: "",
  },
  kind: {
    type: String as PropType<"primary" | "secondary" | "danger">,
    default: "primary",
  },
  icon: {
    type: Object as PropType<FunctionalComponent>,
    default: null
  }
});
</script>

<template>
  <button :type="type" :disabled="disabled || loading" :class="kind">
    <Loader2Icon v-if="loading" class="icon-loader" />
    <component :is="icon" v-else-if="icon" />
    <span>
      <slot />
    </span>
  </button>
</template>

<style scoped>
button {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 1rem;
  border: 1px solid transparent;
  border-radius: 2px;
  font-size: var(--text-small);
  font-family: var(--font-family);
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

button svg {
  height: var(--text-small);
  width: var(--text-small);
}

button:hover {
  background-color: var(--color-primary-dark);
}

button:focus {
  outline: none;
  border-color: var(--color-accent);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.icon-loader {
  animation: spin 1s linear infinite;
}
  
button.primary {
  background-color: var(--color-primary);
}

button.primary:hover {
  background-color: var(--color-primary-dark);
}
  
button.secondary {
  background-color: unset;
  border: 1px solid var(--color-primary);
  color: var(--color-primary);
}
  
button.secondary:hover {
  background-color: var(--color-background-muted);
  border: 1px solid var(--color-primary-dark);
  color: var(--color-primary-dark);
}
  
button.danger {
  background-color: var(--color-error);
}

button.danger:hover {
  background-color: var(--color-error-dark);
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}
</style>
