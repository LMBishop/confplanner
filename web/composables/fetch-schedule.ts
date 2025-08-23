import { useConferenceStore } from "~/stores/conference";
import { expireAuth } from "./expire-auth";

export default async function() {
  const conferenceStore = useConferenceStore()
  const scheduleStore = useScheduleStore();
  const errorStore = useErrorStore();
  const config = useRuntimeConfig();

  scheduleStore.status = 'pending'
  
  return $api(config.public.baseURL + '/conference/' + conferenceStore.id, {
    method: 'GET',
    onResponse: ({ response }) => {
      if (!response.ok) {
        if (response.status === 401) {
          expireAuth()
          return
        } else {
          errorStore.setError(response._data.message || 'An unknown error occurred');
        }
      }

      if (response._data) {
        let schedule = (response._data as any).data.schedule
        scheduleStore.setSchedule(schedule);

        conferenceStore.venue = schedule.conference.venue
        conferenceStore.title = schedule.conference.title
        conferenceStore.city = schedule.conference.city

        scheduleStore.status = 'idle'
      }
    },
  }).catch(() => {
    // todo do this better
    errorStore.setError('An unknown error occurred');
  });
}