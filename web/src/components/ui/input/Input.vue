<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { useAttrs } from 'vue'

interface Props {
  id?: string
  modelValue?: string
  class?: HTMLAttributes['class']
  type?: string
  disabled?: boolean
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
})

const attrs = useAttrs()
</script>

<template>
  <input
    :id="id"
    :type="type"
    :value="modelValue"
    :disabled="disabled"
    :placeholder="placeholder"
    :class="cn(
      'flex h-9 w-full rounded-sm border border-default bg-canvas px-3 py-1 text-[14px] text-ink transition-colors duration-[var(--duration-fast)] file:border-0 file:bg-transparent file:text-[14px] file:font-medium placeholder:text-[color:var(--color-ink-subtle)] hover:border-strong focus-visible:outline-none focus-visible:border-focus focus-visible:shadow-[var(--shadow-focus-ring)] disabled:cursor-not-allowed disabled:bg-sunken disabled:text-[color:var(--color-ink-disabled)]',
      props.class,
    )"
    v-bind="attrs"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>
