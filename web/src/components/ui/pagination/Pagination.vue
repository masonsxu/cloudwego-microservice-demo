<script setup lang="ts">
import { PaginationRoot } from 'radix-vue'
import type { HTMLAttributes } from 'vue'

const props = defineProps<{
  total: number
  itemsPerPage?: number
  defaultPage?: number
  siblingCount?: number
  showEdges?: boolean
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <PaginationRoot
    :total="total"
    :items-per-page="itemsPerPage"
    :default-page="defaultPage"
    :sibling-count="siblingCount"
    :show-edges="showEdges"
  >
    <slot />
  </PaginationRoot>
</template>
