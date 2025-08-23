export function expireAuth() {
  const authStore = useAuthStore()

  authStore.admin = false;
  authStore.username = null;
  authStore.token = null;
  navigateTo({ path: '/login', state: { error: 'Sorry, your session has expired' } });
}