<template>
  <Teleport to="body">
    <Transition name="tuyile-fade">
      <div
        v-if="open"
        class="tuyile-root fixed inset-0 z-[1000] flex items-end justify-center bg-black/50 p-0 backdrop-blur-sm sm:items-center sm:p-4"
        role="dialog"
        aria-modal="true"
        aria-labelledby="tuyile-share-title"
        @click.self="handleClose"
      >
        <div class="w-full max-w-md rounded-t-2xl bg-white p-5 shadow-xl dark:bg-gray-800 sm:rounded-2xl">
          <!-- 确认阶段 -->
          <template v-if="phase === 'confirm'">
            <h3 id="tuyile-share-title" class="tuyile-serif mb-1 text-lg font-bold text-gray-900 dark:text-gray-100">分享这段对话</h3>
            <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">
              将生成一个<strong class="text-gray-700 dark:text-gray-200">只读</strong>链接，别人打开只能看，不能接着聊。
            </p>
            <div class="mb-3 rounded-xl bg-rose-50/70 p-3 text-sm text-gray-700 dark:bg-gray-900/50 dark:text-gray-200">
              <div class="flex items-center justify-between">
                <span>将分享</span>
                <span class="font-bold text-rose-600 dark:text-rose-400">{{ plan.messages.length }} 条消息</span>
              </div>
              <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                你 {{ userCount }} 条 · 图一乐 {{ assistantCount }} 条
              </div>
            </div>
            <ul v-if="plan.dropped || plan.truncated" class="mb-3 space-y-1 rounded-lg bg-amber-50 px-3 py-2 text-xs leading-relaxed text-amber-700 dark:bg-amber-500/10 dark:text-amber-300">
              <li v-if="plan.dropped">对话共 {{ plan.total }} 条，超出上限，只分享最近 {{ plan.messages.length }} 条。</li>
              <li v-if="plan.truncated">有 {{ plan.truncated }} 条消息超过 {{ SHARE_MAX_CONTENT_LENGTH }} 字，已截断。</li>
            </ul>
            <div class="mb-4 max-h-40 space-y-1.5 overflow-y-auto rounded-lg border border-gray-100 p-2 text-xs dark:border-gray-700">
              <div v-for="(m, i) in previewList" :key="i" class="truncate" :class="m.role === 'user' ? 'text-rose-600 dark:text-rose-300' : 'text-gray-600 dark:text-gray-300'">
                <span class="font-semibold">{{ m.role === 'user' ? '我' : '图一乐' }}：</span>{{ m.preview }}
              </div>
            </div>
            <p class="mb-4 text-xs text-gray-400 dark:text-gray-500">提醒：分享内容将公开可见，请勿包含真实姓名、生日等隐私信息。</p>
            <div v-if="error" class="mb-3 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-500/10 dark:text-rose-300" role="alert">
              {{ error }}<template v-if="cooldown > 0">（{{ cooldown }}s 后可重试）</template>
            </div>
            <div class="flex gap-2">
              <button type="button" class="h-10 flex-1 rounded-xl border border-gray-200 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-700" @click="handleClose">
                取消
              </button>
              <button
                type="button"
                class="h-10 flex-1 rounded-xl bg-rose-600 text-sm font-medium text-white transition-colors hover:bg-rose-700 disabled:cursor-not-allowed disabled:opacity-50"
                :disabled="submitting || cooldown > 0 || !canShare"
                @click="handleSubmit"
              >
                {{ submitting ? '生成中…' : (inWechat ? '生成并设置微信卡片' : '生成分享链接') }}
              </button>
            </div>
            <p v-if="!canShare" class="mt-2 text-center text-xs text-gray-400">至少需要 1 条图一乐的回复才能分享</p>
          </template>

          <!-- 成功阶段 -->
          <template v-else>
            <div class="mb-3 text-center">
              <div class="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-emerald-50 text-2xl dark:bg-emerald-500/10">🎉</div>
              <h3 id="tuyile-share-title" class="tuyile-serif text-lg font-bold text-gray-900 dark:text-gray-100">分享成功</h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ inWechat ? (wechatReady ? '微信分享卡片已设置好～' : '链接已生成，卡片设置可能未生效，可直接复制链接') : '把链接发给朋友，一起图个乐～' }}
              </p>
            </div>
            <div class="mb-4 flex items-center gap-2 rounded-xl border border-gray-200 bg-gray-50 p-2 dark:border-gray-700 dark:bg-gray-900/50">
              <input
                ref="urlInput"
                class="min-w-0 flex-1 bg-transparent px-1 text-sm text-gray-700 outline-none dark:text-gray-200"
                :value="result?.url"
                readonly
                aria-label="分享链接"
                @focus="selectUrl"
              >
            </div>
            <div class="flex gap-2">
              <button type="button" class="h-10 flex-1 rounded-xl bg-rose-600 text-sm font-medium text-white transition-colors hover:bg-rose-700" @click="handleCopy">
                复制链接
              </button>
              <a
                :href="result?.url"
                target="_blank"
                rel="noopener"
                class="flex h-10 flex-1 items-center justify-center rounded-xl border border-gray-200 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-700"
              >
                打开看看
              </a>
            </div>
            <button
              v-if="inWechat"
              type="button"
              class="mt-2 h-10 w-full rounded-xl bg-emerald-600 text-sm font-medium text-white transition-colors hover:bg-emerald-700"
              @click="showGuide = true"
            >
              去微信转发
            </button>
            <button type="button" class="mt-3 w-full text-center text-sm text-gray-400 hover:text-gray-600 dark:hover:text-gray-200" @click="handleClose">
              关闭
            </button>
          </template>
        </div>
      </div>
    </Transition>

    <!-- 微信右上角转发引导蒙层 -->
    <Transition name="tuyile-fade">
      <div v-if="showGuide" class="fixed inset-0 z-[1100] bg-black/75" role="dialog" aria-label="微信转发引导" @click="showGuide = false">
        <div class="absolute right-5 top-2 text-right text-white">
          <svg class="ml-auto h-20 w-20" viewBox="0 0 80 80" fill="none" aria-hidden="true">
            <path d="M14 70 C 30 40, 50 30, 66 12" stroke="white" stroke-width="3" stroke-linecap="round" stroke-dasharray="6 6" />
            <path d="M54 12 L 67 11 L 64 24" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <p class="mt-2 text-lg font-bold">点击右上角「…」</p>
          <p class="mt-1 text-sm text-white/85">选择「发送给朋友」或「分享到朋友圈」</p>
          <p class="mt-6 text-xs text-white/60">轻触任意处关闭</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { ChatBubble, CreateShareResponse } from '../types'
