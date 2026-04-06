<template>
  <div class="min-h-screen">
    <!-- 顶部阅读进度条 -->
    <div class="fixed left-0 top-0 z-[60] h-0.5 w-full">
      <div
        class="h-full bg-primary-500 transition-[width] duration-100"
        :style="{ width: readProgress + '%' }"
      />
    </div>

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
        to="/home"
        class="flex h-9 items-center gap-1.5 rounded-lg bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        <span class="hidden sm:inline">返回首页</span>
      </NuxtLink>
    </div>

    <!-- 文章内容 + 右侧 TOC -->
    <div class="mx-auto flex max-w-5xl gap-8 px-4 pt-16 pb-10">
      <!-- 文章主体 -->
      <main class="min-w-0 flex-1 max-w-3xl">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>

        <!-- 错误状态 -->
        <div v-else-if="error" class="py-20 text-center">
          <p class="text-red-500">文章加载失败</p>
          <NuxtLink to="/home" class="mt-3 inline-block text-sm text-primary-500 hover:text-primary-600">返回首页</NuxtLink>
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
            <span class="flex items-center gap-1">
              · 👁 {{ articleViews }}
            </span>
          </div>
          <p v-if="article.summary" class="mt-4 text-sm leading-relaxed text-gray-500 sm:text-base dark:text-gray-400">
            {{ article.summary }}
          </p>

          <!-- Tiptap 渲染区域 -->
          <div ref="articleContentRef" class="tiptap-renderer mt-8">
            <EditorContent v-if="editor" :editor="editor" />
          </div>

          <!-- 点赞按钮 -->
          <div class="mt-10 flex justify-center border-t border-gray-200/60 pt-6 dark:border-gray-700/60">
            <button
              class="group flex items-center gap-2 rounded-full border border-gray-200 px-6 py-3 transition-all duration-200 hover:border-red-200 hover:bg-red-50 dark:border-gray-600 dark:hover:border-red-800 dark:hover:bg-red-900/20"
              @click="handleArticleLike"
            >
              <svg
                class="h-5 w-5 transition-all duration-300"
                :class="[
                  articleLiked ? 'text-red-500 fill-red-500' : 'text-gray-400 fill-none stroke-gray-400 dark:text-gray-500 dark:stroke-gray-500 group-hover:text-red-400 group-hover:stroke-red-400',
                  likeAnimating ? 'scale-125' : '',
                ]"
                viewBox="0 0 24 24"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
              </svg>
              <span class="text-sm font-medium" :class="articleLiked ? 'text-red-500' : 'text-gray-500 dark:text-gray-400 group-hover:text-red-400'">
                {{ articleLikeCount > 0 ? articleLikeCount : '点赞' }}
              </span>
            </button>
          </div>
        </article>
      </main>

      <!-- 右侧 TOC（仅 PC 端） -->
      <aside v-if="tocItems.length > 0" class="hidden w-52 shrink-0 md:block">
        <nav class="sticky top-20">
          <h4 class="mb-3 text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">目录</h4>
          <ul class="space-y-1 border-l border-gray-200 dark:border-gray-700">
            <li v-for="item in tocItems" :key="item.id">
              <button
                class="block w-full truncate border-l-2 py-1 text-left text-xs transition-colors"
                :class="[
                  activeTocId === item.id
                    ? 'border-primary-500 font-medium text-primary-600 dark:text-primary-400'
                    : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300',
                  item.level === 1 ? 'pl-3' : item.level === 2 ? 'pl-5' : 'pl-7',
                ]"
                @click="scrollToHeading(item.id)"
              >
                {{ item.text }}
              </button>
            </li>
          </ul>
        </nav>
      </aside>
    </div>

    <!-- 页脚 -->
    <footer class="border-t border-gray-200/60 py-8 text-center text-sm text-gray-400 transition-colors dark:border-gray-700/60 dark:text-gray-500">
      <p>&copy; {{ new Date().getFullYear() }} fatwill. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import type { ArticleDetail } from '~/types'
import { apiFetchArticle, apiToggleArticleLike, apiGetArticleLikeStatus, apiRecordArticleView } from '~/utils/api'

const route = useRoute()
const { isDark, toggleTheme } = useTheme()

// SSR 预取文章数据
const { data: articleData } = await useAsyncData(
  `article-${route.params.id}`,
  () => apiFetchArticle(route.params.id as string),
  { default: () => null }
)

