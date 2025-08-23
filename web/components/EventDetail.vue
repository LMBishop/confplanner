<script setup lang="ts">
import { format } from 'date-fns';
import { type Event as ScheduledEvent } from '~/stores/schedule';

const { event } = defineProps<{
  event: ScheduledEvent;
}>();

const getHostname = (url: string) => new URL(url).hostname;

</script>

<template>
  <div class="event">
    <div class="event-details">
      <span class="event-info">
        {{ format(event.start, "eeee kk:mm") }} - {{ format(event.end, "kk:mm") }}, {{ event.room }}
      </span>
      <span class="event-title">{{ event.title }}</span>
      <span class="event-speaker">{{ event.persons.map(p => p.name).join(", ") }}</span>
    </div>

    <div class="event-abstract" v-html="event.abstract" />

    <div v-if="event.links.length > 0 || event.attachments.length > 0" class="event-supplementary">
      <div class="event-links" v-if="event.links.length > 0">
        <span class="event-links-header">Links</span>
        <ul>
          <li v-for="link in event.links" :key="link.href" class="event-link">
            <a :href="link.href" target="_blank">{{ link.name }}</a>
          </li>
        </ul>
      </div>

      <div class="event-attachments" v-if="event.attachments.length > 0">
        <span class="event-attachments-header">Attachments</span>
        <ul>
          <li v-for="attachment in event.attachments" :key="attachment.href" class="event-attachment">
            <a :href="attachment.href" target="_blank">{{ attachment.name }}</a>
          </li>
        </ul>
      </div>
    </div>

    <div>
      <span v-if="event.track" class="event-track"><NuxtLink :to="'/tracks/' + event.track.slug">{{ event.track.name }}</NuxtLink> &bull;</span> <a v-if="event.url" class="event-url" :href="event.url" target="_blank">view on {{ getHostname(event.url)}}</a>
    </div>
  </div>
</template>

<style scoped>
.event {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.event-details {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-direction: column;
  line-height: 1.3;
}

.event-info {
  font-size: var(--text-normal); 
  color: var(--color-text); 
  margin: 0;
}

.event-title {
  font-size: var(--text-larger);
  font-weight: 600; 
  margin: 0;
}

.event-supplementary {
  display: flex;
  flex-direction: row;
  gap: 1rem;
}

.event-links-header, .event-attachments-header, .event-track-header {
  font-weight: 600;
}

.event-speaker {
  font-size: var(--text-small);
  color: var(--color-text-muted);
}

.event-track, .event-url {
  font-size: var(--text-small);
}

.event-button {
  background: none;
  border: none;
  color: var(--color-primary); 
  cursor: pointer;
  outline: none;
}

.event-button:hover {
  color: var(--color-primary-dark); 
}

.event-icon {
  height: 1.5rem; 
  width: 1.5rem; 
}
</style>