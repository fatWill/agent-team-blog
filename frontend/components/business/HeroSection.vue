<template>
  <section class="relative overflow-hidden bg-gradient-to-br from-primary-500 via-primary-600 to-primary-700">
    <!-- 装饰背景 -->
    <div class="absolute inset-0">
      <div class="absolute -top-24 -right-24 h-96 w-96 rounded-full bg-white/5 blur-3xl" />
      <div class="absolute -bottom-24 -left-24 h-96 w-96 rounded-full bg-accent-400/10 blur-3xl" />
    </div>

    <div class="relative mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28 lg:px-8">
      <div class="max-w-2xl">
        <div class="mb-4 inline-flex items-center rounded-full bg-white/10 px-4 py-1.5 text-sm font-medium text-white/90 backdrop-blur-sm">
          ✨ 欢迎来到 fatwill
        </div>
        <h1 class="mb-6 text-4xl font-extrabold tracking-tight text-white sm:text-5xl lg:text-6xl">
          fatwill最棒
          <br />
          <span class="text-primary-200">无限可能</span>
        </h1>
        <p class="mb-8 max-w-lg text-lg leading-relaxed text-primary-100">
          在这里，我们分享前沿技术趋势、深度实践经验与独到的工程思考。与志同道合的开发者一起，探索技术的无限可能。
        </p>
        <div class="flex flex-wrap gap-3">
          <NuxtLink
            to="/articles"
            class="inline-flex items-center gap-2 rounded-lg bg-white px-6 py-3 text-sm font-semibold text-primary-600 shadow-lg shadow-primary-900/20 transition-all duration-200 hover:bg-primary-50 hover:-translate-y-0.5"
          >
            开始阅读
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </NuxtLink>
          <NuxtLink
            to="/about"
            class="inline-flex items-center gap-2 rounded-lg border border-white/20 bg-white/10 px-6 py-3 text-sm font-semibold text-white backdrop-blur-sm transition-all duration-200 hover:bg-white/20 hover:-translate-y-0.5"
          >
            了解更多
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
</script>
