<template>
  <div class="mx-auto flex h-[100dvh] max-w-3xl flex-col px-4 pb-4 pt-16">
    <div class="mb-3 flex items-center justify-between gap-2">
      <div class="min-w-0">
        <h2 class="tuyile-serif truncate text-lg font-bold text-gray-900 dark:text-gray-100">🔮 图一乐 · 运势小摊</h2>
        <p class="truncate text-xs text-gray-400 dark:text-gray-500">大罗天仙在线摆摊，搬家跳槽取名整活，啥都能算</p>
      </div>
      <div class="flex shrink-0 items-center gap-2">
        <button
          type="button"
          class="flex items-center gap-1 rounded-lg border border-rose-200 bg-white px-3 py-1.5 text-sm text-rose-600 transition-all duration-200 hover:border-rose-400 hover:bg-rose-50 disabled:cursor-not-allowed disabled:opacity-40 dark:border-rose-500/30 dark:bg-gray-800 dark:text-rose-300 dark:hover:bg-gray-700"
          :disabled="!canShare"
          :title="shareTitle"
          @click="shareOpen = true"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z" />
          </svg>
          分享
        </button>
        <button
          type="button"
          class="rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-sm text-gray-600 transition-all duration-200 hover:border-rose-300 hover:text-rose-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
          @click="emit('back')"
        >
          ← 黄历
        </button>
      </div>
    </div>

    <div class="mb-3 rounded-lg bg-amber-50 px-3 py-2 text-xs leading-relaxed text-amber-700 dark:bg-amber-500/10 dark:text-amber-300">
      🎲 仅供娱乐，不构成医疗、投资、法律建议。对话只保存在当前页面，刷新或离开即清空；主动点「分享」才会生成只读快照。
    </div>

    <!-- 消息区 -->
    <div
      ref="scrollEl"
      class="flex-1 space-y-4 overflow-y-auto rounded-xl border border-gray-100 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-4"
      aria-live="polite"
      @scroll="onScroll"
    >
      <div v-if="!bubbles.length && !isStreaming && notice" class="py-6 text-center">
        <button type="button" class="rounded-lg bg-rose-600 px-4 py-2 text-sm text-white hover:bg-rose-700 disabled:opacity-50" :disabled="cooldown > 0" @click="retryGreet">
          {{ cooldown > 0 ? `${cooldown}s 后可重试` : '再叫一下图一乐' }}
        </button>
      </div>

      <ChatMessage
        v-for="b in bubbles"
        :key="b.id"
        :role="b.role"
        :content="b.content"
        :streaming="b.streaming"
      />

      <!-- 思考状态条：收到 TEXT_MESSAGE_START 即消失 -->
      <Transition name="tuyile-fade">
        <div v-if="isStreaming && status" class="flex items-center gap-2" role="status">
          <span class="tuyile-avatar">乐</span>
          <div class="tuyile-status inline-flex max-w-[85%] items-center gap-2 rounded-2xl rounded-tl-sm px-3.5 py-2 text-sm text-rose-700 dark:text-rose-200" :class="isDark ? 'tuyile-status-dark' : ''">
            <span class="tuyile-status-icon" aria-hidden="true">{{ status.stage === 'querying' ? '📖' : '🔮' }}</span>
            <span class="truncate">
              {{ status.text }}<span v-if="status.range" class="ml-1 text-xs text-rose-500/80 dark:text-rose-300/70">({{ formatRange(status.range) }})</span>
            </span>
            <span class="tuyile-dots shrink-0"><i /><i /><i /></span>
          </div>
        </div>
      </Transition>
      <!-- 已发起但尚未收到任何状态事件 -->
      <div v-if="isStreaming && !status && !hasStreamingBubble" class="flex items-center gap-2 text-sm text-gray-400">
        <span class="tuyile-avatar">乐</span>
        <span class="tuyile-dots"><i /><i /><i /></span>
      </div>
    </div>

    <!-- 提示条 -->
    <Transition name="tuyile-fade">
      <div
        v-if="notice"
        class="mt-3 flex items-start justify-between gap-2 rounded-lg px-3 py-2 text-sm"
        :class="notice.kind === 'rate'
          ? 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-300'
          : 'bg-rose-50 text-rose-700 dark:bg-rose-500/10 dark:text-rose-300'"
        role="alert"
      >
        <span>
          {{ notice.message }}
          <template v-if="notice.kind === 'rate' && cooldown > 0">（{{ formatWait(cooldown) }}后再试）</template>
        </span>
        <button
          v-if="!(notice.kind === 'rate' && cooldown > 0)"
          type="button"
          class="shrink-0 opacity-60 hover:opacity-100"
          aria-label="关闭提示"
          @click="dismissNotice"
        >
          ✕
        </button>
      </div>
    </Transition>

    <!-- 快捷提问 -->
    <div class="tuyile-chips -mx-4 mt-3 flex gap-2 overflow-x-auto px-4 pb-0.5" role="group" aria-label="快捷提问">
      <button
        v-for="q in QUICK_QUESTIONS"
        :key="q.text"
        type="button"
        class="shrink-0 whitespace-nowrap rounded-full border border-rose-200 bg-white px-3 py-1.5 text-xs text-rose-600 transition-all duration-200 hover:-translate-y-0.5 hover:border-rose-400 hover:bg-rose-50 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:translate-y-0 dark:border-rose-500/30 dark:bg-gray-800 dark:text-rose-300 dark:hover:bg-gray-700"
        :disabled="isStreaming || cooldown > 0"
        @click="handleQuick(q.text)"
      >
        {{ q.icon }} {{ q.text }}
      </button>
    </div>

    <!-- 输入区 -->
    <form class="mt-2 flex items-end gap-2" @submit.prevent="handleSend">
      <div class="relative flex-1">
        <textarea
          ref="inputEl"
          v-model="input"
          rows="1"
          :maxlength="TUYILE_MAX_MESSAGE_LENGTH"
          class="block max-h-36 min-h-[44px] w-full resize-none rounded-xl border border-gray-200 bg-white px-3.5 py-2.5 pr-14 text-sm leading-relaxed text-gray-800 outline-none transition-colors placeholder:truncate focus:border-rose-400 focus:ring-2 focus:ring-rose-400/20 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100"
          :placeholder="placeholder"
          aria-label="输入消息"
          @input="autoGrow"
          @keydown.enter.exact="handleEnter"
        />
        <span
          class="pointer-events-none absolute bottom-2 right-3 text-[11px]"
          :class="inputLength >= TUYILE_MAX_MESSAGE_LENGTH ? 'text-rose-500' : 'text-gray-400'"
        >
          {{ inputLength }}/{{ TUYILE_MAX_MESSAGE_LENGTH }}
        </span>
      </div>
      <button
        v-if="isStreaming"
        type="button"
        class="h-11 shrink-0 rounded-xl border border-rose-300 bg-white px-4 text-sm font-medium text-rose-600 transition-all duration-200 hover:bg-rose-50 dark:border-rose-500/40 dark:bg-gray-800 dark:text-rose-300"
        @click="stop"
      >
        ■ 停止
      </button>
      <button
        v-else
        type="submit"
        class="h-11 shrink-0 rounded-xl bg-rose-600 px-5 text-sm font-medium text-white transition-all duration-200 hover:bg-rose-700 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="!canSend"
      >
        {{ cooldown > 0 ? `${cooldown}s` : '发送' }}
      </button>
    </form>

    <ShareDialog v-model:open="shareOpen" :bubbles="bubbles" />
  </div>
