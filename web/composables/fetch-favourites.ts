export default function() {
  const conferenceStore = useConferenceStore()
  const favouritesStore = useFavouritesStore();
  const errorStore = useErrorStore();
  const config = useRuntimeConfig();
  
  favouritesStore.status = 'pending'

  return $api(config.public.baseURL + '/favourites/' + conferenceStore.id, {
    method: 'GET',
    onResponse: ({ response }) => {
      favouritesStore.status = 'idle'
      if (!response.ok) {
        errorStore.setError(response._data.message || 'An unknown error occurred');
      }
      if (response._data) {
        favouritesStore.setFavourites((response._data as any).data);
      }
    },
  });
}