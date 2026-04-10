<template>
  <div class="min-h-screen">
    <!-- 顶部阅读进度条 -->
    <div class="fixed left-0 top-0 z-[60] h-0.5 w-full">
      <div
        class="h-full bg-primary-500 transition-[width] duration-100"
        :style="{ width: readProgress + '%' }"
      />
    </div>

    <!-- 右上角：主题切换按钮（滚动隐藏） -->
    <div
      class="fixed right-4 top-4 z-50 transition-all duration-300"
      :style="{
        transform: headerVisible ? 'translateY(0)' : 'translateY(-200%)',
        opacity: headerVisible ? 1 : 0,
      }"
    >
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

    <!-- 左侧面板：返回按钮 + TOC 目录（PC 端，横向统一面板） -->
    <aside
      class="fixed left-0 top-0 z-40 hidden h-full md:block"
      :style="{
        transform: tocPanelVisible ? 'translateX(0)' : 'translateX(-100%)',
        opacity: tocPanelVisible ? 1 : 0,
        transition: 'transform 0.3s ease, opacity 0.3s ease',
      }"
    >
      <!-- 返回图标 + TOC 横向排列（始终渲染，返回箭头不受 tocExpanded 影响） -->
      <div class="flex h-full flex-row pt-5 pb-8 pl-3">
        <!-- 左列：返回图标按钮（始终可见） -->
        <div class="flex flex-shrink-0 flex-col items-center">
          <NuxtLink
            to="/home"
            class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-700 dark:hover:text-gray-300"
            aria-label="返回首页"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
            </svg>
          </NuxtLink>
          <!-- 收起状态下的展开按钮（紧贴返回箭头下方） -->
          <button
            v-if="tocItems.length > 0 && !tocExpanded"
            class="mt-2 flex h-8 w-8 items-center justify-center rounded-lg text-gray-400 transition-all duration-200 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-300"
            aria-label="展开目录"
            @click="tocExpanded = true"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h10M4 18h16" />
            </svg>
          </button>
        </div>

        <!-- 右列：TOC 目录区域（展开/收起独立控制） -->
        <nav
          v-if="tocItems.length > 0 && tocExpanded"
          class="flex min-h-0 min-w-0 flex-col pr-2"
          style="width: 190px"
        >
          <div class="mb-3 flex items-center justify-between pr-1">
            <h4 class="text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">目录</h4>
            <button
              class="flex h-6 w-6 items-center justify-center rounded text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-300"
              aria-label="收起目录"
              @click="tocExpanded = false"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11 19l-7-7 7-7" />
              </svg>
            </button>
          </div>
          <ul class="flex-1 space-y-1 overflow-y-auto border-l border-gray-200 overscroll-contain dark:border-gray-700">
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
      </div>
    </aside>

    <!-- 移动端：左上角返回 + 目录工具栏（滚动隐藏） -->
    <div
      class="fixed left-4 top-4 z-50 flex items-center gap-1 rounded-lg bg-white/80 shadow-sm backdrop-blur-lg dark:bg-gray-800/80 md:hidden transition-all duration-300"
      :style="{
        transform: headerVisible ? 'translateY(0)' : 'translateY(-200%)',
        opacity: headerVisible ? 1 : 0,
      }"
    >
      <!-- 返回按钮 -->
      <NuxtLink
        to="/home"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
        aria-label="返回首页"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
      </NuxtLink>
      <!-- 分割线 + 目录按钮（仅有目录时显示） -->
      <template v-if="tocItems.length > 0">
        <div class="h-4 w-px bg-gray-200 dark:bg-gray-600" />
        <button
          class="flex h-9 items-center gap-1 rounded-lg px-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
          aria-label="打开目录"
          @click="mobileTocOpen = true"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h10M4 18h16" />
          </svg>
          <span class="text-xs">目录</span>
        </button>
      </template>
    </div>

    <!-- 文章内容区域（居中布局） -->
    <div class="mx-auto max-w-3xl px-4 pt-16 pb-10">
      <!-- 文章主体 -->
      <main class="min-w-0">
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
    </div>

    <!-- 移动端：底部抽屉式目录 -->
    <Teleport to="body">
      <Transition name="mobile-toc-overlay">
        <div
          v-if="mobileTocOpen"
          class="fixed inset-0 z-[70] bg-black/40 backdrop-blur-sm md:hidden"
          @click="mobileTocOpen = false"
        />
      </Transition>
      <Transition name="mobile-toc-drawer">
        <div
          v-if="mobileTocOpen"
          class="fixed bottom-0 left-0 right-0 z-[71] max-h-[60vh] overflow-y-auto rounded-t-2xl bg-white pb-safe dark:bg-gray-800 md:hidden"
        >
          <!-- 抽屉把手 -->
          <div class="sticky top-0 flex items-center justify-between bg-white/95 px-5 py-4 backdrop-blur dark:bg-gray-800/95">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">目录</h4>
            <button
              class="flex h-7 w-7 items-center justify-center rounded-full text-gray-400 transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
              @click="mobileTocOpen = false"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <!-- 目录列表 -->
          <ul class="space-y-0.5 px-5 pb-6">
            <li v-for="item in tocItems" :key="item.id">
              <button
                class="block w-full truncate rounded-lg py-2.5 text-left text-sm transition-colors"
                :class="[
                  activeTocId === item.id
                    ? 'bg-primary-50 font-medium text-primary-600 dark:bg-primary-900/20 dark:text-primary-400'
                    : 'text-gray-600 hover:bg-gray-50 dark:text-gray-400 dark:hover:bg-gray-700/50',
                  item.level === 1 ? 'pl-3' : item.level === 2 ? 'pl-6' : 'pl-9',
                ]"
                @click="handleMobileTocClick(item.id)"
              >
                {{ item.text }}
              </button>
            </li>
          </ul>
        </div>
      </Transition>
    </Teleport>

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

