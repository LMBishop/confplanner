<script setup lang="ts">
import { TrainTrack } from 'lucide-vue-next';
import Panel from '~/components/Panel.vue';

definePageMeta({
  middleware: ['logged-in', 'conference-selected']
})
const scheduleStore = useScheduleStore();
</script>

<template>
  <div v-if="scheduleStore.status === 'pending'" class="loading">
    <span class="loading-text">
      <Spinner color="var(--color-text-muted)" />Updating schedule...
    </span>
  </div>
  <Panel v-else title="Tracks" :icon="TrainTrack">
    <ul class="tracks-list">
      <li 
        v-for="track in scheduleStore.schedule?.tracks" 
        :key="track.name" 
        class="tracks-item"
      >
        <NuxtLink :to="'/tracks/' +  track.slug" class="track-item">
          {{ track.name }}
        </NuxtLink>
      </li>
    </ul>
  </Panel>    
</template>

<style scoped>
.tracks-list {
  list-style: none;
  margin: -0.5rem 0 0 0;
  padding: 0;
  display: grid;
}
          
.track-item {
  position: relative;
  border-bottom: 1px solid var(--color-background-muted); 
  padding: 0.5rem 1rem;
  left: -1rem;
  width: calc(100%);
  display: block;
  text-decoration: none;
}

.track-item:last-child {
  border-bottom: none; 
}

.track-item:hover {
  background-color: var(--color-background-muted);
}
</style>