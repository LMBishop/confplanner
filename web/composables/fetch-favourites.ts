export default function() {
    const favouritesStore = useFavouritesStore();
    const errorStore = useErrorStore();
    const config = useRuntimeConfig();
    
    favouritesStore.setStatus('pending')

    useFetch(config.public.baseURL + '/favourites', {
      method: 'GET',
      server: false,
      lazy: true,
      onResponseError: ({ response }) => {
        favouritesStore.setStatus('idle')
        errorStore.setError(response._data.message || 'An unknown error occurred');
      },
      onResponse: ({ response }) => {
        if (response._data) {
          favouritesStore.setFavourites((response._data as any).data);
        }
        favouritesStore.setStatus('idle')
      },
    });
}