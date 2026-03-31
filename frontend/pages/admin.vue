<template>
  <div class="min-h-screen bg-gray-50 transition-colors dark:bg-gray-900">
    <!-- 顶部导航 -->
    <header class="sticky top-0 z-50 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg transition-colors duration-300 dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-4xl items-center justify-between px-4">
        <div class="flex items-center gap-3">
          <NuxtLink to="/" class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            ← 首页
          </NuxtLink>
          <span class="text-gray-300 dark:text-gray-600">/</span>
          <span class="text-sm font-medium text-gray-900 dark:text-gray-100">管理后台</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
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
          <button
            class="rounded-lg px-3 py-1.5 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
            @click="handleLogout"
          >
            退出登录
          </button>
        </div>
      </div>
    </header>

    <!-- Tab 切换 -->
    <nav class="border-b border-gray-200/60 bg-white/80 backdrop-blur-lg dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex max-w-4xl gap-0 px-4">
        <button
          v-for="tab in adminTabs"
          :key="tab.key"
          class="relative px-5 py-3 text-sm font-medium transition-colors"
          :class="activeTab === tab.key
            ? 'text-gray-900 dark:text-gray-100'
            : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
          <span
            v-if="activeTab === tab.key"
            class="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary-500"
          />
        </button>
      </div>
    </nav>

    <!-- 内容区域 -->
    <main class="mx-auto max-w-4xl px-4 py-8">

      <!-- ========== Tab 1: 写文章 / 编辑文章 ========== -->
      <div v-if="activeTab === 'write'">
        <div class="mb-6 flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">
            {{ editingArticleId ? '编辑文章' : '写文章' }}
          </h2>
          <button
            v-if="editingArticleId"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
            @click="cancelEdit"
          >
            取消编辑
          </button>
        </div>

        <form class="space-y-5" @submit.prevent="handlePublish">
          <!-- 标题 -->
          <div>
            <input
              v-model="form.title"
              type="text"
              placeholder="文章标题"
              required
              class="w-full rounded-lg border border-gray-200 bg-white px-4 py-3 text-lg font-medium text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
            />
          </div>

          <!-- 封面图 URL -->
          <div>
            <input
              v-model="form.coverImage"
              type="text"
              placeholder="封面图 URL（可选）"
              class="w-full rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
            />
            <!-- 封面图预览 -->
            <div v-if="form.coverImage" class="mt-2">
              <img
                :src="form.coverImage"
                alt="封面图预览"
                class="h-32 w-auto rounded-lg border border-gray-200 object-cover dark:border-gray-600"
                @error="(e: Event) => (e.target as HTMLImageElement).style.display = 'none'"
              />
            </div>
          </div>

          <!-- 摘要 -->
          <div>
            <textarea
              v-model="form.summary"
              placeholder="文章摘要（可选）"
              rows="2"
              class="w-full resize-none rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
            />
          </div>

          <!-- Tiptap 编辑器 -->
          <div class="tiptap-editor overflow-hidden rounded-lg border border-gray-200 bg-white transition-colors dark:border-gray-600 dark:bg-gray-800">
            <!-- 工具栏 -->
            <div v-if="editor" class="flex flex-wrap gap-1 border-b border-gray-200 px-3 py-2 dark:border-gray-600">
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('bold') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleBold().run()"
              >
                <strong>B</strong>
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('italic') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleItalic().run()"
              >
                <em>I</em>
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('strike') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleStrike().run()"
              >
                <s>S</s>
              </button>
              <span class="mx-1 w-px bg-gray-200 dark:bg-gray-600" />
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('heading', { level: 1 }) ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
              >
                H1
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('heading', { level: 2 }) ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
              >
                H2
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('heading', { level: 3 }) ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
              >
                H3
              </button>
              <span class="mx-1 w-px bg-gray-200 dark:bg-gray-600" />
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('bulletList') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleBulletList().run()"
              >
                • 列表
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('orderedList') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleOrderedList().run()"
              >
                1. 列表
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('blockquote') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleBlockquote().run()"
              >
                引用
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm transition-colors"
                :class="editor.isActive('codeBlock') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="editor.chain().focus().toggleCodeBlock().run()"
              >
                代码块
              </button>
              <span class="mx-1 w-px bg-gray-200 dark:bg-gray-600" />
              <button
                type="button"
                class="rounded px-2 py-1 text-sm text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                @click="addImage"
              >
                图片
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                @click="addLink"
              >
                链接
              </button>
              <button
                type="button"
                class="rounded px-2 py-1 text-sm text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                @click="editor.chain().focus().setHorizontalRule().run()"
              >
                分割线
              </button>
            </div>

            <!-- 编辑区域 -->
            <div class="text-gray-900 dark:text-gray-100">
              <EditorContent :editor="editor" />
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-end gap-3">
            <span v-if="publishSuccess" class="text-sm text-green-600 dark:text-green-400">
              ✓ {{ editingArticleId ? '保存成功！' : '发布成功！' }}
            </span>
            <span v-if="publishError" class="text-sm text-red-500">
              {{ publishError }}
            </span>
            <button
              type="submit"
              :disabled="publishing || !form.title.trim()"
              class="rounded-lg bg-primary-500 px-6 py-2.5 text-sm font-medium text-white transition-all duration-200 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-300 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {{ publishing ? (editingArticleId ? '保存中...' : '发布中...') : (editingArticleId ? '保存修改' : '发布文章') }}
            </button>
          </div>
        </form>
      </div>

      <!-- ========== Tab 2: 文章管理 ========== -->
      <div v-if="activeTab === 'manage'">
        <h2 class="mb-6 text-xl font-bold text-gray-900 dark:text-gray-100">文章管理</h2>

        <!-- 搜索框 -->
        <div class="mb-6">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="按标题搜索文章..."
            class="w-full rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
          />
        </div>

        <!-- 加载状态 -->
        <div v-if="manageLoading" class="space-y-3">
          <div v-for="i in 3" :key="i" class="animate-pulse rounded-lg border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <div class="flex items-center gap-4">
              <div class="h-12 w-20 rounded bg-gray-200 dark:bg-gray-700" />
              <div class="flex-1">
                <div class="mb-2 h-4 w-1/2 rounded bg-gray-200 dark:bg-gray-700" />
                <div class="h-3 w-1/4 rounded bg-gray-100 dark:bg-gray-700/50" />
              </div>
            </div>
          </div>
        </div>

        <!-- 文章列表 -->
        <div v-else-if="manageArticles.length > 0" class="space-y-3">
          <div
            v-for="article in manageArticles"
            :key="article.id"
            class="rounded-lg border border-gray-100 bg-white p-4 transition-colors dark:border-gray-700 dark:bg-gray-800"
          >
            <div class="flex items-start gap-4">
              <!-- 封面缩略图 -->
              <img
                v-if="article.coverImage"
                :src="article.coverImage"
                :alt="article.title"
                class="h-16 w-24 flex-shrink-0 rounded-md object-cover"
              />
              <div v-else class="flex h-16 w-24 flex-shrink-0 items-center justify-center rounded-md bg-gray-100 dark:bg-gray-700">
                <svg class="h-6 w-6 text-gray-300 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3.75 21h16.5a1.5 1.5 0 001.5-1.5V5.25a1.5 1.5 0 00-1.5-1.5H3.75a1.5 1.5 0 00-1.5 1.5v14.25a1.5 1.5 0 001.5 1.5z" />
                </svg>
              </div>

              <!-- 文章信息 -->
              <div class="min-w-0 flex-1">
                <h3 class="truncate text-sm font-semibold text-gray-900 dark:text-gray-100">
                  {{ article.title }}
                </h3>
                <p v-if="article.summary" class="mt-1 truncate text-xs text-gray-500 dark:text-gray-400">
                  {{ article.summary }}
                </p>
                <time class="mt-1 block text-xs text-gray-400 dark:text-gray-500">
                  {{ formatDate(article.createdAt) }}
                </time>
              </div>

              <!-- 操作按钮 -->
              <div class="flex flex-shrink-0 items-center gap-2">
                <button
                  class="rounded-md px-3 py-1.5 text-xs font-medium text-primary-500 transition-colors hover:bg-primary-50 dark:hover:bg-primary-900/20"
                  @click="startEdit(article)"
                >
                  编辑
                </button>
                <button
                  class="rounded-md px-3 py-1.5 text-xs font-medium text-red-500 transition-colors hover:bg-red-50 dark:hover:bg-red-900/20"
                  :disabled="deletingId === article.id"
                  @click="handleDelete(article)"
                >
                  {{ deletingId === article.id ? '删除中...' : '删除' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="py-20 text-center">
          <p class="text-gray-400 dark:text-gray-500">
            {{ searchKeyword ? '没有找到匹配的文章' : '暂无文章' }}
          </p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import type { ArticleListItem } from '~/types'
import { useAuthStore } from '~/stores/auth'
import {
  apiCreateArticle,
  apiUpdateArticle,
  apiDeleteArticle,
  apiGetArticles,
  apiFetchArticle,
} from '~/utils/api'

const authStore = useAuthStore()
const router = useRouter()
const { isDark, toggleTheme } = useTheme()

// ====== Tab 管理 ======
const adminTabs = [
  { key: 'write', label: '写文章' },
  { key: 'manage', label: '文章管理' },
]
const activeTab = ref('write')

// ====== 写文章 / 编辑文章 ======
const editingArticleId = ref<string | null>(null)

const form = reactive({
  title: '',
  summary: '',
  coverImage: '',
})

const publishing = ref(false)
const publishSuccess = ref(false)
const publishError = ref('')

// Tiptap 编辑器
const editor = useEditor({
  extensions: [
    StarterKit,
    Image,
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        target: '_blank',
        rel: 'noopener noreferrer',
      },
    }),
  ],
  content: '<p>开始写作...</p>',
  editorProps: {
    attributes: {
      class: 'prose prose-sm max-w-none',
    },
  },
})

