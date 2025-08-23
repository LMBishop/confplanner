export function $api<T>(
  request: Parameters<typeof $fetch<T>>[0],
  opts?: Parameters<typeof $fetch<T>>[1],
) {
  const authStore = useAuthStore()

  return $fetch<T>(request, {
    ...opts,
    headers: {
      Authorization: authStore.isLoggedIn() ? `Bearer ${authStore.token}` : '',
      ...opts?.headers,
    },
  })
}