</template>

<script setup lang="ts">
import { TUYILE_MAX_MESSAGE_LENGTH } from '../api'
import { useTuyileChat } from '../composables/useTuyileChat'
import ChatMessage from './ChatMessage.vue'
import ShareDialog from './ShareDialog.vue'

const emit = defineEmits<{ (e: 'back'): void }>()

const { bubbles, isStreaming, notice, cooldown, status, greet, retryGreet, send, stop, dismissNotice } = useTuyileChat()
const { isDark } = useTheme()

const QUICK_QUESTIONS = [
  { icon: '🏠', text: '下周哪天适合搬家' },
  { icon: '🪪', text: '用我的网名算算性格' },
  { icon: '🚪', text: '明天几点出门最吉利' },
  { icon: '💨', text: '今天往哪放屁最吉利' },
  { icon: '💼', text: '这个月哪天适合提离职' },
]

const PLACEHOLDERS = [
  '报个网名，本仙给你批一卦～',
  '想算搬家、跳槽、表白，还是今天几点放屁最吉利？',
  '取名、测名、择吉日，本仙全都包了',
  '有啥想不开的？说来让本仙乐呵乐呵',
  '问吧，算不准不要钱（本来也不要钱）',
]

const input = ref('')
const scrollEl = ref<HTMLElement | null>(null)
const inputEl = ref<HTMLTextAreaElement | null>(null)
const stickBottom = ref(true)
const shareOpen = ref(false)
const placeholderIndex = ref(0)
let placeholderTimer: ReturnType<typeof setInterval> | null = null

