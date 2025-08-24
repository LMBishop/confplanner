<script setup lang="ts">
import { MapPinPlus, MapPin } from 'lucide-vue-next';
import { expireAuth } from '~/composables/expire-auth';

definePageMeta({
  middleware: ['logged-in']
})

interface Conference {
  id: number,
  url: string,
  title: string,
  venue: string,
  city: string,
}

const setting = ref(false)
const conferences = ref([] as Conference[])

const conferenceStore = useConferenceStore();
const errorStore = useErrorStore();
const authStore = useAuthStore();
const config = useRuntimeConfig();

const status = ref('idle' as 'loading' | 'idle')

const refConfirmDeleteDialog = ref<typeof Dialog>();

const deleteConferenceId = ref(null as number | null)
const deleteConferenceName = ref(null as string | null)
const deleteAction = ref(false);

const fetchConferences = () => {
  status.value = 'loading'
  $api(config.public.baseURL + '/conference', {
    method: 'GET',
    onResponse: ({ response }) => {
      if (!response.ok) {
        if (response.status === 401) {
          expireAuth()
          return
        }
        errorStore.setError(response._data.message || 'An unknown error occurred');
      }
      status.value = 'idle'
      conferences.value = response._data.data
    },
  })
}

const deleteConference = () => {
  $api(config.public.baseURL + '/conference', {
    method: 'DELETE',
    body: {
      id: deleteConferenceId.value
    },
    onResponse: ({ response }) => {
      if (!response.ok) {
        errorStore.setError(response._data.message || 'An unknown error occurred');
      }
      fetchConferences()
    },
  })
}

const selectConference = async (c: any) => {
  setting.value = true
  conferenceStore.id = c.id
  try {
    await fetchSchedule()
    await fetchFavourites()
    navigateTo({ path: "/events" })
  } catch (e) {
    conferenceStore.clear()
    setting.value = false
  }
}

onMounted(fetchConferences)
</script>

<template>
  <template v-if="!setting">
    <Panel title="Conferences" :icon="MapPin">
      <span class="loading-text" v-if="status === 'loading'"><Spinner color="var(--color-text-muted)" />Fetching conferences...</span>
      <div class="conference-list" v-if="conferences?.length > 0 && status !== 'loading'">
        <template v-for="conference of conferences">
          <span class="title">{{ conference.title }}</span>
          <span>{{ conference.city }}</span>
          <span>{{ conference.venue }}</span>
          <span class="actions">
            <Button v-if="authStore.admin" kind="secondary" @click="() => { 
              deleteConferenceId = conference.id
              deleteConferenceName = conference.title
              refConfirmDeleteDialog!.show()
            }">Delete</Button>
            <Button @click="() => { selectConference(conference) }">Select</Button>
          </span>
        </template>
      </div>
      <p v-if="(conferences == null || conferences?.length == 0) && status !== 'loading'">
        There are no conferences to display.
      </p>
    </Panel>
    
    <Panel v-if="authStore.admin" title="Add conference" :icon="MapPinPlus">
      <AddConference @update="fetchConferences" />
    </Panel>
  </template>
  <template v-else>
    <div class="loading">
      <span class="loading-text"><Spinner color="var(--color-text-muted)" />Setting conference...</span>
    </div>
  </template>

  <Dialog ref="refConfirmDeleteDialog" title="Delete conference" :confirmation="true" @submit="deleteConference" :fit-contents="true">
    <span>Are you sure you want to delete "{{ deleteConferenceName }}"?</span>
    <span>All favourites for all users for this conference will be lost and irrecoverable.</span>
    <template v-slot:actions>
      <Button kind="secondary" type="button" @click="refConfirmDeleteDialog!.close()">Cancel</Button>
      <Button kind="danger" type="submit" :loading="deleteAction">Delete</Button>
    </template>
  </Dialog>
</template>

<style>
.conference-list {
  display: grid;
  grid-template-columns: 1fr 1fr 2fr 1fr;
  align-items: center;
  gap: 0.5rem;
}

.conference-list > .title {
  font-weight: bold;
}

.conference-list > .actions {
  display: flex;
  gap: 0.5rem;
  justify-self: flex-end;
}
</style>