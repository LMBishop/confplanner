import { type Event } from "./schedule";
import { defineStore } from "pinia";

export const useSelectedEventStore = defineStore('selected-event', () => {
  const selectedEvent = ref(null as Event | null)
  const setSelectedEvent = (event: Event) => {
    selectedEvent.value = event
  }
  const clearSelectedEvent = () => {
    selectedEvent.value = null
  }
  
  return {selectedEvent, setSelectedEvent, clearSelectedEvent}
})