const article = ref<ArticleDetail | null>(articleData.value)
const loading = ref(!articleData.value)
const error = ref(false)
const articleContentRef = ref<HTMLElement | null>(null)

// ====== 阅读量 ======
const articleViews = ref(articleData.value?.views ?? 0)

// ====== 阅读进度条 ======
const readProgress = ref(0)

function updateReadProgress() {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  readProgress.value = docHeight > 0 ? Math.min(100, (scrollTop / docHeight) * 100) : 0
}

// ====== TOC ======
interface TocItem {
  id: string
  text: string
  level: number
}

const tocItems = ref<TocItem[]>([])
const activeTocId = ref('')

function buildToc() {
  if (!articleContentRef.value) return
  const headings = articleContentRef.value.querySelectorAll('h1, h2, h3')
  const items: TocItem[] = []
  headings.forEach((el, index) => {
    const id = `heading-${index}`
    el.setAttribute('id', id)
    items.push({
      id,
      text: el.textContent?.trim() || '',
      level: parseInt(el.tagName[1]),
    })
  })
  tocItems.value = items
}

function scrollToHeading(id: string) {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

let tocObserver: IntersectionObserver | null = null

function setupTocObserver() {
  if (tocObserver) tocObserver.disconnect()
  const headings = tocItems.value.map(item => document.getElementById(item.id)).filter(Boolean) as HTMLElement[]
  if (headings.length === 0) return

  tocObserver = new IntersectionObserver(
    (entries) => {
      // 找到最近进入视口的标题
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeTocId.value = entry.target.id
          break
        }
      }
    },
    { rootMargin: '-80px 0px -70% 0px', threshold: 0 }
  )

  headings.forEach(el => tocObserver!.observe(el))
}

// ====== 代码块增强 ======
async function enhanceCodeBlocks() {
  if (!articleContentRef.value) return
  const hljs = (await import('highlight.js')).default
  await import('highlight.js/styles/github-dark.css')

  const codeBlocks = articleContentRef.value.querySelectorAll('pre code')
  codeBlocks.forEach((block) => {
    hljs.highlightElement(block as HTMLElement)

    const pre = block.parentElement
    if (!pre) return

    // 让 pre 相对定位以放置按钮
    pre.style.position = 'relative'

    // 检测语言
    const langClass = Array.from(block.classList).find(c => c.startsWith('language-'))
    const lang = langClass ? langClass.replace('language-', '') : ''

    // 语言标签
    if (lang) {
      const langLabel = document.createElement('span')
      langLabel.className = 'code-lang-label'
      langLabel.textContent = lang
      pre.appendChild(langLabel)
    }

    // 复制按钮
    const copyBtn = document.createElement('button')
    copyBtn.className = 'code-copy-btn'
    copyBtn.textContent = '复制'
    copyBtn.addEventListener('click', async () => {
      try {
        await navigator.clipboard.writeText(block.textContent || '')
        copyBtn.textContent = '已复制 ✓'
        setTimeout(() => { copyBtn.textContent = '复制' }, 2000)
      } catch {
        copyBtn.textContent = '失败'
        setTimeout(() => { copyBtn.textContent = '复制' }, 2000)
      }
    })
    pre.appendChild(copyBtn)
  })
}

// SEO meta（SSR 阶段即生效）
useSeoMeta({
  title: () => articleData.value ? `${articleData.value.title} - fatwill` : 'fatwill - 个人博客',
  description: () => articleData.value?.summary || articleData.value?.title || 'fatwill 的个人博客',
  ogTitle: () => articleData.value?.title || 'fatwill - 个人博客',
  ogDescription: () => articleData.value?.summary || articleData.value?.title || '',
  ogType: 'article',
  ogUrl: () => `https://fatwill.cloud/articles/${route.params.id}`,
  ogSiteName: 'fatwill',
  twitterCard: 'summary',
  twitterTitle: () => articleData.value?.title || 'fatwill - 个人博客',
  twitterDescription: () => articleData.value?.summary || '',
})

