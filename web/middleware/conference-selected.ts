import { useConferenceStore } from "~/stores/conference";

const conferenceStore = useConferenceStore();
const scheduleStore = useScheduleStore()

export default defineNuxtRouteMiddleware((to, from) => {
  if (conferenceStore.id === null) {
    return navigateTo("/conferences");
  }

  if (scheduleStore.schedule === null) {
    fetchSchedule();
    fetchFavourites();
  }
});