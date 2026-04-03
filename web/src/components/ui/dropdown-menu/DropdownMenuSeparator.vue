<script setup lang="ts">
import { cn } from '@/lib/utils'
import { DropdownMenuSeparator, type DropdownMenuSeparatorProps } from 'radix-vue'
import type { HTMLAttributes } from 'vue'

const props = defineProps<DropdownMenuSeparatorProps & {
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <DropdownMenuSeparator :class="cn('-mx-1 my-1 h-px bg-muted', props.class)" v-bind="props" />
</template>
