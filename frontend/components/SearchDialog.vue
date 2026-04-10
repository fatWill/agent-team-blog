<template>
  <Teleport to="body">
    <Transition name="search-overlay">
      <div
        v-if="visible"
        class="fixed inset-0 z-[200] flex items-start justify-center bg-black/50 pt-[15vh]"
        @click.self="close"
        @keydown.esc="close"
      >
        <div
          class="w-full max-w-lg rounded-xl border border-gray-200 bg-white shadow-2xl dark:border-gray-700 dark:bg-gray-800 mx-4"
          @keydown.esc="close"
        >
          <!-- 搜索输入框 -->
          <div class="flex items-center gap-3 border-b border-gray-100 px-4 py-3 dark:border-gray-700">
            <svg class="h-5 w-5 shrink-0 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
            <input
              ref="inputRef"
              v-model="keyword"
              type="text"
              placeholder="搜索文章..."
              class="flex-1 bg-transparent text-sm text-gray-900 outline-none placeholder:text-gray-400 dark:text-gray-100 dark:placeholder:text-gray-500"
              @keydown.down.prevent="moveDown"
              @keydown.up.prevent="moveUp"
              @keydown.enter.prevent="confirmSelect"
            />
            <kbd class="hidden rounded border border-gray-200 px-1.5 py-0.5 text-[10px] font-medium text-gray-400 dark:border-gray-600 dark:text-gray-500 sm:inline-block">ESC</kbd>
          </div>

          <!-- 搜索结果 -->
          <div class="max-h-[50vh] overflow-y-auto px-2 py-2">
            <!-- 加载中 -->
            <div v-if="loading" class="flex items-center justify-center py-8">
              <svg class="h-5 w-5 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
            </div>

            <!-- 结果列表 -->
            <template v-else-if="results.length > 0">
              <button
                v-for="(article, index) in results"
                :key="article.id"
                class="flex w-full flex-col rounded-lg px-3 py-2.5 text-left transition-colors"
                :class="selectedIndex === index
                  ? 'bg-primary-50 dark:bg-primary-900/20'
                  : 'hover:bg-gray-50 dark:hover:bg-gray-700/50'"
                @click="goToArticle(article.id)"
                @mouseenter="selectedIndex = index"
              >
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100 line-clamp-1" v-html="highlightKeyword(article.title, keyword)" />
                <span v-if="article.summary" class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 line-clamp-1">{{ article.summary }}</span>
              </button>
            </template>

            <!-- 无结果 -->
            <div v-else-if="keyword.trim() && !loading" class="py-8 text-center">
              <p class="text-sm text-gray-400 dark:text-gray-500">没有找到相关文章</p>
            </div>

            <!-- 空状态 -->
            <div v-else class="py-8 text-center">
              <p class="text-sm text-gray-400 dark:text-gray-500">输入关键词搜索文章</p>
            </div>
          </div>

          <!-- 底部快捷键提示 -->
          <div class="flex items-center gap-4 border-t border-gray-100 px-4 py-2 dark:border-gray-700">
            <span class="flex items-center gap-1 text-[10px] text-gray-400 dark:text-gray-500">
              <kbd class="rounded border border-gray-200 px-1 py-0.5 dark:border-gray-600">↑↓</kbd> 选择
            </span>
            <span class="flex items-center gap-1 text-[10px] text-gray-400 dark:text-gray-500">
              <kbd class="rounded border border-gray-200 px-1 py-0.5 dark:border-gray-600">↵</kbd> 跳转
            </span>
            <span class="flex items-center gap-1 text-[10px] text-gray-400 dark:text-gray-500">
              <kbd class="rounded border border-gray-200 px-1 py-0.5 dark:border-gray-600">esc</kbd> 关闭
            </span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { ArticleListItem } from '~/types'
import { apiSearchArticles } from '~/utils/api'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
}>()

const router = useRouter()
const inputRef = ref<HTMLInputElement | null>(null)
const keyword = ref('')
const results = ref<ArticleListItem[]>([])
const loading = ref(false)
const selectedIndex = ref(0)

let searchTimer: ReturnType<typeof setTimeout> | null = null

function close() {
  emit('update:visible', false)
  keyword.value = ''
  results.value = []
  selectedIndex.value = 0
}

function goToArticle(id: string) {
  close()
  router.push(`/articles/${id}`)
}

function moveDown() {
  if (results.value.length > 0) {
    selectedIndex.value = (selectedIndex.value + 1) % results.value.length
  }
}

function moveUp() {
  if (results.value.length > 0) {
    selectedIndex.value = (selectedIndex.value - 1 + results.value.length) % results.value.length
  }
}

function confirmSelect() {
  if (results.value.length > 0 && selectedIndex.value >= 0) {
    goToArticle(results.value[selectedIndex.value].id)
  }
}

/** 高亮搜索关键词 */
function highlightKeyword(text: string, kw: string): string {
  if (!kw.trim()) return escapeHtml(text)
  const escaped = kw.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(`(${escaped})`, 'gi')
  return escapeHtml(text).replace(regex, '<mark class="search-highlight">$1</mark>')
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

// Debounce 搜索
watch(keyword, (val) => {
  if (searchTimer) clearTimeout(searchTimer)
  selectedIndex.value = 0
  if (!val.trim()) {
    results.value = []
    loading.value = false
    return
  }
  loading.value = true
  searchTimer = setTimeout(async () => {
    try {
      const res = await apiSearchArticles(val.trim())
      results.value = res.list
    } catch {
      results.value = []
    } finally {
      loading.value = false
    }
  }, 300)
})

// 打开时自动聚焦
watch(() => props.visible, (val) => {
  if (val) {
    nextTick(() => inputRef.value?.focus())
  }
})

onUnmounted(() => {
  if (searchTimer) clearTimeout(searchTimer)
})
</script>

<style>
.search-overlay-enter-active {
  transition: opacity 0.2s ease;
}
.search-overlay-leave-active {
  transition: opacity 0.15s ease;
}
.search-overlay-enter-from,
.search-overlay-leave-to {
  opacity: 0;
}

.search-highlight {
  background-color: rgba(234, 179, 8, 0.3);
  color: inherit;
  padding: 0 1px;
  border-radius: 2px;
}

.dark .search-highlight {
  background-color: rgba(234, 179, 8, 0.25);
}
</style>
