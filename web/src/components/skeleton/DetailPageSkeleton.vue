<template>
  <div class="grid gap-5 lg:grid-cols-3">
    <div class="lg:col-span-1">
      <div
        v-for="n in sideCards"
        :key="'s' + n"
        class="sk-card"
        :style="{ animationDelay: `${(n - 1) * 0.1}s` }"
      >
        <div class="sk-card-header">
          <div class="skeleton-bone" style="width: 18px; height: 18px; border-radius: 4px" />
          <div class="skeleton-bone" style="width: 35%; height: 14px" />
        </div>
        <div class="sk-card-body">
          <div v-if="n === 1 && showAvatar" class="sk-avatar-section">
            <div class="skeleton-bone sk-avatar" />
            <div class="skeleton-bone" style="width: 60px; height: 24px; border-radius: 11px" />
          </div>
          <div
            v-for="i in (sideItemCounts[n - 1] ?? 4)"
            :key="i"
            class="sk-field"
            :style="{ animationDelay: `${i * 0.04}s` }"
          >
            <div class="skeleton-bone sk-label" />
            <div class="skeleton-bone sk-value" :style="{ width: widths[(i - 1) % widths.length] }" />
          </div>
        </div>
      </div>
    </div>
    <div class="lg:col-span-2">
      <div
        v-for="n in mainCards"
        :key="'m' + n"
        class="sk-card"
        :style="{ animationDelay: `${(n - 1) * 0.1 + 0.05}s` }"
      >
        <div class="sk-card-header">
          <div class="skeleton-bone" style="width: 18px; height: 18px; border-radius: 4px" />
          <div class="skeleton-bone" style="width: 28%; height: 14px" />
        </div>
        <div class="sk-card-body">
          <div
            v-for="i in (mainItemCounts[n - 1] ?? 6)"
            :key="i"
            class="sk-field"
            :style="{ animationDelay: `${i * 0.04}s` }"
          >
            <div class="skeleton-bone sk-label" />
            <div class="skeleton-bone sk-value" :style="{ width: widths[(i + 2) % widths.length] }" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  sideSpan?: number
  mainSpan?: number
  sideCards?: number
  mainCards?: number
  showAvatar?: boolean
  sideItemCounts?: number[]
  mainItemCounts?: number[]
}>(), {
  sideSpan: 8,
  mainSpan: 16,
  sideCards: 1,
  mainCards: 2,
  showAvatar: true,
  sideItemCounts: () => [4],
  mainItemCounts: () => [6, 4],
})

const widths = ['60%', '75%', '50%', '45%', '65%', '55%', '70%']
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

.sk-card {
  background: var(--color-canvas);
  border: 1px solid var(--color-border-subtle);
  border-radius: 8px;
  margin-bottom: 16px;
  animation: fadeSlideIn 0.3s ease-out backwards;
  overflow: hidden;
}

.sk-card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--color-border-subtle);
}

.sk-card-body {
  padding: 20px;
}

.sk-avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.sk-avatar {
  width: 100px;
  height: 100px;
  border-radius: 9999px;
}

.sk-field {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--color-border-subtle);
  animation: fadeSlideIn 0.25s ease-out backwards;
}

.sk-field:last-child {
  border-bottom: none;
}

.sk-label {
  width: 80px;
  height: 14px;
  flex-shrink: 0;
  margin-right: 20px;
}

.sk-value {
  height: 14px;
}
</style>
