<template>
  <div class="table-skeleton">
    <!-- Table header -->
    <div class="flex gap-4 px-4 py-3 border-b border-subtle">
      <div
        v-for="i in columns"
        :key="i"
        class="skeleton-bone h-3 w-[60%] rounded-xs"
        :style="{ flex: i === 1 ? 1.5 : 1 }"
      />
    </div>
    <!-- Table rows -->
    <div
      v-for="row in rows"
      :key="row"
      class="flex gap-4 px-4 py-3 border-b border-subtle animate-fadeSlideIn"
      :style="{ animationDelay: `${row * 0.04}s`, animationFillMode: 'backwards' }"
    >
      <div
        v-for="col in columns"
        :key="col"
        class="flex items-center min-w-0"
        :style="{ flex: col === 1 ? 1.5 : 1 }"
      >
        <div v-if="col === 1" class="flex items-center gap-2.5 w-full">
          <div class="skeleton-bone w-6 h-6 rounded-full flex-shrink-0" />
          <div class="skeleton-bone h-3.5 w-[70%]" />
        </div>
        <div v-else class="skeleton-bone h-3.5" :style="{ width: getCellWidth(col) }" />
      </div>
    </div>
    <!-- Pagination -->
    <div class="flex justify-end px-4 py-3 border-t border-subtle">
      <div class="skeleton-bone w-80 h-7 rounded-sm" />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  rows?: number
  columns?: number
}>(), {
  rows: 8,
  columns: 6,
})

function getCellWidth(col: number): string {
  const widths = ['60%', '70%', '50%', '45%', '55%', '40%', '65%']
  return widths[(col - 1) % widths.length] ?? '60%'
}
</script>

<style scoped>
@keyframes shimmer {
  0% { background-position: -400px 0; }
  100% { background-position: 400px 0; }
}

@keyframes fadeSlideIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

.skeleton-bone {
  background: linear-gradient(
    90deg,
    var(--color-sunken) 0%,
    var(--color-border-subtle) 50%,
    var(--color-sunken) 100%
  );
  background-size: 800px 100%;
  animation: shimmer 1.6s ease-in-out infinite;
  border-radius: 4px;
}

.animate-fadeSlideIn {
  animation-name: fadeSlideIn;
  animation-duration: 0.25s;
  animation-timing-function: ease-out;
}
</style>
