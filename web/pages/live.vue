<script setup lang="ts">
import { Radio } from 'lucide-vue-next';
import EventListing from '~/components/EventListing.vue';
import Panel from '~/components/Panel.vue';
import { type Event } from '~/stores/schedule';

const favouritesStore = useFavouritesStore();
const scheduleStore = useScheduleStore();

const showAllHappeningNow = ref(false);
const showAllUpcoming = ref(false);
const timer = ref();
const refreshKey = ref(0);

const todayEvents = computed(() => {
  return scheduleStore.events.filter((event) => isToday(new Date(event.start)));
});
const happeningNow = computed(() => {
  refreshKey.value;
  return todayEvents.value.filter((event) => isEventHappeningNow(event));
});
const favouritesHappeningNow = computed(() => {
  refreshKey.value;
  return happeningNow.value.filter((event) => favouritesStore.isFavourite(event));
});
const upcomingToday = computed(() => {
  refreshKey.value;
  return todayEvents.value.filter((event) => !isEventHappeningNow(event) && isInFuture(event));
});
const favouritesUpcomingToday = computed(() => {
  refreshKey.value;
  return upcomingToday.value.filter((event) => favouritesStore.isFavourite(event));
});

onMounted(() => {
  timer.value = setInterval(() => {
    refreshKey.value++;
  }, 1000);
});

function isEventHappeningNow(event: Event): boolean {
  const now = new Date();
  return event.start <= now && event.end >= now;
}

function isInFuture(event: Event): boolean {
  return event.start > new Date();
}

function isToday(date: Date): boolean {
  const today = new Date();
  return date.getDate() === today.getDate() &&
         date.getMonth() === today.getMonth() &&
         date.getFullYear() === today.getFullYear();
}
</script>

<template>
  <Panel kind="emphasis" class="ongoing" v-if="happeningNow.length > 0" title="Now" :icon="Radio">
    <ul class="events-list">
      <li v-for="event in showAllHappeningNow ? happeningNow : favouritesHappeningNow" :key="event.id">
        <EventListing :event="event" />
      </li>
      <span class="other" v-if="!showAllHappeningNow && happeningNow.length != favouritesHappeningNow.length">{{ happeningNow.length }} other events happening right now</span>
    </ul> 
    
    <template #actions>
      <div class="actions">
        <span class="status" v-if="!showAllHappeningNow">Showing only agenda items</span>
        <Button @click="showAllHappeningNow = !showAllHappeningNow" :kind="showAllHappeningNow ? 'secondary' : 'primary'">
          {{ showAllHappeningNow ? 'Show only agenda items' : 'Show all' }}
        </Button>
      </div>
    </template> 
  </Panel>
  
  <Panel class="upcoming" v-if="upcomingToday.length > 0" title="Upcoming">
    <ul class="events-list">
      <li v-for="event in showAllUpcoming ? upcomingToday : favouritesUpcomingToday" :key="event.id">
        <EventListing :event="event" :show-relative-time="true" />
      </li>
      <span class="other" v-if="!showAllUpcoming">{{ upcomingToday.length }} other upcoming events today</span>
    </ul>

    <template #actions>
      <div class="actions">
        <span class="status" v-if="!showAllUpcoming">Showing only agenda items</span>
        <Button @click="showAllUpcoming = !showAllUpcoming" :kind="showAllUpcoming ? 'secondary' : 'primary'">
          {{ showAllUpcoming ? 'Show only agenda items' : 'Show all' }}
        </Button>
      </div>
    </template> 
  </Panel> 
  
  <Panel kind="error" v-if="todayEvents.length === 0">
    <span>There are no more events today</span>
  </Panel>

</template>

<style scoped>
ul.events-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  margin-top: -1rem;
}

ul.events-list > .other {
  font-style: italic;
}

ul.events-list > li {
  border-bottom: 1px solid var(--color-hover); 
}

ul.events-list > li:last-of-type {
  border-bottom: none;
}

div.actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

div.actions > span.status {
  color: var(--color-text);
  font-style: italic;
  font-size: var(--text-small);
}

@media (max-width: 800px) {
  div.actions > span.status {
    display: none;
  }  
}
</style>