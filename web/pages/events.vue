<script setup lang="ts">
import { Calendar, SquareGanttChart } from 'lucide-vue-next';
import { useScheduleStore } from '~/stores/schedule';

const scheduleStore = useScheduleStore();

</script>

<template>
  <Panel title="Events" :icon="SquareGanttChart" v-if="scheduleStore.schedule">
    <div v-for="[day, events] of Object.entries(scheduleStore.eventsPerDay)" :key="day" class="events-container">
      <ul class="events-list">
        <li v-for="event in events" :key="event.id" class="event-item" :data-index="event.id">
          <EventListing :event="event" />
        </li>
      </ul>
    </div>
  </Panel>
</template>

<style>
.events-container {
  margin: -1rem 0;
}

.events-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
}

.event-item {
  border-bottom: 1px solid var(--color-background-muted);
}

.event-item:last-child {
  border-bottom: none;
}
</style>