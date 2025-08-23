<script setup lang="ts">
import { ref } from 'vue'
import { FetchError } from 'ofetch'
import Input from '~/components/Input.vue'
import { useLoginOptionsStore } from '~/stores/login-options'

definePageMeta({
  layout: 'none'
})

const authenticating = ref(false)
const authenticatingProvider = ref('')
const error = ref("")
const completingJourney = ref(false)
const basicAuthEnabled = ref(false)

const route = useRoute()
const config = useRuntimeConfig()
const authStore = useAuthStore()
const loginOptionsStore = useLoginOptionsStore()
const headers = useRequestHeaders(['cookie'])

const { loginOptions, status } = storeToRefs(loginOptionsStore)

watch(loginOptions, (options) => {
  basicAuthEnabled.value = options.some(o => o.type === 'basic')
})

const authFail = (e: any) => {
  if ((e as FetchError).data) {
    error.value = e.data.message
  } else {
    error.value = "An unknown error occurred"
  }

  authenticating.value = false
  authenticatingProvider.value = ''
}

const handleBasicAuth = async (e: Event, providerName: string) => {
  const target = e.target as HTMLFormElement;
  const formData = new FormData(target);

  authenticating.value = true
  authenticatingProvider.value = providerName
  
  $fetch(config.public.baseURL + '/login/' + providerName, {
    method: 'POST',
    body: JSON.stringify(Object.fromEntries(formData)),
    headers: headers,
    server: false,
    onResponse: ({ response }) => {
      authStore.token = response._data.data.token
      authStore.username = response._data.data.username
      authStore.admin = response._data.data.admin

      navigateTo("/");
    },
    onResponseError: authFail
  });
  
}

const handleOIDCAuth = async (providerName: string) => {
  authenticating.value = true
  authenticatingProvider.value = providerName
  
  $fetch(config.public.baseURL + '/login/' + providerName, {
    method: 'POST',
    headers: headers,
    server: false,
    onResponse: ({ response }) => {
      navigateTo(response._data.data.url, { external: true })
    },
    onResponseError: authFail
  });
}

onMounted(async () => {
  if (history.state.error) {
    error.value = history.state.error as string
  }

  if (route.params.provider) {
    completingJourney.value = true

    try {
      let state = route.query.state
      let code = route.query.code


      let response: any = await $fetch(config.public.baseURL + '/login/' + route.params.provider, {
        method: 'POST',
        headers: headers,
        server: false,
        body: {
          state: state,
          code: code,
        }
      });

      if (response.code === 307) {
        throw Error()
      }

      authStore.token = response.data.token
      authStore.username = response.data.username
      authStore.admin = response.data.admin

      navigateTo("/");
    } catch (e: any) {
      if ((e as FetchError).data) {
        error.value = e.data.message
      } else {
        error.value = "An unknown error occurred"
      }

      completingJourney.value = false
      fetchLogin()
    }
    return
  }

  fetchLogin()
})

</script>

<template>
  <div class="auth-container">
    <div class="auth-header">
      <h2 class="auth-title">Sign in</h2>

      <div v-if="error" class="auth-error">
        {{ error }}
      </div>
    </div>

    <div class="auth-body">
      <Panel>
        <div class="auth-form">
          <div v-if="completingJourney" class="spinner">
            <Spinner color="var(--color-text-muted)" />Completing login...
          </div>
          <div v-if="status === 'pending'" class="spinner">
            <Spinner color="var(--color-text-muted)" />Getting login options...
          </div>
          <div v-for="option in loginOptions">
            <form v-if="option.type === 'basic'" class="basic-form" @submit.prevent="(e) => handleBasicAuth(e, option.identifier)">
              <div class="form-group">
                <label for="username" class="form-label">
                  Username
                </label>
                <div class="form-input-container">
                  <Input id="username" name="username" required />
                </div>
              </div>

              <div class="form-group">
                <label for="password" class="form-label">
                  Password
                </label>
                <div class="form-input-container">
                  <Input id="password" name="password" type="password" autocomplete="current-password" required />
                </div>
              </div>


              <div class="form-submit">
                <Button type="submit" :loading="authenticatingProvider === option.identifier" :disabled="authenticating">
                  Sign in
                </Button>
              </div>
            </form>

            <div v-if="option.type === 'oidc'" class="auth-provider">
              <Button type="button" :loading="authenticatingProvider === option.identifier" :disabled="authenticating" @click="(e) => handleOIDCAuth(option.identifier)">
                Sign in with {{ option.name }}
              </Button>
            </div>
          </div>

          <Version class="version" />
        </div>
      </Panel>

    </div>

    <div v-if="basicAuthEnabled" class="form-footer">
      <NuxtLink to="/register" class="register-link">
        Register
      </NuxtLink>
    </div>

  </div>
</template>

<style scoped>
div.auth-container {
  min-height: 100vh;
  background-color: var(--color-background-muted);
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1rem;
}

div.auth-header {
  margin: 0 auto;
  width: 100%;
  max-width: 28rem;
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-direction: column;
}

h2.auth-title {
  margin-top: 1.5rem;
  font-size: 1.875rem;
  font-weight: 800;
  color: #1f2937; 
}

div.auth-body {
  margin-top: 2rem; 
  margin: 0 auto;
  width: 100%;
  max-width: 28rem; 
}

div.auth-form {
  display: grid;
  gap: 1.5rem; 
}

div.auth-error {
  color: var(--color-text-error);
  font-style: italic;
}

div.form-footer {
  display: flex;
  justify-content: flex-end;
  margin: 0 auto;
  max-width: 28rem;
}

div.form-submit {
  display: flex;
}

div.form-submit button {
  width: 100%;
}
    
.version {
  font-size: var(--text-smaller);
  margin: 0 auto;
  color: var(--color-text-muted-light);
}

.auth-provider button {
  display: flex;
  width: 100%;
}

.register-link {
  font-size: var(--text-small); 
  font-weight: 500; 
}

input[name="username"] {
  text-transform: lowercase;
}

.spinner {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  font-size: var(--text-normal);
  color: var(--color-text-muted);
}
</style>
