<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
const { colorMode } = useTheme()
const { siteUrl, assetsUrl, cdnUrl } = useSiteUrl()

// SSR 阶段直接把 dark class 注入 <html>，客户端接管时已经有正确的类名，彻底消灭主题闪烁
// canonical / preconnect / dns-prefetch 从 nuxt.config.ts 迁移到此处，
// 以便通过 runtimeConfig 在运行时切换域名（页面级 canonical 会覆盖此处的站点根 canonical）
useHead({
  htmlAttrs: {
    class: computed(() => colorMode.value === 'dark' ? 'dark' : ''),
  },
  link: [
    { rel: 'canonical', href: siteUrl },
    { rel: 'preconnect', href: assetsUrl, crossorigin: '' },
    { rel: 'dns-prefetch', href: assetsUrl },
    { rel: 'preconnect', href: cdnUrl, crossorigin: '' },
    { rel: 'dns-prefetch', href: cdnUrl },
  ],
})
</script>
