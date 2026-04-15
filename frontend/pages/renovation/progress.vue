<template>
  <div class="flex min-h-screen flex-col bg-gray-50 transition-colors dark:bg-gray-900">
    <!-- 顶部导航 -->
    <header class="sticky top-0 z-30 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-4xl items-center px-4">
        <NuxtLink
          to="/renovation"
          class="flex items-center gap-2 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
          返回装修
        </NuxtLink>
        <span class="mx-3 text-gray-300 dark:text-gray-600">/</span>
        <span class="text-sm font-medium text-gray-900 dark:text-gray-100">装修进度</span>
      </div>
    </header>

    <main class="mx-auto w-full max-w-4xl flex-1 px-4 py-10">
      <div class="py-20 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-blue-50 dark:bg-blue-900/20">
          <span class="text-3xl">📋</span>
        </div>
        <h1 class="text-xl font-bold text-gray-900 dark:text-gray-100">装修进度</h1>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">功能正在开发中，敬请期待...</p>
      </div>
    </main>
  </div>
</template>