function addImage() {
  const url = window.prompt('请输入图片 URL')
  if (url && editor.value) {
    editor.value.chain().focus().setImage({ src: url }).run()
  }
}

function addLink() {
  const url = window.prompt('请输入链接 URL')
  if (url && editor.value) {
    editor.value.chain().focus().setLink({ href: url }).run()
  }
}

/** 重置表单为新建模式 */
function resetForm() {
  editingArticleId.value = null
  form.title = ''
  form.summary = ''
  form.coverImage = ''
  editor.value?.commands.setContent('<p>开始写作...</p>')
}

/** 取消编辑 */
function cancelEdit() {
  resetForm()
}

/** 发布 / 保存 */
async function handlePublish() {
  if (!editor.value || !form.title.trim()) return

  publishing.value = true
  publishSuccess.value = false
  publishError.value = ''

  try {
    const content = editor.value.getJSON()
    const payload = {
      title: form.title.trim(),
      summary: form.summary.trim() || undefined,
      coverImage: form.coverImage.trim() || undefined,
      content,
    }

    if (editingArticleId.value) {
      // 编辑模式：调用 PUT
      await apiUpdateArticle(editingArticleId.value, payload)
      publishSuccess.value = true
      // 保存成功后切回文章管理 Tab
      setTimeout(() => {
        resetForm()
        activeTab.value = 'manage'
        publishSuccess.value = false
        fetchManageArticles()
      }, 1500)
    } else {
      // 新建模式：调用 POST
      await apiCreateArticle(payload)
      publishSuccess.value = true
      resetForm()
      setTimeout(() => {
        publishSuccess.value = false
      }, 3000)
    }
  } catch (err: unknown) {
    const fetchErr = err as { data?: { statusMessage?: string }; statusCode?: number }
    if (fetchErr?.statusCode === 401 || (fetchErr?.data as Record<string, unknown>)?.statusCode === 401) {
      publishError.value = '登录已过期，请重新登录'
      authStore.setLoggedIn(false)
      setTimeout(() => router.push('/login'), 1500)
    } else {
      publishError.value = editingArticleId.value ? '保存失败，请重试' : '发布失败，请重试'
    }
  } finally {
    publishing.value = false
  }
}

