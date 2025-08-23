import { useLocalStorage } from "@vueuse/core";
import { defineStore } from "pinia";

export const useConferenceStore = defineStore('conference', () => {
  const id = useLocalStorage('conference/id', null)
  const title = useLocalStorage('conference/title', null)
  const venue = useLocalStorage('conference/venue', null)
  const city = useLocalStorage('conference/city', null)

  const clear = () => {
    id.value = null
    title.value = null
    venue.value = null
    city.value = null
  }
  
  return {id, title, venue, city, clear}
})
