<template>
  <!-- PC 端：使用 ant-design-vue Spin -->
  <div v-if="!isMobile" class="flex items-center justify-center py-20">
    <ASpin size="large" :tip="tip" />
  </div>
  <!-- 移动端：自定义 loading -->
  <div v-else class="flex items-center justify-center py-20">
    <div class="mobile-loading-wrap">
      <svg class="mobile-loading-spinner" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      <span v-if="tip" class="mobile-loading-tip">{{ tip }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Spin as ASpin } from 'ant-design-vue'
import 'ant-design-vue/es/spin/style'

defineProps<{
  tip?: string
}>()

const { isMobile } = useDevice()
</script>

<style scoped>
.mobile-loading-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 20px 0;
}
.mobile-loading-spinner {
  width: 32px;
  height: 32px;
  color: #1677ff;
  animation: spin 0.8s linear infinite;
}
.mobile-loading-tip {
  font-size: 14px;
  color: #999;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

:root.dark .mobile-loading-spinner {
  color: #4096ff;
}
:root.dark .mobile-loading-tip {
  color: #666;
}
</style>
