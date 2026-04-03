<script setup lang="ts">
import { cn } from '@/lib/utils'
import { ChevronRight } from 'lucide-vue-next'
import { PaginationNext, type PaginationNextProps } from 'radix-vue'
import { Button } from '@/components/ui/button'
import type { HTMLAttributes } from 'vue'

const props = withDefaults(
  defineProps<PaginationNextProps & { class?: HTMLAttributes['class'] }>(),
  {
    asChild: true,
  },
)
</script>

<template>
  <PaginationNext v-bind="props">
    <Button :class="cn('w-9 h-9 p-0', props.class)" variant="outline">
      <slot>
        <ChevronRight />
      </slot>
    </Button>
  </PaginationNext>
</template>
