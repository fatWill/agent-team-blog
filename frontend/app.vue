<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
const { colorMode } = useTheme()

// SSR 阶段直接把 dark class 注入 <html>，客户端接管时已经有正确的类名，彻底消灭主题闪烁
useHead({
  htmlAttrs: {
    class: computed(() => colorMode.value === 'dark' ? 'dark' : ''),
  },
})
</script>
