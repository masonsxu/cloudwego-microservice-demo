<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { cva, type VariantProps } from 'class-variance-authority'

const badgeVariants = cva(
  'inline-flex items-center rounded-xs px-1.5 py-0.5 text-[12px] font-medium leading-tight transition-colors focus:outline-none',
  {
    variants: {
      variant: {
        default:
          'bg-sunken text-[color:var(--color-ink-muted)]',
        primary:
          'bg-[color:var(--color-primary-soft)] text-[color:var(--color-primary-active)]',
        success:
          'bg-[color:var(--color-success-soft)] text-[color:var(--color-success-ink)]',
        warning:
          'bg-[color:var(--color-warning-soft)] text-[color:var(--color-warning-ink)]',
        destructive:
          'bg-[color:var(--color-danger-soft)] text-[color:var(--color-danger-ink)]',
        info:
          'bg-[color:var(--color-info-soft)] text-[color:var(--color-info-ink)]',
        secondary:
          'bg-sunken text-[color:var(--color-ink-muted)]',
        outline:
          'border border-default text-ink',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  },
)

type BadgeVariants = VariantProps<typeof badgeVariants>

interface Props {
  variant?: NonNullable<BadgeVariants['variant']>
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
  <div :class="cn(badgeVariants({ variant }), props.class)">
    <slot />
  </div>
</template>