import {
  apiCreateTuyileShare,
  SHARE_MAX_BODY_BYTES,
  SHARE_MAX_CONTENT_LENGTH,
  SHARE_MAX_MESSAGES,
  TuyileHttpError,
} from '../api'
import { isWechatUA, markdownSummary, planShare } from '../utils'
import { setTuyileWechatCard } from '../composables/useTuyileWechat'
import { showError, showSuccess } from '~/utils/ui'

const props = defineProps<{ open: boolean, bubbles: ChatBubble[] }>()
const emit = defineEmits<{ (e: 'update:open', v: boolean): void }>()

const phase = ref<'confirm' | 'done'>('confirm')
const submitting = ref(false)
const error = ref('')
const cooldown = ref(0)
const result = ref<CreateShareResponse | null>(null)
const wechatReady = ref(false)
const showGuide = ref(false)
const urlInput = ref<HTMLInputElement | null>(null)
const inWechat = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

const plan = computed(() => planShare(
  props.bubbles.filter(b => !b.streaming),
  SHARE_MAX_MESSAGES,
  SHARE_MAX_CONTENT_LENGTH,
  SHARE_MAX_BODY_BYTES,
))
const userCount = computed(() => plan.value.messages.filter(m => m.role === 'user').length)
const assistantCount = computed(() => plan.value.messages.filter(m => m.role === 'assistant').length)
const canShare = computed(() => assistantCount.value > 0)
const previewList = computed(() => plan.value.messages.map(m => ({
  role: m.role,
  preview: m.role === 'assistant' ? markdownSummary(m.content, 60) : Array.from(m.content).slice(0, 60).join(''),
})))

watch(() => props.open, (v) => {
  if (v) {
    inWechat.value = isWechatUA()
    phase.value = 'confirm'
    error.value = ''
    result.value = null
    wechatReady.value = false
  }
})

function startCooldown(sec: number) {
  if (timer) clearInterval(timer)
  cooldown.value = sec
  timer = setInterval(() => {
    cooldown.value -= 1
    if (cooldown.value <= 0) {
      cooldown.value = 0
      if (timer) clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function handleSubmit() {
  if (submitting.value || !canShare.value) return
  submitting.value = true
  error.value = ''
  try {
    const res = await apiCreateTuyileShare(plan.value.messages)
    result.value = res
    phase.value = 'done'
    if (inWechat.value) {
      const firstAssistant = plan.value.messages.find(m => m.role === 'assistant')
      wechatReady.value = await setTuyileWechatCard({
        title: '图一乐给我批了一卦，快来围观',
        desc: firstAssistant ? markdownSummary(firstAssistant.content, 50).replace(/</g, '＜').replace(/>/g, '＞') : '大罗天仙在线算命，图个乐呵～',
        link: res.url,
        imgUrl: shareIconUrl(),
      })
      showGuide.value = true
    }
  } catch (err: unknown) {
    if (err instanceof TuyileHttpError) {
      error.value = err.message
      if (err.status === 429 && err.retryAfter) startCooldown(err.retryAfter)
    } else {
      error.value = '网络开小差了，稍后再试～'
    }
  } finally {
    submitting.value = false
  }
}

const { siteUrl } = useSiteUrl()
function shareIconUrl(): string {
  return `${siteUrl}/android-chrome-512x512.png`
}

function selectUrl() {
  urlInput.value?.select()
}

async function handleCopy() {
  const url = result.value?.url
  if (!url) return
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(url)
    } else {
      // 兜底：非安全上下文 / 微信旧内核
      const el = urlInput.value
      if (!el) throw new Error('no input')
      el.select()
      if (!document.execCommand('copy')) throw new Error('copy failed')
    }
    showSuccess('链接已复制，快去分享吧～')
  } catch {
    selectUrl()
    showError('复制失败，请长按链接手动复制')
  }
}

function handleClose() {
  if (submitting.value) return
  emit('update:open', false)
}

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>
