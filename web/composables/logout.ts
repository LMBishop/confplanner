export function logout() {
  const authStore = useAuthStore()
  const conferenceStore = useConferenceStore()
  const config = useRuntimeConfig();

  $api(config.public.baseURL + '/logout', { method: 'POST' }).finally(() => {
    authStore.admin = false;
    authStore.username = null;
    authStore.token = null;

    conferenceStore.clear()

    navigateTo({ path: '/login', state: { error: 'You have logged out' } });
  })
}