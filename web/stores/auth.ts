import { useLocalStorage } from "@vueuse/core";
import { defineStore } from "pinia";

export const useAuthStore = defineStore('auth', () => {
  const token = useLocalStorage('auth/token', null)
  const username = useLocalStorage('auth/username', null)
  const admin = useLocalStorage('auth/admin', false)

  const isLoggedIn = () => token.value != null

  return {token, username, admin, isLoggedIn}
})
