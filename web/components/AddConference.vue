<script setup lang="ts">
import { format } from 'date-fns';
import { type Event as ScheduledEvent } from '~/stores/schedule';

const errorStore = useErrorStore();
const config = useRuntimeConfig();
const loading = ref(false)

const emit = defineEmits(['update']);

const addConference = async (e: Event) => {
  const target = e.target as HTMLFormElement;
  const formData = new FormData(target);
  loading.value = true

  $api(config.public.baseURL + '/conference', {
    method: 'POST',
    body: JSON.stringify(Object.fromEntries(formData)),
    onResponse: ({ response }) => {
      loading.value = false
      if (!response.ok) {
        errorStore.setError(response._data?.message || 'An unknown error occurred');
        return
      }
      emit('update')
    },
  });
}

</script>

<template>
  <div>
    <form @submit.prevent="(e) => addConference(e)">
      <div class="form-group">
        <label for="url" class="form-label">
          Schedule data URL
        </label>
        <div class="form-input-container">
          <Input id="url" name="url" required />
        </div>
      </div>

      <div class="form-submit">
        <Button type="submit" :loading="loading">
          Add
        </Button>
      </div>
    </form>
  </div>
</template>

<style scoped>
</style>