<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { Primitive, type PrimitiveProps } from 'radix-vue'
import { cva, type VariantProps } from 'class-variance-authority'

const buttonVariants = cva(
  'inline-flex items-center justify-center gap-1.5 whitespace-nowrap rounded-sm text-[14px] font-semibold transition-colors duration-[var(--duration-fast)] focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0',
  {
    variants: {
      variant: {
        default:
          'bg-[color:var(--color-primary)] text-[color:var(--color-ink-on-primary)] hover:bg-[color:var(--color-primary-hover)] active:bg-[color:var(--color-primary-active)]',
        destructive:
          'bg-[color:var(--color-danger)] text-white hover:bg-[color:var(--color-danger-ink)]',
        outline:
          'border border-default bg-canvas text-ink hover:bg-sunken hover:border-strong',
        secondary:
          'bg-sunken text-ink hover:bg-[color:var(--color-border-subtle)]',
        ghost:
          'text-[color:var(--color-ink-muted)] hover:bg-sunken hover:text-ink',
        link:
          'text-[color:var(--color-primary)] underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-9 px-3.5',
        sm: 'h-8 px-3 text-[13px]',
        lg: 'h-10 px-5',
        icon: 'h-8 w-8',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

type ButtonVariants = VariantProps<typeof buttonVariants>

interface Props extends PrimitiveProps {
  variant?: NonNullable<ButtonVariants['variant']>
  size?: NonNullable<ButtonVariants['size']>
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  as: 'button',
})
</script>

<template>
  <Primitive
    :as="as"
    :as-child="asChild"
    :class="cn(buttonVariants({ variant, size }), props.class)"
  >
    <slot />
  </Primitive>
</template>
