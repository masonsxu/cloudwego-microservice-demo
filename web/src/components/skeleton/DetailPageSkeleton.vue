<template>
  <el-row :gutter="20">
    <el-col :xs="24" :lg="sideSpan">
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
    </el-col>
    <el-col :xs="24" :lg="mainSpan">
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
    </el-col>
  </el-row>
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

.sk-card {
  background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
  margin-bottom: 20px;
  animation: fadeSlideIn 0.3s ease-out backwards;
  overflow: hidden;
}

.sk-card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(212, 175, 55, 0.2);
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
  border-radius: 50%;
}

.sk-field {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  animation: fadeSlideIn 0.25s ease-out backwards;

  &:last-child {
    border-bottom: none;
  }
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
