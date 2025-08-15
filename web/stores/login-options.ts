import { defineStore } from "pinia";

interface LoginOption {
  name: string;
  identifier: string;
  type: string;
}

export const useLoginOptionsStore = defineStore('loginOptions', () => {
  const loginOptions = ref([] as LoginOption[])
  const status = ref('idle' as 'idle' | 'pending')

  const setLoginOptions = (newLoginOptions: LoginOption[]) => {
    loginOptions.value = newLoginOptions
  } 

  const setStatus = (newStatus: 'idle' | 'pending') => {
    status.value = newStatus
  }

  return {loginOptions, status, setLoginOptions, setStatus}
})