// ====== 滚动方向检测 & 顶栏/TOC 可见性 ======
const headerVisible = ref(true)
const tocPanelVisible = ref(true)
const tocExpanded = ref(true)
const mobileTocOpen = ref(false)
let lastScrollY = 0
const SCROLL_THRESHOLD = 10 // 防抖阈值，避免微小滚动触发

function updateReadProgress() {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  readProgress.value = docHeight > 0 ? Math.min(100, (scrollTop / docHeight) * 100) : 0

  // 滚动方向检测
  const delta = scrollTop - lastScrollY
  if (Math.abs(delta) > SCROLL_THRESHOLD) {
    if (delta > 0 && scrollTop > 80) {
      // 向下滚动且超过一定距离：隐藏顶栏和 TOC
      headerVisible.value = false
      tocPanelVisible.value = false
    } else if (delta < 0) {
      // 向上滚动：显示顶栏和 TOC
      headerVisible.value = true
      tocPanelVisible.value = true
    }
    lastScrollY = scrollTop
  }
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

function handleMobileTocClick(id: string) {
  mobileTocOpen.value = false
  // 等抽屉关闭动画结束后再滚动
  setTimeout(() => {
    scrollToHeading(id)
  }, 300)
}

let tocObserver: IntersectionObserver | null = null

function setupTocObserver() {
  if (tocObserver) tocObserver.disconnect()
  const headings = tocItems.value.map(item => document.getElementById(item.id)).filter(Boolean) as HTMLElement[]
  if (headings.length === 0) return

  tocObserver = new IntersectionObserver(
    (entries) => {
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

    pre.style.position = 'relative'

    const langClass = Array.from(block.classList).find(c => c.startsWith('language-'))
    const lang = langClass ? langClass.replace('language-', '') : ''

    if (lang) {
      const langLabel = document.createElement('span')
      langLabel.className = 'code-lang-label'
      langLabel.textContent = lang
      pre.appendChild(langLabel)
    }

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

  const wasLiked = articleLiked.value
  articleLiked.value = !wasLiked
  articleLikeCount.value = wasLiked ? Math.max(0, articleLikeCount.value - 1) : articleLikeCount.value + 1

  likeAnimating.value = true
  setTimeout(() => { likeAnimating.value = false }, 600)

  try {
    const res = await apiToggleArticleLike(article.value.id, deviceId.value)
    articleLiked.value = res.liked
    articleLikeCount.value = res.likeCount
  } catch {
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
    nextTick(() => {
      fetchLikeStatus()
      recordView()
      enhanceCodeBlocks()
      nextTick(() => {
        buildToc()
        nextTick(() => setupTocObserver())
      })
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
  lastScrollY = window.scrollY
  window.addEventListener('scroll', updateReadProgress, { passive: true })

  if (article.value) {
    if (editor.value && article.value.content) {
      editor.value.commands.setContent(article.value.content)
    }
    loading.value = false
    nextTick(() => {
      fetchLikeStatus()
      recordView()
      enhanceCodeBlocks()
      nextTick(() => {
        buildToc()
        nextTick(() => setupTocObserver())
      })
    })
  } else {
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

/* 移动端目录抽屉动画 */
.mobile-toc-overlay-enter-active,
.mobile-toc-overlay-leave-active {
  transition: opacity 0.3s ease;
}

.mobile-toc-overlay-enter-from,
.mobile-toc-overlay-leave-to {
  opacity: 0;
}

.mobile-toc-drawer-enter-active,
.mobile-toc-drawer-leave-active {
  transition: transform 0.3s ease;
}

.mobile-toc-drawer-enter-from,
.mobile-toc-drawer-leave-to {
  transform: translateY(100%);
}

/* iOS 安全区域适配 */
.pb-safe {
  padding-bottom: max(1.5rem, env(safe-area-inset-bottom));
}
</style>
