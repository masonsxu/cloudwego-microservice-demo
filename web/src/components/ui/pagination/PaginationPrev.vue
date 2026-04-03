<script setup lang="ts">
import { cn } from '@/lib/utils'
import { ChevronLeft } from 'lucide-vue-next'
import { PaginationPrev, type PaginationPrevProps } from 'radix-vue'
import { Button } from '@/components/ui/button'
import type { HTMLAttributes } from 'vue'

const props = withDefaults(
  defineProps<PaginationPrevProps & { class?: HTMLAttributes['class'] }>(),
  {
    asChild: true,
  },
)
</script>

<template>
  <PaginationPrev v-bind="props">
    <Button :class="cn('w-9 h-9 p-0', props.class)" variant="outline">
      <slot>
        <ChevronLeft />
      </slot>
    </Button>
  </PaginationPrev>
</template>
