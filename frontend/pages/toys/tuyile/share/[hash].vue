<template>
  <TuyileShell>
    <div class="mx-auto max-w-3xl px-4 pb-16 pt-16">
      <!-- 404 / 错误 -->
      <div v-if="notFound || loadError" class="mt-10 rounded-2xl bg-white p-10 text-center shadow-sm dark:bg-gray-800">
        <div class="mb-3 text-5xl" aria-hidden="true">{{ notFound ? '🫥' : '😵' }}</div>
        <h1 class="tuyile-serif mb-2 text-xl font-bold text-gray-900 dark:text-gray-100">
          {{ notFound ? '这一卦找不到啦' : '卦象一时没加载出来' }}
        </h1>
        <p class="mb-6 text-sm leading-relaxed text-gray-500 dark:text-gray-400">
          {{ notFound ? '分享链接不存在或已失效，可能是链接复制不完整。' : '网络开小差了，刷新试试～' }}
        </p>
        <div class="flex flex-col justify-center gap-2 sm:flex-row">
          <NuxtLink to="/toys/tuyile?view=chat" class="rounded-xl bg-rose-600 px-5 py-2.5 text-sm font-medium text-white transition-colors hover:bg-rose-700">
            🔮 我也去算一卦
          </NuxtLink>
          <button v-if="loadError" type="button" class="rounded-xl border border-gray-200 px-5 py-2.5 text-sm text-gray-600 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-700" @click="refresh()">
            重新加载
          </button>
        </div>
      </div>

      <template v-else-if="data">
        <header class="mb-4 text-center">
          <div class="mb-2 inline-flex items-center gap-1.5 rounded-full bg-rose-500/10 px-3 py-1 text-xs font-medium text-rose-600 dark:bg-rose-500/20 dark:text-rose-300">
            🔗 对话分享 · 只读
          </div>
          <h1 class="tuyile-serif text-2xl font-black tracking-wide text-gray-900 dark:text-gray-100">图一乐 · 一卦实录</h1>
          <p v-if="createdText" class="mt-1 text-xs text-gray-400 dark:text-gray-500">分享于 {{ createdText }}</p>
        </header>

        <article class="space-y-4 rounded-xl border border-gray-100 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-4" aria-label="分享的对话">
          <ChatMessage v-for="(m, i) in data.messages" :key="i" :role="m.role" :content="m.content" />
        </article>

        <p class="mt-4 text-center text-xs leading-relaxed text-gray-400 dark:text-gray-500">
          🎲 以上内容由 AI 生成，仅供娱乐，不构成任何医疗、投资、法律建议。
        </p>

        <div class="mt-6 text-center">
          <NuxtLink
            to="/toys/tuyile?view=chat"
            class="inline-flex items-center gap-2 rounded-2xl bg-gradient-to-r from-rose-600 to-amber-500 px-7 py-3.5 text-base font-bold text-white shadow-lg shadow-rose-500/25 transition-all duration-200 hover:scale-105 hover:shadow-xl active:scale-95"
          >
            🔮 我也去算一卦
          </NuxtLink>
          <p class="mt-2 text-xs text-gray-400 dark:text-gray-500">开启一段全新的对话</p>
        </div>
      </template>

      <div v-else class="space-y-4">
        <div class="mx-auto h-8 w-48 animate-pulse rounded-lg bg-white dark:bg-gray-800" />
        <div class="h-64 animate-pulse rounded-xl bg-white shadow-sm dark:bg-gray-800" />
      </div>
    </div>
  </TuyileShell>
</template>

<script setup lang="ts">
import { ChatMessage, TuyileShell, apiFetchTuyileShare, markdownSummary, setTuyileWechatCard } from '~/features/tuyile'

definePageMeta({ layout: false })

const route = useRoute()
const hash = computed(() => String(route.params.hash || ''))
const HASH_RE = /^[0-9A-Za-z]{12}$/

const { data, error, refresh } = await useAsyncData(
  () => `tuyile-share-${hash.value}`,
  () => {
    if (!HASH_RE.test(hash.value)) {
      throw createError({ statusCode: 404, statusMessage: '分享不存在或已失效' })
    }
    return apiFetchTuyileShare(hash.value)
  },
  { watch: [hash] },
)

function errStatus(e: unknown): number {
  if (typeof e === 'object' && e !== null && 'statusCode' in e) {
    const s = (e as { statusCode: unknown }).statusCode
    if (typeof s === 'number') return s
  }
  return 0
}

const notFound = computed(() => !!error.value && errStatus(error.value) === 404)
const loadError = computed(() => !!error.value && !notFound.value)

// SSR 下 404 返回真实 404 状态码，避免搜索引擎收录失效链接
if (import.meta.server && notFound.value) {
  setResponseStatus(useRequestEvent()!, 404)
}

const createdText = computed(() => {
  const s = data.value?.createdAt
  if (!s) return ''
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    timeZone: 'Asia/Shanghai',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d)
})

const firstAssistant = computed(() => data.value?.messages.find(m => m.role === 'assistant')?.content || '')
const firstUser = computed(() => data.value?.messages.find(m => m.role === 'user')?.content || '')

const { siteUrl, pageUrl } = useSiteUrl()
const pageLink = computed(() => pageUrl(`/toys/tuyile/share/${hash.value}`))
const shareImage = `${siteUrl}/android-chrome-512x512.png`

const seoTitle = computed(() => {
  if (notFound.value) return '分享不存在 · 图一乐'
  const q = firstUser.value ? metaSafe(markdownSummary(firstUser.value, 24)) : ''
  return q ? `「${q}」图一乐这样说 · 图一乐` : '图一乐 · 一卦实录'
})
const seoDesc = computed(() => (firstAssistant.value
  ? metaSafe(markdownSummary(firstAssistant.value, 100))
  : '大罗天仙在线算命，搬家跳槽取名整活，啥都能算，仅供娱乐。'))

/** meta 文本纵深防御：属性值虽已转义引号，仍把尖括号替换为全角，避免任何解析歧义 */
function metaSafe(s: string): string {
  return s.replace(/</g, '＜').replace(/>/g, '＞')
}

useSeoMeta({
  title: () => seoTitle.value,
  description: () => seoDesc.value,
  ogTitle: () => seoTitle.value,
  ogDescription: () => seoDesc.value,
  ogUrl: () => pageLink.value,
  ogImage: shareImage,
  ogType: 'article',
  twitterCard: 'summary',
  twitterTitle: () => seoTitle.value,
  twitterDescription: () => seoDesc.value,
  // 用户生成的快照内容，不希望被搜索引擎收录
  robots: 'noindex, follow',
})

// 微信内打开分享页：设置卡片，方便二次转发；SPA 内跳转到新 hash 时在新 URL 上重签
function applyWechatCard() {
  if (!data.value) return
  setTuyileWechatCard({
    title: seoTitle.value,
    desc: seoDesc.value,
    link: pageLink.value,
    imgUrl: shareImage,
  })
}
onMounted(applyWechatCard)
watch(() => data.value?.hash, () => nextTick(applyWechatCard))
</script>
