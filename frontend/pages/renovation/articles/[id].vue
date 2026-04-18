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
        <span class="truncate text-sm font-medium text-gray-900 dark:text-gray-100">{{ article?.title || '加载中...' }}</span>
      </div>
    </header>

    <!-- 加载状态 -->
    <main v-if="loading" class="mx-auto w-full max-w-3xl flex-1 px-4 py-10">
      <div class="animate-pulse space-y-4">
        <div class="h-8 w-2/3 rounded bg-gray-200 dark:bg-gray-700" />
        <div class="h-4 w-full rounded bg-gray-200 dark:bg-gray-700" />
        <div class="h-4 w-5/6 rounded bg-gray-200 dark:bg-gray-700" />
        <div class="h-4 w-4/5 rounded bg-gray-200 dark:bg-gray-700" />
      </div>
    </main>

    <!-- 错误状态 -->
    <main v-else-if="error" class="mx-auto w-full max-w-3xl flex-1 px-4 py-20 text-center">
      <p class="text-lg font-medium text-gray-600 dark:text-gray-400">文章不存在或加载失败</p>
      <NuxtLink to="/renovation" class="mt-4 inline-block text-sm text-primary-500 hover:text-primary-600">
        返回装修主页
      </NuxtLink>
    </main>

    <!-- 文章内容 -->
    <main v-else-if="article" class="mx-auto w-full max-w-3xl flex-1 px-4 py-10">
      <article>
        <h1 class="mb-6 text-2xl font-bold leading-tight text-gray-900 dark:text-gray-100 md:text-3xl">
          {{ article.title }}
        </h1>
        <div
          class="markdown-body"
          v-html="renderedContent"
        />
      </article>
    </main>
  </div>
</template>

<script setup lang="ts">
import MarkdownIt from 'markdown-it'
import { apiFetchRenovationArticle } from '~/features/renovation'

const route = useRoute()
const articleId = route.params.id as string

const loading = ref(true)
const error = ref(false)
const article = ref<{ id: number; title: string; content: string } | null>(null)

// 初始化 markdown-it
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
})

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  return md.render(article.value.content)
})

// SSR 预取
try {
  const data = await apiFetchRenovationArticle(articleId)
  if (data) {
    article.value = data
  } else {
    error.value = true
  }
} catch {
  error.value = true
} finally {
  loading.value = false
}

// SEO
useSeoMeta({
  title: () => article.value ? `${article.value.title} - 装修知识` : '装修知识',
  description: () => article.value?.title || '装修知识文章',
})
</script>

<style>
/* Markdown 渲染样式 */
.markdown-body {
  color: #1f2937;
  line-height: 1.75;
  font-size: 0.9375rem;
}
.dark .markdown-body {
  color: #d1d5db;
}

.markdown-body h1 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 2rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e5e7eb;
}
.dark .markdown-body h1 {
  border-bottom-color: #374151;
}

.markdown-body h2 {
  font-size: 1.25rem;
  font-weight: 700;
  margin-top: 1.75rem;
  margin-bottom: 0.75rem;
  color: #111827;
}
.dark .markdown-body h2 {
  color: #f3f4f6;
}

.markdown-body h3 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
  color: #1f2937;
}
.dark .markdown-body h3 {
  color: #e5e7eb;
}

.markdown-body p {
  margin-bottom: 1rem;
}

.markdown-body ul, .markdown-body ol {
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}
.markdown-body ul {
  list-style-type: disc;
}
.markdown-body ol {
  list-style-type: decimal;
}
.markdown-body li {
  margin-bottom: 0.375rem;
}
.markdown-body li > ul, .markdown-body li > ol {
  margin-top: 0.375rem;
  margin-bottom: 0;
}

.markdown-body strong {
  font-weight: 700;
  color: #111827;
}
.dark .markdown-body strong {
  color: #f9fafb;
}

.markdown-body blockquote {
  border-left: 4px solid #6366f1;
  padding: 0.75rem 1rem;
  margin: 1rem 0;
  background: #f5f3ff;
  border-radius: 0 0.5rem 0.5rem 0;
  color: #4338ca;
}
.dark .markdown-body blockquote {
  background: rgba(99, 102, 241, 0.1);
  border-left-color: #818cf8;
  color: #a5b4fc;
}

.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
  font-size: 0.875rem;
}
.markdown-body th {
  background: #f3f4f6;
  font-weight: 600;
  text-align: left;
  padding: 0.625rem 0.75rem;
  border: 1px solid #e5e7eb;
}
.dark .markdown-body th {
  background: #1f2937;
  border-color: #374151;
}
.markdown-body td {
  padding: 0.5rem 0.75rem;
  border: 1px solid #e5e7eb;
}
.dark .markdown-body td {
  border-color: #374151;
}

.markdown-body code {
  background: #f3f4f6;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-size: 0.85em;
  color: #dc2626;
}
.dark .markdown-body code {
  background: #1f2937;
  color: #f87171;
}

.markdown-body pre {
  background: #1f2937;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1rem 0;
}
.markdown-body pre code {
  background: none;
  padding: 0;
  color: #e5e7eb;
}

.markdown-body hr {
  border: none;
  border-top: 1px solid #e5e7eb;
  margin: 2rem 0;
}
.dark .markdown-body hr {
  border-top-color: #374151;
}

.markdown-body a {
  color: #6366f1;
  text-decoration: underline;
  text-underline-offset: 2px;
}
.dark .markdown-body a {
  color: #818cf8;
}

/* 任务列表样式 */
.markdown-body input[type="checkbox"] {
  margin-right: 0.5rem;
  accent-color: #6366f1;
}
</style>
