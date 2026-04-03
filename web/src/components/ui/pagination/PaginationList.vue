<script setup lang="ts">
import { cn } from '@/lib/utils'
import { PaginationList, type PaginationListProps } from 'radix-vue'
import type { HTMLAttributes } from 'vue'

const props = defineProps<PaginationListProps & { class?: HTMLAttributes['class'] }>()
</script>

<template>
  <PaginationList v-slot="{ items }" :class="cn('flex items-center gap-1', props.class)">
    <slot :items="items" />
  </PaginationList>
</template>
