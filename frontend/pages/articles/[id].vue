<template>
  <div class="min-h-screen">
    <!-- 右上角固定：主题切换 -->
    <div class="fixed right-4 top-4 z-50">
      <button
        class="flex h-9 w-9 items-center justify-center rounded-lg bg-white/80 text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
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
    </div>

    <!-- 左上角固定：返回首页 -->
    <div class="fixed left-4 top-4 z-50">
      <NuxtLink
        to="/"
        class="flex h-9 items-center gap-1.5 rounded-lg bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        <span class="hidden sm:inline">返回首页</span>
      </NuxtLink>
    </div>

    <!-- 文章内容 -->
    <main class="mx-auto max-w-3xl px-4 pt-16 pb-10">
      <!-- 加载状态 -->
      <div v-if="loading" class="space-y-4">
        <div class="h-8 w-3/4 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
        <div class="h-4 w-1/3 animate-pulse rounded bg-gray-100 dark:bg-gray-700/50" />
        <div class="mt-8 space-y-3">
          <div v-for="i in 8" :key="i" class="h-4 animate-pulse rounded bg-gray-100 dark:bg-gray-700/50" :class="i === 8 ? 'w-2/3' : 'w-full'" />
        </div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="py-20 text-center">
        <p class="text-red-500">文章加载失败</p>
        <NuxtLink to="/" class="mt-3 inline-block text-sm text-primary-500 hover:text-primary-600">返回首页</NuxtLink>
      </div>

      <!-- 文章详情 -->
      <article v-else-if="article">
        <h1 class="text-2xl font-bold leading-tight text-gray-900 sm:text-3xl dark:text-gray-100">
          {{ article.title }}
        </h1>
        <div class="mt-3 flex flex-wrap items-center gap-2 text-sm text-gray-400 sm:gap-3 dark:text-gray-500">
          <time>{{ formatDate(article.createdAt) }}</time>
          <span v-if="article.updatedAt !== article.createdAt">
            · 更新于 {{ formatDate(article.updatedAt) }}
          </span>
        </div>
        <p v-if="article.summary" class="mt-4 text-sm leading-relaxed text-gray-500 sm:text-base dark:text-gray-400">
          {{ article.summary }}
        </p>

        <!-- Tiptap 渲染区域 -->
        <div class="tiptap-renderer mt-8">
          <EditorContent v-if="editor" :editor="editor" />
        </div>
      </article>
    </main>

    <!-- 页脚 -->
    <footer class="border-t border-gray-200/60 py-8 text-center text-sm text-gray-400 transition-colors dark:border-gray-700/60 dark:text-gray-500">
      <p>&copy; {{ new Date().getFullYear() }} fatwillzeng. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import type { ArticleDetail } from '~/types'
import { apiFetchArticle } from '~/utils/api'

const route = useRoute()
const { isDark, toggleTheme } = useTheme()

const article = ref<ArticleDetail | null>(null)
const loading = ref(true)
const error = ref(false)

// Tiptap 只读编辑器
const editor = useEditor({
  editable: false,
  extensions: [
    StarterKit,
    Image,
    Link.configure({
      openOnClick: true,
      HTMLAttributes: {
        target: '_blank',
        rel: 'noopener noreferrer',
      },
    }),
  ],
})

async function loadArticle() {
  const id = route.params.id as string
  loading.value = true
  error.value = false
  try {
    const data = await apiFetchArticle(id)
    article.value = data
    if (editor.value && data.content) {
      editor.value.commands.setContent(data.content)
    }
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

onMounted(() => {
  loadArticle()
})

watch(article, (val) => {
  if (val) {
    useHead({
      title: `${val.title} - fatwillzeng`,
      meta: [
        { name: 'description', content: val.summary || val.title },
      ],
    })
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>
