<script setup lang="ts">
definePageMeta({
  middleware: ['logged-in', 'conference-selected']
})

const scheduleStore = useScheduleStore();

const destination = ref()

if (scheduleStore.isConferenceOngoing()) {
  destination.value = "/live";
  navigateTo('/live');
} else {
  destination.value = "/events";
  navigateTo('/events');
}
</script>

<template>
  <div v-if="scheduleStore.status === 'pending'" class="loading">
    <span class="loading-text">
      <Spinner color="var(--color-text-muted)" />Updating schedule...
    </span>
  </div>
  <Panel kind="success">
    <span class="text-icon">
      <Spinner />
      <span>Successfully logged in. Navigating to {{ destination }}...</span>
    </span>
  </Panel>
</template>
