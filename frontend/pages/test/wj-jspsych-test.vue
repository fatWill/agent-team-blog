<template>
  <div class="min-h-screen flex items-center justify-center bg-white text-gray-500">
    实验加载中，请稍候…
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

useHead({ title: 'wj-jspsych-test | 素材测试' })

onMounted(() => {
  // 保留 query 参数，兼容问卷平台带 group_id/sample_id 之类透传
  const target = '/test/wj-jspsych-test/index.html' + (window.location.search || '')
  window.location.replace(target)
})
</script>
