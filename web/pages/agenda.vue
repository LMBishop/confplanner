<script setup lang="ts">
import Panel from '~/components/Panel.vue';
import Dialog from '~/components/Dialog.vue';
import { Calendar } from 'lucide-vue-next';

const favouritesStore = useFavouritesStore();
const scheduleStore = useScheduleStore();
const errorStore = useErrorStore();
const config = useRuntimeConfig();

const favouriteEvents = computed(() => {
  return scheduleStore.events.filter((event) => favouritesStore.isFavourite(event));
});

const calendarStatus = ref('pending' as 'pending' | 'idle');
const calendarLink = ref('')
const calendarLinkWithPageProtocol = computed(() => {
  return window.location.protocol + '//' + calendarLink.value;
});

const refConfirmDeleteDialog = ref<typeof Dialog>();

const calendarAction = ref(false);

useFetch(config.public.baseURL + '/calendar', {
  method: 'GET',
  server: false,
  lazy: true,
  onResponse: ({ response }) => {
    if (!response.ok) {
      if (response.status !== 404) {
        errorStore.setError(response._data.message || 'An unknown error occurred');
      }
    } else if (response._data) {
      calendarLink.value = (response._data as any).data.url;
    }
    calendarStatus.value = 'idle';
  },
});

function generateCalendar() {
  calendarAction.value = true;
  useFetch(config.public.baseURL + '/calendar', {
    method: 'POST',
    server: false,
    lazy: true,
    onResponseError: ({ response }) => {
      errorStore.setError(response._data.message || 'An unknown error occurred');
      calendarAction.value = false;
    },
    onResponse: ({ response }) => {
      if (response._data) {
        calendarLink.value = (response._data as any).data.url;
      }
      calendarAction.value = false;
    },
  });
}

function deleteCalendar() {
  calendarAction.value = true;
  useFetch(config.public.baseURL + '/calendar', {
    method: 'DELETE',
    server: false,
    onResponseError: ({ response }) => {
      errorStore.setError(response._data.message || 'An unknown error occurred');
      calendarAction.value = false;
    },
    onResponse: () => {
      calendarLink.value = '';
      calendarAction.value = false;
    },
  });
}

</script>

<template>
  <Panel v-if="favouritesStore.status === 'pending'">
    <span>Updating favourites...</span>
  </Panel>

  <template v-else-if="favouriteEvents.length > 0" >
    <div class="page">
      <Panel title="Agenda" :icon="Calendar">
        <ul class="agenda-list">
          <li v-for="event in favouriteEvents" :key="event.id" class="agenda-item" >
            <EventListing :event="event" />
          </li>
        </ul>
      </Panel>
      <Panel>
        <template v-if="calendarStatus === 'pending'">
          <span>Fetching calendar status...</span>
        </template>
        
        <div v-else-if="calendarStatus === 'idle'" class="calendar">
          <template v-if="calendarLink">
            <span>You can add your agenda to your own calendar using the iCal link below</span>
            <Input :value="calendarLinkWithPageProtocol" readonly/>
            <Button @click="refConfirmDeleteDialog!.show()" :loading="calendarAction">Delete calendar</Button>
          </template>
          <template v-else>
            <span>You do not have a calendar link yet. Use the button below to request a calendar link to subscribe to on your own calendar app.</span>
            <Button @click="generateCalendar" :loading="calendarAction">Request calendar</Button>
          </template>
        </div>
      </Panel>
    </div>
  </template>

  <Panel v-else>
    <span>You have not added any favourites yet.</span>
  </Panel>
  
  <Dialog ref="refConfirmDeleteDialog" title="Delete calendar" :confirmation="true" @submit="deleteCalendar" :fit-contents="true">
    <span>Are you sure you want to delete your calendar?</span>
    <span>Your unique link cannot be recovered if you do so.</span>
    <template v-slot:actions>
      <Button kind="secondary" type="button" @click="refConfirmDeleteDialog!.close()">Cancel</Button>
      <Button kind="danger" type="submit" :loading="calendarAction">Delete</Button>
    </template>
  </Dialog>
</template>

<style scoped>
.agenda-list {
  list-style: none;
  margin: -1rem 0;
  padding: 0;
  display: grid;
}

.agenda-item {
  border-bottom: 1px solid var(--color-background-muted); 
}

.agenda-title, .calendar-title {
  margin-bottom: 1rem;
}

.agenda-item:last-child {
  border-bottom: none; 
}

.page, .calendar {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  line-height: 1.3;
}
</style>