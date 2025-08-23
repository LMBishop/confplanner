import { defineStore } from "pinia";

interface LoginOption {
  name: string;
  identifier: string;
  type: string;
}

export const useLoginOptionsStore = defineStore('loginOptions', () => {
  const loginOptions = ref([] as LoginOption[])
  const status = ref('idle' as 'idle' | 'pending')

  return {loginOptions, status}
})
