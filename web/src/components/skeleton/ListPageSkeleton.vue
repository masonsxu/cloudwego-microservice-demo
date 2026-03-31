<template>
  <div class="table-skeleton">
    <!-- Table header -->
    <div class="skeleton-table-header">
      <div
        v-for="i in columns"
        :key="i"
        class="skeleton-bone table-header-bone"
        :style="{ flex: i === 1 ? 1.5 : 1 }"
      />
    </div>
    <!-- Table rows -->
    <div
      v-for="row in rows"
      :key="row"
      class="skeleton-table-row"
      :style="{ animationDelay: `${row * 0.04}s` }"
    >
      <div
        v-for="col in columns"
        :key="col"
        class="skeleton-table-cell"
        :style="{ flex: col === 1 ? 1.5 : 1 }"
      >
        <div v-if="col === 1" class="cell-with-avatar">
          <div class="skeleton-bone avatar-bone" />
          <div class="skeleton-bone cell-text-bone" />
        </div>
        <div v-else class="skeleton-bone cell-text-bone" :style="{ width: getCellWidth(col) }" />
      </div>
    </div>
    <!-- Pagination -->
    <div class="skeleton-pagination">
      <div class="skeleton-bone pagination-bone" />
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

<style scoped lang="scss">
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
    rgba(212, 175, 55, 0.04) 0%,
    rgba(212, 175, 55, 0.10) 40%,
    rgba(212, 175, 55, 0.04) 80%
  );
  background-size: 800px 100%;
  animation: shimmer 1.8s ease-in-out infinite;
  border-radius: 6px;
}

.table-skeleton {
  .skeleton-table-header {
    display: flex;
    gap: 16px;
    padding: 16px 20px;
    border-bottom: 1px solid var(--c-border);

    .table-header-bone {
      height: 12px;
      width: 60%;
      border-radius: 3px;
    }
  }

  .skeleton-table-row {
    display: flex;
    gap: 16px;
    padding: 14px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.02);
    animation: fadeSlideIn 0.25s ease-out backwards;

    &:last-of-type { border-bottom: none; }
  }

  .skeleton-table-cell {
    display: flex;
    align-items: center;
    min-width: 0;
  }

  .cell-with-avatar {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;

    .avatar-bone {
      width: 30px;
      height: 30px;
      border-radius: 8px;
      flex-shrink: 0;
    }

    .cell-text-bone {
      width: 70%;
      height: 14px;
    }
  }

  .cell-text-bone {
    height: 14px;
  }

  .skeleton-pagination {
    display: flex;
    justify-content: flex-end;
    padding: 16px 20px;
    border-top: 1px solid var(--c-border);

    .pagination-bone {
      width: 320px;
      height: 28px;
      border-radius: 6px;
    }
  }
}
</style>
