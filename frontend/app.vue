<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <!-- 全局悬浮留言板气泡（非 admin/login 页面显示） -->
  <GuestbookBubble v-if="showGuestbook" />
</template>

<script setup lang="ts">
const { colorMode } = useTheme()
const route = useRoute()

const showGuestbook = computed(() => {
  const path = route.path
  const allowedPaths = ['/home', '/articles', '/life', '/toys', '/agent-team', '/guestbook', '/changelog']
  return allowedPaths.some(p => path === p || path.startsWith(p + '/'))
})

// SSR 阶段直接把 dark class 注入 <html>，客户端接管时已经有正确的类名，彻底消灭主题闪烁
useHead({
  htmlAttrs: {
    class: computed(() => colorMode.value === 'dark' ? 'dark' : ''),
  },
})
</script>
