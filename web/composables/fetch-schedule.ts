export default function useFetchFavourites() {
    const scheduleStore = useScheduleStore();
    const errorStore = useErrorStore();
    const config = useRuntimeConfig();
    
    useFetch(config.public.baseURL + '/schedule', {
      method: 'GET',
      server: false,
      lazy: true,
      onResponse: ({ response }) => {
        if (!response.ok) {
          if (response.status === 401) {
            navigateTo({ name: 'login', state: { error: 'Sorry, your session has expired' } });
          } else {
            errorStore.setError(response._data.message || 'An unknown error occurred');
          }
        }

        if (response._data) {
          scheduleStore.setSchedule((response._data as any).data.schedule);
          errorStore.setError("Schedule set");
        } else {
          errorStore.setError("Invalid response returned by server");
        }
      },
    });
}