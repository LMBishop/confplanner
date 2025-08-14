<script setup lang="ts">
import { TrainTrack } from 'lucide-vue-next';
import { useScheduleStore } from '~/stores/schedule';

const route = useRoute();
const scheduleStore = useScheduleStore();

const track = scheduleStore.schedule?.tracks.find((track) => track.slug === route.params.slug);
</script>

<template>
  <Panel v-if="track" :title="track.name" :breadcrumbs="[{ text: 'Tracks', to: '/tracks' }]" :icon="TrainTrack">
    <ul class="events-list">
      <li 
        v-for="event in scheduleStore.eventsPerTrack[track.name]" 
        :key="event.id" 
        class="event-item"
      >
        <EventListing :event="event" />
      </li>
    </ul>
  </Panel>
</template>

<style>
.events-list {
  list-style: none;
  margin: -1rem 0;
  padding: 0;
  display: grid;
}

.event-item {
  border-bottom: 1px solid var(--color-background-muted); 
}

.events-title {
  margin-bottom: 1rem;
}

.event-item:last-child {
  border-bottom: none; 
}

</style>