<script setup lang="ts">
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
  <Panel kind="success">
    <span class="text-icon">
      <Spinner />
      <span>Successfully logged in. Navigating to {{ destination }}...</span>
    </span>
  </Panel>
</template>
