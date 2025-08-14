<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next';
import { defineProps, type FunctionalComponent, type PropType } from 'vue';

defineProps({
  kind: {
    type: String as PropType<"normal" | "error" | "success" | "emphasis">,
    required: false,
    default: 'normal',
  },
  title: {
    type: String,
    required: false,
    default: '',
  },
  breadcrumbs: {
    type: Array as PropType<Array<{ text: string, to: string }>>,
    required: false,
    default: () => [],
  },
  icon: {
    type: Object as PropType<FunctionalComponent>,
    default: null
  }
});
</script>
<template>
  <div class="card" :class="kind">
    <div v-if="title" class="card-header">
      <div class="header-left">
        <component :is="icon" v-if="icon" />
        <div class="card-title-container">
          <span v-for="(breadcrumb, index) in breadcrumbs" class="breadcrumb" :key="index">
            <NuxtLink :to="breadcrumb.to">{{ breadcrumb.text }}</NuxtLink>
            <ArrowRight v-if="index != breadcrumbs.length - 1" />
          </span>
          <span class="card-title"> {{ title }} </span>
        </div>
      </div>
      
      <slot name="actions" />
    </div>
    <slot />
  </div>
</template>

<style>
div.card {
  padding: 1rem; 
  border-radius: 0.5rem; 
}

div.card-header {
  padding: 1rem 1rem;
  margin: -1rem -1rem 0 -1rem;
  border-top-left-radius: 0.5rem;
  border-top-right-radius: 0.5rem;;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

div.header-left {
  display: flex;
  line-height: 1;
  align-items: center;
  gap: 0.5rem;
}

div.header-left > svg {
  width: 1.3rem;
  height: 1.3rem;
}

span.card-title {
  font-size: 1.5rem;
  font-weight: 700;
}

span.breadcrumb {
  font-size: 1.2rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

span.breadcrumb > svg {
  width: 1.2rem;
  height: 1.2rem;
  color: var(--color-text-muted);
}

span.breadcrumb > a {
  color: var(--color-text-muted);
}

span.breadcrumb > a:hover {
  color: var(--color-primary);
}

div.normal {
  background-color: white;
  border: 0.1rem solid var(--color-border);
}

div.normal .card-header {
  background-color: var(--color-background-muted);
  border-bottom: 1px solid var(--color-border);
  margin: -1rem -1rem 1rem -1rem;
}

div.error {
  background-color: var(--color-error-light);
  border: 0.1rem solid var(--color-border-error-light);
}

div.success {
  background-color: var(--color-success-light);
  border: 0.1rem solid var(--color-border-success-light);
}

div.emphasis {
  background-color: var(--color-background-primary);
  border: 0.1rem solid var(--color-border-primary-light);
}
</style>