import { useLoginOptionsStore } from "~/stores/login-options";

export default function() {
    const loginOptionsStore = useLoginOptionsStore();
    const errorStore = useErrorStore();
    const config = useRuntimeConfig();
    
    loginOptionsStore.status = 'pending'
    
    $fetch(config.public.baseURL + '/login', {
      method: 'GET',
      server: false,
      lazy: true,
      onResponse: ({ response }) => {
        if (!response.ok) {
          errorStore.setError(response._data.message || 'An unknown error occurred');
        }

        if (response._data) {
          loginOptionsStore.loginOptions = (response._data as any).data.options;
          loginOptionsStore.status = 'idle'
        }
      },
    });
}