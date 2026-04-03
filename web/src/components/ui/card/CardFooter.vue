<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <footer
    :class="cn('flex items-center p-6 pt-0', props.class)"
  >
    <slot />
  </footer>
</template>