const placeholder = computed(() => PLACEHOLDERS[placeholderIndex.value])
const inputLength = computed(() => Array.from(input.value).length)
const canSend = computed(() => !isStreaming.value && cooldown.value === 0 && input.value.trim() !== '' && inputLength.value <= TUYILE_MAX_MESSAGE_LENGTH)
const hasAssistant = computed(() => bubbles.value.some(b => b.role === 'assistant' && !b.streaming && b.content.trim() !== ''))
const canShare = computed(() => hasAssistant.value && !isStreaming.value)
const shareTitle = computed(() => {
  if (isStreaming.value) return '图一乐还在说，等说完再分享～'
  if (!hasAssistant.value) return '等图一乐回一句再分享～'
  return '分享这段对话（只读链接）'
})
const hasStreamingBubble = computed(() => bubbles.value.some(b => b.streaming))

function onScroll() {
  const el = scrollEl.value
  if (!el) return
  stickBottom.value = el.scrollHeight - el.scrollTop - el.clientHeight < 60
}

function scrollToBottom(force = false) {
  nextTick(() => {
    const el = scrollEl.value
    if (el && (force || stickBottom.value)) el.scrollTop = el.scrollHeight
  })
}

watch(() => `${bubbles.value.map(b => b.content.length).join(',')}|${status.value?.text ?? ''}`, () => scrollToBottom())

function autoGrow() {
  const el = inputEl.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, 144)}px`
}

function handleEnter(e: KeyboardEvent) {
  // 输入法组字中的回车不发送
  if (e.isComposing) return
  e.preventDefault()
  handleSend()
}

async function submit(text: string, fromInput: boolean) {
  stickBottom.value = true
  scrollToBottom(true)
  const rollback = await send(text)
  if (rollback !== null && fromInput && input.value === '') {
    input.value = rollback
    nextTick(autoGrow)
  }
}

async function handleSend() {
  if (!canSend.value) return
  const text = input.value
  input.value = ''
  nextTick(autoGrow)
  await submit(text, true)
}

async function handleQuick(text: string) {
  if (isStreaming.value || cooldown.value > 0) return
  await submit(text, false)
}

/** 2028-02-02~2028-02-16 → 2028-02-02~02-16（同年省略年份） */
function formatRange(range: string): string {
  const [a, b] = range.split('~')
  if (!b) return a
  return a.slice(0, 4) === b.slice(0, 4) ? `${a}~${b.slice(5)}` : range
}

function formatWait(sec: number): string {
  if (sec >= 3600) return `约 ${Math.ceil(sec / 3600)} 小时`
  if (sec >= 60) return `约 ${Math.ceil(sec / 60)} 分钟`
  return `${sec} 秒`
}

onMounted(() => {
  placeholderIndex.value = Math.floor(Math.random() * PLACEHOLDERS.length)
  placeholderTimer = setInterval(() => {
    if (input.value === '') placeholderIndex.value = (placeholderIndex.value + 1) % PLACEHOLDERS.length
  }, 6000)
  greet()
})

onBeforeUnmount(() => {
  if (placeholderTimer) clearInterval(placeholderTimer)
})
</script>

<style scoped>
.tuyile-status {
  background: linear-gradient(90deg, rgba(255, 228, 230, 0.9), rgba(254, 243, 199, 0.9), rgba(255, 228, 230, 0.9));
  background-size: 200% 100%;
  animation: tuyile-status-flow 2.4s linear infinite;
}
.tuyile-status.tuyile-status-dark {
  background: linear-gradient(90deg, rgba(225, 29, 72, 0.22), rgba(245, 158, 11, 0.16), rgba(225, 29, 72, 0.22));
  background-size: 200% 100%;
}
.tuyile-status-icon {
  display: inline-block;
  animation: tuyile-status-spin 2s ease-in-out infinite;
}
.tuyile-chips {
  scrollbar-width: none;
}
.tuyile-chips::-webkit-scrollbar {
  display: none;
}
@keyframes tuyile-status-flow {
  from { background-position: 0% 0; }
  to { background-position: 200% 0; }
}
@keyframes tuyile-status-spin {
  0%, 100% { transform: rotate(0deg) scale(1); }
  50% { transform: rotate(-12deg) scale(1.12); }
}
@media (prefers-reduced-motion: reduce) {
  .tuyile-status,
  .tuyile-status-icon {
    animation: none;
  }
}
</style>