// canonical + JSON-LD 结构化数据
useHead({
  link: [
    { rel: 'canonical', href: () => `https://fatwill.cloud/articles/${route.params.id}` },
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: () => articleData.value ? JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'BlogPosting',
        headline: articleData.value.title,
        description: articleData.value.summary || articleData.value.title,
        author: {
          '@type': 'Person',
          name: 'fatwill',
          url: 'https://fatwill.cloud',
        },
        publisher: {
          '@type': 'Person',
          name: 'fatwill',
          url: 'https://fatwill.cloud',
        },
        datePublished: articleData.value.createdAt,
        dateModified: articleData.value.updatedAt || articleData.value.createdAt,
        url: `https://fatwill.cloud/articles/${route.params.id}`,
        mainEntityOfPage: `https://fatwill.cloud/articles/${route.params.id}`,
      }) : '{}',
    },
  ],
})

// ====== 设备唯一 ID ======
const deviceId = ref('')
function getOrCreateDeviceId(): string {
  if (import.meta.server) return ''
  const KEY = 'blog_device_id'
  let id = localStorage.getItem(KEY)
  if (!id) {
    id = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0
      const v = c === 'x' ? r : (r & 0x3) | 0x8
      return v.toString(16)
    })
    localStorage.setItem(KEY, id)
  }
  return id
}

// ====== 记录阅读量 ======
async function recordView() {
  if (!deviceId.value || !article.value) return
  try {
    const res = await apiRecordArticleView(article.value.id, deviceId.value)
    articleViews.value = res.views
  } catch {
    // 静默失败
  }
}

// ====== 文章点赞 ======
const articleLiked = ref(false)
const articleLikeCount = ref(articleData.value?.likeCount ?? 0)
const likeAnimating = ref(false)

async function fetchLikeStatus() {
  if (!deviceId.value || !article.value) return
  try {
    const res = await apiGetArticleLikeStatus(article.value.id, deviceId.value)
    articleLiked.value = res.liked
    articleLikeCount.value = res.likeCount
  } catch {
    articleLikeCount.value = article.value?.likeCount ?? 0
  }
}

async function handleArticleLike() {
  if (!deviceId.value || !article.value) return

  // 乐观更新
  const wasLiked = articleLiked.value
  articleLiked.value = !wasLiked
  articleLikeCount.value = wasLiked ? Math.max(0, articleLikeCount.value - 1) : articleLikeCount.value + 1

  // 弹跳动画
  likeAnimating.value = true
  setTimeout(() => { likeAnimating.value = false }, 600)

  try {
    const res = await apiToggleArticleLike(article.value.id, deviceId.value)
    articleLiked.value = res.liked
    articleLikeCount.value = res.likeCount
  } catch {
    // 回滚
    articleLiked.value = wasLiked
    articleLikeCount.value = wasLiked ? articleLikeCount.value + 1 : Math.max(0, articleLikeCount.value - 1)
  }
}

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
    articleLikeCount.value = data.likeCount ?? 0
    articleViews.value = data.views ?? 0
    if (editor.value && data.content) {
      editor.value.commands.setContent(data.content)
    }
    // 客户端拉取点赞状态 + 记录阅读量
    nextTick(() => {
      fetchLikeStatus()
      recordView()
      enhanceCodeBlocks()
      buildToc()
      nextTick(() => setupTocObserver())
    })
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
  deviceId.value = getOrCreateDeviceId()
  window.addEventListener('scroll', updateReadProgress, { passive: true })

  if (article.value) {
    // SSR 已预取，直接初始化 editor
    if (editor.value && article.value.content) {
      editor.value.commands.setContent(article.value.content)
    }
    loading.value = false
    nextTick(() => {
      fetchLikeStatus()
      recordView()
      enhanceCodeBlocks()
      buildToc()
      nextTick(() => setupTocObserver())
    })
  } else {
    // SSR 预取失败，客户端兜底
    loadArticle()
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
  tocObserver?.disconnect()
  window.removeEventListener('scroll', updateReadProgress)
})
</script>

<style>
/* 代码块增强样式（全局，因为是 DOM 操作插入） */
pre {
  overflow-x: auto;
  border-radius: 0.5rem;
}

.code-lang-label {
  position: absolute;
  top: 0.5rem;
  left: 0.75rem;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.4);
  pointer-events: none;
  line-height: 1;
}

.code-copy-btn {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  padding: 0.2rem 0.6rem;
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  transition: all 0.15s;
  line-height: 1.4;
}

.code-copy-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.2);
}
</style>
