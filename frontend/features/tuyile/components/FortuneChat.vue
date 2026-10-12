<template>
  <div class="mx-auto flex h-[100dvh] max-w-3xl flex-col px-4 pb-4 pt-16">
    <div class="mb-3 flex items-center justify-between">
      <div>
        <h2 class="tuyile-serif text-lg font-bold text-gray-900 dark:text-gray-100">🔮 图一乐 · 运势小摊</h2>
        <p class="text-xs text-gray-400 dark:text-gray-500">报上姓名、生日和想看的年份，图一乐给你唠唠</p>
      </div>
      <button
        type="button"
        class="rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-sm text-gray-600 transition-all duration-200 hover:border-rose-300 hover:text-rose-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
        @click="emit('back')"
      >
        ← 回黄历
      </button>
    </div>

    <div class="mb-3 rounded-lg bg-amber-50 px-3 py-2 text-xs leading-relaxed text-amber-700 dark:bg-amber-500/10 dark:text-amber-300">
      🎲 仅供娱乐，不构成任何医疗、投资、法律建议。对话只保存在当前页面内存，刷新或离开即清空，不会保存你的姓名和生日。
    </div>

    <!-- 消息区 -->
    <div
      ref="scrollEl"
      class="flex-1 space-y-4 overflow-y-auto rounded-xl border border-gray-100 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-4"
      aria-live="polite"
      @scroll="onScroll"
    >
      <div v-if="!bubbles.length && isStreaming" class="flex items-center gap-2 text-sm text-gray-400">
        <span class="tuyile-avatar">乐</span>
        <span class="tuyile-dots"><i /><i /><i /></span>
      </div>

      <div v-if="!bubbles.length && !isStreaming && notice" class="py-6 text-center">
        <button type="button" class="rounded-lg bg-rose-600 px-4 py-2 text-sm text-white hover:bg-rose-700 disabled:opacity-50" :disabled="cooldown > 0" @click="retryGreet">
          {{ cooldown > 0 ? `${cooldown}s 后可重试` : '再叫一下图一乐' }}
        </button>
      </div>

      <div
        v-for="b in bubbles"
        :key="b.id"
        class="flex gap-2"
        :class="b.role === 'user' ? 'flex-row-reverse' : ''"
      >
        <span v-if="b.role === 'assistant'" class="tuyile-avatar">乐</span>
        <div
          class="max-w-[85%] rounded-2xl px-3.5 py-2.5 text-sm leading-relaxed"
          :class="b.role === 'user'
            ? 'whitespace-pre-wrap break-words rounded-tr-sm bg-rose-600 text-white'
            : 'tuyile-md rounded-tl-sm bg-gray-50 text-gray-800 dark:bg-gray-900/50 dark:text-gray-100'"
        >
          <template v-if="b.role === 'user'">{{ b.content }}</template>
          <!-- markdown-it 已禁用原始 HTML，可安全 v-html -->
          <div v-else v-html="renderMarkdown(b.content)" />
          <span v-if="b.streaming" class="tuyile-caret" />
        </div>
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

    <!-- 输入区 -->
    <form class="mt-3 flex items-end gap-2" @submit.prevent="handleSend">
      <div class="relative flex-1">
        <textarea
          ref="inputEl"
          v-model="input"
          rows="1"
          :maxlength="TUYILE_MAX_MESSAGE_LENGTH"
          class="block max-h-36 min-h-[44px] w-full resize-none rounded-xl border border-gray-200 bg-white px-3.5 py-2.5 pr-14 text-sm leading-relaxed text-gray-800 outline-none transition-colors focus:border-rose-400 focus:ring-2 focus:ring-rose-400/20 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100"
          placeholder="报上姓名、生日和想看的年份～"
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
  </div>
</template>

<script setup lang="ts">
import { TUYILE_MAX_MESSAGE_LENGTH } from '../api'
import { useTuyileChat } from '../composables/useTuyileChat'
import { renderMarkdown } from '../utils'

const emit = defineEmits<{ (e: 'back'): void }>()

const { bubbles, isStreaming, notice, cooldown, greet, retryGreet, send, stop, dismissNotice } = useTuyileChat()

const input = ref('')
const scrollEl = ref<HTMLElement | null>(null)
const inputEl = ref<HTMLTextAreaElement | null>(null)
const stickBottom = ref(true)

const inputLength = computed(() => Array.from(input.value).length)
const canSend = computed(() => !isStreaming.value && cooldown.value === 0 && input.value.trim() !== '' && inputLength.value <= TUYILE_MAX_MESSAGE_LENGTH)

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

watch(() => bubbles.value.map(b => b.content.length).join(','), () => scrollToBottom())

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

async function handleSend() {
  if (!canSend.value) return
  const text = input.value
  input.value = ''
  nextTick(autoGrow)
  stickBottom.value = true
  scrollToBottom(true)
  const rollback = await send(text)
  if (rollback !== null && input.value === '') {
    input.value = rollback
    nextTick(autoGrow)
  }
}

function formatWait(sec: number): string {
  if (sec >= 3600) return `约 ${Math.ceil(sec / 3600)} 小时`
  if (sec >= 60) return `约 ${Math.ceil(sec / 60)} 分钟`
  return `${sec} 秒`
}

onMounted(() => {
  greet()
})
</script>
