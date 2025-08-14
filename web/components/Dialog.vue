<script setup lang="ts">
import { ref } from 'vue';

const visible = ref(false);
const refDialog = ref<HTMLDialogElement | null>(null);

const props = defineProps<{
  kind?: 'normal' | 'error';
  fitContents?: boolean;
  title?: string;
}>();

const showModal = () => {
  refDialog.value?.showModal();
  visible.value = true;
};

const closeModal = () => {
  refDialog.value?.close();
};

const emit = defineEmits(['close', 'submit']);

defineExpose({
  show: showModal,
  close: closeModal,
  visible,
});

const onClose = () => {
  visible.value = false;
  emit('close');
};

const onSubmit = (e: Event) => {
  e.preventDefault();
  const formData = new FormData(e.target as HTMLFormElement);
  const formValue = Object.fromEntries(formData.entries());
  emit('submit', formValue);
  
  closeModal();
};

const onDivClick = (e: MouseEvent) => {
  e.stopPropagation()
};

const onDialogClick = (e: MouseEvent) => {
  if (e.target === refDialog.value) {
    refDialog.value?.close();
  }
};
</script>

<template>
  <dialog ref="refDialog" @click="onDialogClick" @close="onClose" @submit="onSubmit" :class="[props.kind ?? 'normal', { fit: props.fitContents }]">
    <div @click="onDivClick">
      <form v-if="visible" method="dialog">
        <div class="title" v-if="title">{{ props.title }}</div>

        <slot />
        
        <div class="actions" v-if="$slots.actions">
          <slot name="actions" class="actions" />
        </div>
      </form>
    </div>
  </dialog>
</template>


<style scoped>
dialog {
  outline: none;
  border-radius: 0.5rem;
  padding: 1rem;
  width: 1000px;
  margin: 0;
  top: 80px;
  max-height: calc(100vh - 160px);
  left: 50%;
  transform: translateX(-50%);
}

dialog.normal {
  border: 2px solid var(--color-border);
  background-color: var(--color-background);
}

dialog.error {
  border: 2px solid var(--color-border-error-light);
  background-color: var(--color-error-light);
}

dialog.fit {
  width: fit-content;
  max-width: 1000px;
}

dialog::backdrop {
  backdrop-filter: blur(4px);
  background-color: rgba(0, 0, 0, 0.1);
}

form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

div.title {
  font-size: var(--text-larger);
  font-weight: 700;
}

.actions {
  display: flex;
  gap: 1rem;
  align-self: flex-end;
  justify-content: flex-end;
}
</style>