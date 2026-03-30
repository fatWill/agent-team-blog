<template>
  <div class="min-h-screen">
    <!-- 顶部导航栏 -->
    <header class="sticky top-0 z-50 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg transition-colors duration-300 dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-3xl items-center justify-between px-4">
        <NuxtLink to="/" class="text-lg font-bold text-gray-900 dark:text-gray-100">
          fatwillzeng
        </NuxtLink>
        <div class="flex items-center gap-2">
          <!-- GitHub -->
          <a
            href="https://github.com/fatwillzeng"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
            aria-label="GitHub"
          >
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
          </a>
          <!-- 主题切换 -->
          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
            aria-label="切换主题"
            @click="toggleTheme"
          >
            <!-- 太阳图标（dark 模式下显示，点击切换为 light） -->
            <svg v-if="isDark" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <!-- 月亮图标（light 模式下显示，点击切换为 dark） -->
            <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- 个人信息区域 -->
    <section class="mx-auto max-w-3xl px-4 pt-10 pb-6">
      <div class="flex items-center gap-5">
        <img
          src="/avatar.png"
          alt="fatwillzeng 头像"
          class="h-20 w-20 rounded-full object-cover ring-2 ring-gray-200 dark:ring-gray-700"
        />
        <div class="flex-1">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">fatwillzeng</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            全栈开发者 / 热爱技术与生活 / 记录成长的点滴
          </p>
        </div>
      </div>
    </section>

    <!-- Tab 导航 -->
    <nav class="sticky top-14 z-40 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg transition-colors duration-300 dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex max-w-3xl gap-0 px-4">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="relative px-4 py-3 text-sm font-medium transition-colors"
          :class="activeTab === tab.key
            ? 'text-gray-900 dark:text-gray-100'
            : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
          <!-- 激活指示条 -->
          <span
            v-if="activeTab === tab.key"
            class="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary-500"
          />
        </button>
      </div>
    </nav>

    <!-- 内容区域 -->
    <main class="mx-auto max-w-3xl px-4 py-8">
      <!-- 文章 Tab -->
      <div v-if="activeTab === 'articles'">
        <!-- 加载状态 -->
        <div v-if="loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="animate-pulse rounded-xl border border-gray-100 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
            <div class="mb-3 h-5 w-3/4 rounded bg-gray-200 dark:bg-gray-700" />
            <div class="mb-2 h-4 w-full rounded bg-gray-100 dark:bg-gray-700/50" />
            <div class="h-4 w-1/2 rounded bg-gray-100 dark:bg-gray-700/50" />
          </div>
        </div>

        <!-- 错误状态 -->
        <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-6 text-center dark:border-red-800 dark:bg-red-900/20">
          <p class="text-red-600 dark:text-red-400">加载文章失败，请稍后重试</p>
          <button class="mt-3 text-sm text-primary-500 hover:text-primary-600" @click="fetchArticles">
            重新加载
          </button>
        </div>

        <!-- 文章列表 -->
        <div v-else-if="articles.length > 0" class="space-y-4">
          <NuxtLink
            v-for="article in articles"
            :key="article.id"
            :to="`/articles/${article.id}`"
            class="block rounded-xl border border-gray-100 bg-white p-5 transition-all duration-200 hover:shadow-md hover:-translate-y-0.5 dark:border-gray-700 dark:bg-gray-800 dark:hover:border-gray-600"
          >
            <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-gray-100">
              {{ article.title }}
            </h2>
            <p v-if="article.summary" class="mb-3 text-sm leading-relaxed text-gray-500 dark:text-gray-400 line-clamp-2">
              {{ article.summary }}
            </p>
            <time class="text-xs text-gray-400 dark:text-gray-500">
              {{ formatDate(article.createdAt) }}
            </time>
          </NuxtLink>
        </div>

        <!-- 空状态 -->
        <div v-else class="py-20 text-center">
          <p class="text-gray-400 dark:text-gray-500">暂无文章</p>
        </div>
      </div>

      <!-- 其他 Tab -->
      <div v-else class="py-20 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
          <svg class="h-8 w-8 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
          </svg>
        </div>
        <p class="text-lg font-medium text-gray-600 dark:text-gray-400">敬请期待</p>
        <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">该板块正在建设中...</p>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="border-t border-gray-200/60 py-8 text-center text-sm text-gray-400 transition-colors dark:border-gray-700/60 dark:text-gray-500">
      <p>© {{ new Date().getFullYear() }} fatwillzeng. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import type { ArticleListItem, TabItem } from '~/types'
import { apiFetchArticles } from '~/utils/api'

const { isDark, toggleTheme } = useTheme()

// Tab 导航
const tabs: TabItem[] = [
  { key: 'articles', label: '文章' },
  { key: 'life', label: '生活' },
  { key: 'tools', label: '小工具·小游戏' },
  { key: 'agent-team', label: 'agent team' },
]
const activeTab = ref('articles')

// 文章数据
const articles = ref<ArticleListItem[]>([])
const loading = ref(false)
const error = ref(false)

async function fetchArticles() {
  loading.value = true
  error.value = false
  try {
    const res = await apiFetchArticles()
    articles.value = res.list
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// 页面挂载时加载文章
onMounted(() => {
  fetchArticles()
})

useHead({
  title: 'fatwillzeng - 个人博客',
  meta: [
    { name: 'description', content: '全栈开发者 fatwillzeng 的个人博客，分享技术与生活。' },
  ],
})
</script>
