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

    <!-- 编辑器主体 -->
    <main class="mx-auto max-w-4xl px-4 py-8">
      <h2 class="mb-6 text-xl font-bold text-gray-900 dark:text-gray-100">写文章</h2>

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
            ✓ 发布成功！
          </span>
          <span v-if="publishError" class="text-sm text-red-500">
            {{ publishError }}
          </span>
          <button
            type="submit"
            :disabled="publishing || !form.title.trim()"
            class="rounded-lg bg-primary-500 px-6 py-2.5 text-sm font-medium text-white transition-all duration-200 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-300 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {{ publishing ? '发布中...' : '发布文章' }}
          </button>
        </div>
      </form>
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
import { useAuthStore } from '~/stores/auth'
import { apiCreateArticle } from '~/utils/api'

const authStore = useAuthStore()
const router = useRouter()
const { isDark, toggleTheme } = useTheme()

const form = reactive({
  title: '',
  summary: '',
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

async function handlePublish() {
  if (!editor.value || !form.title.trim()) return

  publishing.value = true
  publishSuccess.value = false
  publishError.value = ''

  try {
    const content = editor.value.getJSON()
    await apiCreateArticle({
      title: form.title.trim(),
      summary: form.summary.trim() || undefined,
      content,
    })
    publishSuccess.value = true
    // 重置表单
    form.title = ''
    form.summary = ''
    editor.value.commands.setContent('<p>开始写作...</p>')

    // 2 秒后跳转首页
    setTimeout(() => {
      router.push('/')
    }, 2000)
  } catch (err: unknown) {
    const fetchErr = err as { data?: { statusMessage?: string }; statusCode?: number }
    if (fetchErr?.statusCode === 401 || (fetchErr?.data as Record<string, unknown>)?.statusCode === 401) {
      publishError.value = '登录已过期，请重新登录'
      authStore.setLoggedIn(false)
      setTimeout(() => router.push('/login'), 1500)
    } else {
      publishError.value = '发布失败，请重试'
    }
  } finally {
    publishing.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/')
}

onBeforeUnmount(() => {
  editor.value?.destroy()
})

useHead({
  title: '管理后台 - fatwillzeng',
})
</script>
