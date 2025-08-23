const authStore = useAuthStore()

export default defineNuxtRouteMiddleware((to, from) => {
  if (!authStore.isLoggedIn()) {
    return navigateTo("/login");
  }
});

