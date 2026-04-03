<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { useAttrs } from 'vue'

interface Props {
  id?: string
  modelValue?: string
  class?: HTMLAttributes['class']
  disabled?: boolean
  placeholder?: string
  rows?: number
}

const props = withDefaults(defineProps<Props>(), {
  rows: 3,
})

const attrs = useAttrs()
</script>

<template>
  <textarea
    :id="id"
    :value="modelValue"
    :disabled="disabled"
    :placeholder="placeholder"
    :rows="rows"
    :class="cn(
      'flex min-h-[60px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
      props.class,
    )"
    v-bind="attrs"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