// ====== 文章管理 ======
const searchKeyword = ref('')
const manageArticles = ref<ArticleListItem[]>([])
const manageLoading = ref(false)
const deletingId = ref<string | null>(null)

let searchTimer: ReturnType<typeof setTimeout> | null = null

/** 获取文章列表 */
async function fetchManageArticles() {
  manageLoading.value = true
  try {
    const params = searchKeyword.value.trim() ? { title: searchKeyword.value.trim() } : undefined
    const res = await apiGetArticles(params)
    manageArticles.value = res.list
  } catch {
    manageArticles.value = []
  } finally {
    manageLoading.value = false
  }
}

/** 搜索防抖 */
watch(searchKeyword, () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    fetchManageArticles()
  }, 300)
})

/** 切到文章管理 Tab 时加载数据 */
watch(activeTab, (val) => {
  if (val === 'manage') {
    fetchManageArticles()
  }
})

/** 点击编辑：加载文章详情，填充表单，切到写文章 Tab */
async function startEdit(article: ArticleListItem) {
  try {
    // 获取完整文章详情（含 content）
    const detail = await apiFetchArticle(article.id)
    editingArticleId.value = detail.id
    form.title = detail.title
    form.summary = detail.summary || ''
    form.coverImage = detail.coverImage || ''
    editor.value?.commands.setContent(detail.content)
    activeTab.value = 'write'
  } catch {
    alert('加载文章详情失败，请重试')
  }
}

/** 删除文章 */
async function handleDelete(article: ArticleListItem) {
  if (!confirm(`确定要删除「${article.title}」吗？此操作不可撤销。`)) return

  deletingId.value = article.id
  try {
    await apiDeleteArticle(article.id)
    // 刷新列表
    await fetchManageArticles()
  } catch (err: unknown) {
    const fetchErr = err as { statusCode?: number }
    if (fetchErr?.statusCode === 401) {
      alert('登录已过期，请重新登录')
      authStore.setLoggedIn(false)
      router.push('/login')
    } else {
      alert('删除失败，请重试')
    }
  } finally {
    deletingId.value = null
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

function handleLogout() {
  authStore.logout()
  router.push('/')
}

onUnmounted(() => {
  editor.value?.destroy()
  if (searchTimer) clearTimeout(searchTimer)
})

useHead({
  title: '管理后台 - fatwillzeng',
})
</script>
