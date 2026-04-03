<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <div
    :class="
      cn(
        'relative w-full overflow-auto',
        props.class,
      )
    "
  >
    <table :class="cn('w-full caption-bottom text-sm', $attrs.class ?? '')">
      <slot />
    </table>
  </div>
</template>
