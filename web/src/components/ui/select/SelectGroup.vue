<script setup lang="ts">
import { cn } from '@/lib/utils'
import { SelectGroup, type SelectGroupProps } from 'radix-vue'
import { type HTMLAttributes } from 'vue'

const props = defineProps<SelectGroupProps & { class?: HTMLAttributes['class'] }>()
</script>

<template>
  <SelectGroup :class="cn('p-1 w-full', props.class)" v-bind="props">
    <slot />
  </SelectGroup>
</template>
