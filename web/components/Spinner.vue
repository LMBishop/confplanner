<script setup lang="ts">
import { Loader2Icon } from 'lucide-vue-next'

const props = defineProps<{
  color?: string
  size?: number
}>()

</script>

<template>
  <Loader2Icon class="icon-loader" :color="color" :size="size" />
</template>

<style scoped>
.icon-loader {
  animation: spin 1s linear infinite;
  color: var(--color-text);
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
