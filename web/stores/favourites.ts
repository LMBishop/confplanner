import { defineStore } from "pinia";
import { type Event } from "./schedule";

interface Favourite {
  id: number;
  eventGuid?: string;
  eventId?: number;
}

export const useFavouritesStore = defineStore('favourites', () => {
  const favourites = ref([] as Favourite[])
  const status = ref('idle' as 'idle' | 'pending')

  const setFavourites = (newFavourites: Favourite[]) => {
    favourites.value = newFavourites
  } 
  const addFavourite = (favourite: Favourite) => {
    favourites.value.push(favourite)
  }
  const removeFavourite = (favourite: { eventGuid?: string, eventId?: number }) => {
    if (favourite.eventGuid) {
      favourites.value = favourites.value.filter(f => f.eventGuid !== favourite.eventGuid)
    }
    if (favourite.eventId) {
      favourites.value = favourites.value.filter(f => f.eventId !== favourite.eventId)
    }
  }
  const isFavourite = (event: Event) => {
    return favourites.value.some(f => {
      if (f.eventGuid) {
        return f.eventGuid === event.guid
      } else if (f.eventId) {
        return f.eventId === event.id
      }
    })
  }
  
  const setStatus = (newStatus: 'idle' | 'pending') => {
    status.value = newStatus
  }

  return {favourites, status, setFavourites, addFavourite, removeFavourite, isFavourite, setStatus}
})
