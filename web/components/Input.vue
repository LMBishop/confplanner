<script lang="ts">
export default {
  name: 'CustomInput',
  props: {
    id: {
      type: String,
      required: false,
      default: '',
    },
    label: {
      type: String,
      required: false,
      default: '',
    },
    autocomplete: {
      type: String,
      required: false,
      default: '',
    },
    type: {
      type: String,
      required: false,
      default: 'text',
    },
    name: {
      type: String,
      required: false,
      default: '',
    },
    required: {
      type: Boolean,
      required: false,
      default: false,
    },
    modelValue: {
      type: [String, Number],
      required: false,
      default: '',
    },
    placeholder: {
      type: String,
      required: false,
      default: '',
    },
    disabled: {
      type: Boolean,
      required: false,
      default: false,
    },
    readonly: {
      type: Boolean,
      required: false,
      default: false,
    },
    inputClass: {
      type: String,
      required: false,
      default: '',
    },
  },
};
</script>

<template>
  <input
    :id="id"
    :type="type"
    :name="name"
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target! as any).value)"
    :placeholder="placeholder"
    :disabled="disabled"
    :required="required"
    :readonly="readonly"
  />
</template>

<style>
input {
  appearance: none;
  display: block;
  width: 100%;
  padding: 0.5rem 0.75rem; 
  border: 1px solid #d1d5db; 
  border-radius: 2px; 
  box-sizing: border-box;
  /* box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);  */
  font-family: var(--font-family);
  color: var(--color-text); 
  font-size: var(--text-small); 
}

input::placeholder {
  color: #9ca3af; 
}

input:focus {
  outline: none;
  border-color: var(--color-accent); 
  /* box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.4);  */
}
</style>