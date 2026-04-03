<script setup lang="ts">
import { cn } from '@/lib/utils'
import { ChevronsRight } from 'lucide-vue-next'
import { PaginationLast, type PaginationLastProps } from 'radix-vue'
import { Button } from '@/components/ui/button'
import type { HTMLAttributes } from 'vue'

const props = withDefaults(
  defineProps<PaginationLastProps & { class?: HTMLAttributes['class'] }>(),
  {
    asChild: true,
  },
)
</script>

<template>
  <PaginationLast v-bind="props">
    <Button :class="cn('w-9 h-9 p-0', props.class)" variant="outline">
      <slot>
        <ChevronsRight />
      </slot>
    </Button>
  </PaginationLast>
</template>
