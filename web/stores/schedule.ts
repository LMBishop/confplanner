import { TZDate } from "@date-fns/tz";
import { addDays } from "date-fns";
import { defineStore } from "pinia";

interface Schedule {
  conference: Conference;
  tracks: Track[];
  days: Day[];
}

interface Conference {
  title: string;
  venue: string;
  city: string;
  start: Date;
  end: Date;
  days: number;
  dayChange: string;
  timeslotDuration: string;
  baseUrl: string;
  timeZoneName: string;
}

interface Track {
  name: string;
  slug: string;
}

interface Day {
  date: string;
  start: Date;
  end: Date;
  rooms: Room[];
}

interface Room {
  name: string;
  events: Event[];
}

export interface Event {
  id: number;
  guid: string;
  date: string;
  start: Date;
  end: Date;
  duration: string;
  room: string;
  url: string;
  track: Track;
  type: string;
  title: string;
  abstract: string;
  persons: Person[];
  attachments: Attachment[];
  links: Link[];
}

interface Person {
  id: number;
  name: string;
}

interface Attachment {
  type: string;
  href: string;
  name: string;
}

interface Link {
  href: string;
  name: string;
}

export const useScheduleStore = defineStore('schedule', () => {
  const schedule = ref(null as Schedule | null)
  const status = ref('idle' as 'idle' | 'pending')

  const events = ref([] as Event[])
  const eventsPerDay = ref({} as { [key: string]: Event[] })
  const eventsPerTrack = ref({} as { [key: string]: Event[] })
  
  const tracks = ref({} as { [key: string]: Track })
  
  const setSchedule = (newSchedule: Schedule) => {
    tracks.value = {}
    newSchedule.tracks.forEach(track => {
      track.slug = convertToSlug(track.name)
      tracks.value[track.name] = track
    });

    newSchedule.conference.start = new TZDate(newSchedule.conference.start, newSchedule.conference.timeZoneName)
    newSchedule.conference.end = addDays(new TZDate(newSchedule.conference.end, newSchedule.conference.timeZoneName), 1)
    
    events.value = []
    newSchedule.days.forEach(day => {
      day.start = new TZDate(day.start, newSchedule.conference.timeZoneName)
      day.end = new TZDate(day.end, newSchedule.conference.timeZoneName)
      day.rooms.forEach(room => {
        room.events.forEach(event => {
          normalizeDates(event, newSchedule.conference.timeZoneName)

          events.value.push(event)

          if (event.track) {
            event.track = tracks.value[event.track as unknown as string]
          }
        })
      })
    })
    events.value.sort((a, b) => { 
      return a.start.getTime() - b.start.getTime()
    })

    schedule.value = newSchedule

    eventsPerDay.value = {}
    events.value.forEach(event => {
      const date = event.date.split('T')[0]
      if (!eventsPerDay.value[date]) {
        eventsPerDay.value[date] = []
      }
      eventsPerDay.value[date].push(event)
    })
    
    eventsPerTrack.value = {}
    events.value.forEach(event => {
      if (!event.track) return
      if (!eventsPerTrack.value[event.track.name]) {
        eventsPerTrack.value[event.track.name] = []
      }
      eventsPerTrack.value[event.track.name].push(event)
    });
  }
  
  const isConferenceOngoing = () => {
    if (!schedule.value) {
      return false
    }
    return new Date() >= schedule.value.conference.start && new Date() < schedule.value.conference.end
  }
  
  const isConferenceFinished = () => {
    if (!schedule.value) {
      return false
    }
    return new Date() > schedule.value.conference.end
  }
  
  const getStartDate = () => {
    return schedule.value?.conference.start || 0
  }
  
  return {schedule, events, eventsPerDay, eventsPerTrack, status, setSchedule, isConferenceOngoing, isConferenceFinished, getStartDate}
})

function normalizeDates(event: Event, timeZone: string) {
  event.start = new TZDate(event.start, timeZone)
  event.end = new TZDate(event.end, timeZone)
}

function parseDuration(duration: string) {
  const [hours, minutes] = duration.split(':').map(Number)
  return hours * 60 * 60 * 1000 + minutes * 60 * 1000
}

function convertToSlug(text: string) {
  return text
    .toLowerCase()
    .replace(/ /g, '-')
    .replace(/[^\w-]+/g, '')
}