<script setup lang="ts">
import { cn } from '@/lib/utils'
import { DropdownMenuItem, type DropdownMenuItemProps } from 'radix-vue'
import type { HTMLAttributes } from 'vue'

const props = defineProps<DropdownMenuItemProps & { class?: HTMLAttributes['class'], inset?: boolean }>()
</script>

<template>
  <DropdownMenuItem
    v-bind="props"
    :class="cn(
      'relative flex cursor-default select-none items-center rounded-sm gap-2 px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50  [&>svg]:size-4 [&>svg]:shrink-0',
      inset && 'pl-8',
      props.class,
    )"
  >
    <slot />
  </DropdownMenuItem>
</template>
