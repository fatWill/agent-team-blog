<template>
  <!-- 悬浮按钮 -->
  <button
    v-if="!isOpen"
    class="fixed bottom-16 right-6 z-50 flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 text-white shadow-lg transition-all duration-200 hover:scale-110 hover:bg-primary-600 hover:shadow-xl active:scale-95"
    aria-label="打开留言板"
    @click="open"
  >
    <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M12 2.25c5.385 0 9.75 3.806 9.75 8.5 0 4.694-4.365 8.5-9.75 8.5a10.17 10.17 0 01-3.248-.534c-.364-.12-.757-.06-1.064.147l-2.025 1.37a.75.75 0 01-1.163-.633v-1.72a.75.75 0 00-.28-.587C2.528 15.446 2.25 13.24 2.25 10.75c0-4.694 4.365-8.5 9.75-8.5z" />
    </svg>
  </button>

  <!-- 气泡弹窗 -->
  <Teleport to="body">
    <Transition name="bubble">
      <div
        v-if="isOpen"
        class="fixed bottom-16 right-6 z-50 flex w-[360px] max-w-[calc(100vw-48px)] flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-gray-700 dark:bg-gray-800"
        style="max-height: min(500px, calc(100vh - 48px));"
      >
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-gray-100 px-4 py-3 dark:border-gray-700">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">💬 留言板</h3>
          <button
            class="flex h-7 w-7 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-300"
            aria-label="关闭留言板"
            @click="close"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- 留言列表（滚动区） -->
        <div ref="listRef" class="flex-1 overflow-y-auto px-4 py-3" style="min-height: 0;">
          <!-- 加载中 -->
          <div v-if="messagesLoading" class="flex items-center justify-center py-10">
            <AppLoading tip="加载中..." />
          </div>

          <!-- 空状态 -->
          <div v-else-if="messages.length === 0" class="py-10 text-center">
            <span class="mb-2 block text-2xl">💬</span>
            <p class="text-sm text-gray-400 dark:text-gray-500">还没有留言，快来第一个～</p>
          </div>

          <!-- 留言列表 -->
          <div v-else class="space-y-3">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="rounded-lg border border-gray-100 bg-gray-50/50 p-3 transition-colors dark:border-gray-700 dark:bg-gray-750/50"
              :class="{ 'ring-1 ring-primary-200 dark:ring-primary-800': msg.isOwn }"
            >
              <div class="mb-1.5 flex items-center justify-between">
                <div class="flex items-center gap-1.5">
                  <span class="text-xs font-medium text-gray-900 dark:text-gray-100">{{ msg.nickname }}</span>
                  <span
                    v-if="msg.isOwn"
                    class="rounded-full bg-primary-100 px-1.5 py-0.5 text-[9px] font-medium leading-none text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
                  >我</span>
                </div>
                <time class="text-[10px] text-gray-400 dark:text-gray-500">{{ relativeTime(msg.createdAt) }}</time>
              </div>
              <p class="whitespace-pre-wrap text-xs leading-relaxed text-gray-600 dark:text-gray-400">{{ msg.content }}</p>
            </div>
          </div>
        </div>

        <!-- 输入区（条件显示） -->
        <div class="border-t border-gray-100 px-4 py-3 dark:border-gray-700">
          <!-- 已留言提示 -->
          <div v-if="myMessage" class="text-center">
            <p class="text-xs text-gray-400 dark:text-gray-500">你已经留过言了 ✨</p>
          </div>

          <!-- 输入框 -->
          <div v-else class="space-y-2">
            <textarea
              v-model="content"
              maxlength="500"
              rows="3"
              placeholder="写下你想说的话..."
              class="w-full resize-none rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
            />
            <div class="flex items-center justify-between">
              <span class="text-[10px] text-gray-400 dark:text-gray-500">{{ content.length }}/500</span>
              <button
                :disabled="submitting || !content.trim()"
                class="rounded-lg bg-primary-500 px-4 py-1.5 text-xs font-medium text-white transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
                @click="handleSubmit"
              >
                {{ submitting ? '提交中...' : '发布留言' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { MessageItem } from '~/features/guestbook'
import { apiGetMessages, apiCreateMessage } from '~/features/guestbook'

const isOpen = ref(false)
const messages = ref<MessageItem[]>([])
const messagesLoading = ref(false)
const messagesLoaded = ref(false)
const content = ref('')
const submitting = ref(false)
const listRef = ref<HTMLElement | null>(null)

// 设备唯一 ID
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

/** 当前设备已发过的留言 */
const myMessage = computed(() => messages.value.find(m => m.isOwn))

/** 相对时间格式化 */
function relativeTime(dateStr: string): string {
  const now = Date.now()
  const diff = now - new Date(dateStr).getTime()
  const seconds = Math.floor(diff / 1000)
  if (seconds < 60) return '刚刚'
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

async function fetchMessages() {
  messagesLoading.value = true
  try {
    const res = await apiGetMessages(deviceId.value || undefined)
    messages.value = res.list
    messagesLoaded.value = true
  } catch {
    messages.value = []
  } finally {
    messagesLoading.value = false
  }
}

function open() {
  isOpen.value = true
  if (!messagesLoaded.value) {
    fetchMessages()
  }
}

function close() {
  isOpen.value = false
}

async function handleSubmit() {
  if (!content.value.trim() || !deviceId.value) return
  submitting.value = true
  try {
    const created = await apiCreateMessage({
      deviceId: deviceId.value,
      content: content.value.trim(),
    })
    messages.value.unshift(created)
    content.value = ''
    // 滚动到顶部
    nextTick(() => {
      listRef.value?.scrollTo({ top: 0, behavior: 'smooth' })
    })
  } catch (err: unknown) {
    const fetchErr = err as { statusCode?: number; data?: { statusMessage?: string }; statusMessage?: string }
    const errMsg = fetchErr?.data?.statusMessage || fetchErr?.statusMessage || '操作失败'
    if (import.meta.client) {
      alert(errMsg)
    }
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  deviceId.value = getOrCreateDeviceId()
})
</script>

<style scoped>
/* 气泡弹窗动画：从右下角 scale + fade */
.bubble-enter-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.bubble-leave-active {
  transition: all 0.15s ease-in;
}
.bubble-enter-from {
  opacity: 0;
  transform: scale(0.8) translateY(20px);
  transform-origin: bottom right;
}
.bubble-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(10px);
  transform-origin: bottom right;
}
</style>
