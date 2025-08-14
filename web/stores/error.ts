import { defineStore } from "pinia";

export const useErrorStore = defineStore('error', () => {
  const error = ref(null as string | null)
  const setError = (newError: string) => {
    error.value = newError
  }
  const clearError = () => {
    error.value = null
  }
  
  return {error, setError, clearError}
})
