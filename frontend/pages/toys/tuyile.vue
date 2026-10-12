<template>
  <div class="tuyile-root min-h-screen bg-stone-50 text-gray-800 transition-colors duration-300 dark:bg-gray-900 dark:text-gray-200">
    <!-- 返回玩具区 -->
    <NuxtLink
      to="/toys"
      class="fixed left-4 top-4 z-50 flex h-9 items-center gap-1.5 rounded-full bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
    >
      <span aria-hidden="true">←</span>
      <span>玩具</span>
    </NuxtLink>

    <button
      v-if="view !== 'home'"
      type="button"
      class="fixed left-24 top-4 z-50 flex h-9 items-center rounded-full bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
      @click="go('home')"
    >
      图一乐
    </button>

    <!-- 主题切换 -->
    <button
      type="button"
      class="fixed right-4 top-4 z-50 flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
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

    <Transition name="tuyile-fade" mode="out-in">
      <TuyileHome v-if="view === 'home'" key="home" @open="go('almanac')" />
      <AlmanacView
        v-else-if="view === 'almanac'"
        key="almanac"
        :date="date"
        :today="today"
        @update:date="setDate"
        @chat="go('chat')"
      />
      <div v-else key="chat">
        <ClientOnly>
          <FortuneChat @back="go('almanac')" />
        </ClientOnly>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { AlmanacView, FortuneChat, TuyileHome, beijingToday, isValidDateStr } from '~/features/tuyile'

definePageMeta({ layout: false })

type View = 'home' | 'almanac' | 'chat'

const route = useRoute()
const router = useRouter()
const { isDark, toggleTheme } = useTheme()

const today = beijingToday()

const view = computed<View>(() => {
  const v = route.query.view
  return v === 'almanac' || v === 'chat' ? v : 'home'
})

const date = computed(() => {
  const d = route.query.date
  return typeof d === 'string' && isValidDateStr(d) && d >= '1900-01-01' && d <= '2100-12-31' ? d : today
})

function go(v: View) {
  const query: Record<string, string> = {}
  if (v !== 'home') query.view = v
  if (v === 'almanac' && date.value !== today) query.date = date.value
  router.push({ query })
  if (import.meta.client) window.scrollTo({ top: 0 })
}

function setDate(d: string) {
  const query: Record<string, string> = { view: 'almanac' }
  if (d !== today) query.date = d
  router.replace({ query })
}

const { pageUrl } = useSiteUrl()
useSeoMeta({
  title: '图一乐 · 老黄历与 AI 运势',
  description: '图一乐：每日老黄历宜忌、冲煞、时辰吉凶，再找 AI 唠唠你的流年运势。仅供娱乐。',
  ogTitle: '图一乐 · 老黄历与 AI 运势',
  ogDescription: '每日宜忌 · 十二时辰吉凶 · AI 运势聊天，仅供娱乐。',
  ogUrl: pageUrl('/toys/tuyile'),
})
</script>

<style>
.tuyile-root .tuyile-serif {
  font-family: 'Noto Serif SC', 'Source Han Serif SC', 'Songti SC', 'STSong', 'SimSun', serif;
}
.tuyile-root .tuyile-paper {
  background:
    radial-gradient(circle at 12% 18%, rgba(225, 29, 72, 0.06), transparent 40%),
    radial-gradient(circle at 88% 82%, rgba(245, 158, 11, 0.08), transparent 45%),
    #fffaf3;
}
.dark .tuyile-root .tuyile-paper {
  background:
    radial-gradient(circle at 12% 18%, rgba(225, 29, 72, 0.12), transparent 40%),
    radial-gradient(circle at 88% 82%, rgba(245, 158, 11, 0.08), transparent 45%),
    #1f2937;
}
.tuyile-fade-enter-active,
.tuyile-fade-leave-active {
  transition: opacity 0.2s ease;
}
.tuyile-fade-enter-from,
.tuyile-fade-leave-to {
  opacity: 0;
}
.tuyile-root .tuyile-avatar {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: 9999px;
  background: linear-gradient(135deg, #e11d48, #f59e0b);
  color: #fff;
  font-size: 0.8rem;
  font-weight: 700;
}
.tuyile-root .tuyile-caret {
  display: inline-block;
  width: 0.45em;
  height: 1em;
  margin-left: 2px;
  vertical-align: -0.15em;
  background: currentColor;
  opacity: 0.6;
  animation: tuyile-blink 1s steps(1) infinite;
}
.tuyile-root .tuyile-dots i {
  display: inline-block;
  width: 6px;
  height: 6px;
  margin: 0 2px;
  border-radius: 9999px;
  background: currentColor;
  animation: tuyile-bounce 1.2s infinite ease-in-out;
}
.tuyile-root .tuyile-dots i:nth-child(2) { animation-delay: 0.15s; }
.tuyile-root .tuyile-dots i:nth-child(3) { animation-delay: 0.3s; }
@keyframes tuyile-blink {
  50% { opacity: 0; }
}
@keyframes tuyile-bounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}
.tuyile-root .tuyile-md > div > :first-child { margin-top: 0; }
.tuyile-root .tuyile-md > div > :last-child { margin-bottom: 0; }
.tuyile-root .tuyile-md p { margin: 0.4em 0; }
.tuyile-root .tuyile-md ul,
.tuyile-root .tuyile-md ol { margin: 0.4em 0; padding-left: 1.3em; }
.tuyile-root .tuyile-md ul { list-style: disc; }
.tuyile-root .tuyile-md ol { list-style: decimal; }
.tuyile-root .tuyile-md strong { font-weight: 700; }
.tuyile-root .tuyile-md h1,
.tuyile-root .tuyile-md h2,
.tuyile-root .tuyile-md h3,
.tuyile-root .tuyile-md h4 { margin: 0.6em 0 0.3em; font-weight: 700; }
.tuyile-root .tuyile-md a { color: #e11d48; text-decoration: underline; }
.tuyile-root .tuyile-md code { padding: 0 0.3em; border-radius: 4px; background: rgba(0, 0, 0, 0.06); }
.tuyile-root .tuyile-md blockquote { margin: 0.4em 0; padding-left: 0.8em; border-left: 3px solid #fda4af; opacity: 0.85; }
.tuyile-root .tuyile-md hr { margin: 0.6em 0; border-color: rgba(0, 0, 0, 0.1); }
</style>
