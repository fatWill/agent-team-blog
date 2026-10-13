<template>
  <div class="tuyile-root min-h-screen bg-stone-50 text-gray-800 transition-colors duration-300 dark:bg-gray-900 dark:text-gray-200">
    <!-- 返回玩具区 -->
    <NuxtLink
      to="/toys"
      class="fixed left-4 top-4 z-50 flex h-9 items-center gap-1.5 rounded-full bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
    >
      <span aria-hidden="true">←</span>
      <span>玩具</span>
    </NuxtLink>

    <slot name="nav" />

    <!-- 主题切换 -->
    <button
      type="button"
      class="fixed right-4 top-4 z-50 flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
      aria-label="切换主题"
      @click="toggleTheme"
    >
      <svg v-if="isDark" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
      <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
    </button>

    <slot />
  </div>
</template>

<script setup lang="ts">
import '../styles/tuyile.css'

const { isDark, toggleTheme } = useTheme()
</script>
