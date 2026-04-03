<script setup lang="ts">
import { cn } from '@/lib/utils'
import { PaginationListItem, type PaginationListItemProps } from 'radix-vue'
import { Button } from '@/components/ui/button'
import type { HTMLAttributes } from 'vue'

const props = defineProps<PaginationListItemProps & { class?: HTMLAttributes['class'], isSelected: boolean }>()
</script>

<template>
  <PaginationListItem v-bind="props">
    <Button
      :class="cn('w-9 h-9 p-0', props.isSelected ? 'bg-primary text-primary-foreground' : '', props.class)"
      variant="outline"
    >
      <slot />
    </Button>
  </PaginationListItem>
</template>
