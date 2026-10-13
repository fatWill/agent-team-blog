<template>
  <TuyileShell>
    <template #nav>
      <button
        v-if="view !== 'home'"
        type="button"
        class="fixed left-24 top-4 z-50 flex h-9 items-center rounded-full bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
        @click="go('home')"
      >
        图一乐
      </button>
    </template>

    <Transition name="tuyile-fade" mode="out-in">
      <TuyileHome v-if="view === 'home'" key="home" @open="go('almanac')" @chat="go('chat')" />
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
  </TuyileShell>
</template>

<script setup lang="ts">
import {
  AlmanacView,
  FortuneChat,
  TuyileHome,
  TuyileShell,
  beijingToday,
  isValidDateStr,
  setTuyileWechatCard,
} from '~/features/tuyile'

definePageMeta({ layout: false })

type View = 'home' | 'almanac' | 'chat'

const route = useRoute()
const router = useRouter()

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

const { pageUrl, siteUrl } = useSiteUrl()
useSeoMeta({
  title: '图一乐 · 老黄历与 AI 运势',
  description: '图一乐：每日老黄历宜忌、冲煞、时辰吉凶，再找大罗天仙唠唠搬家跳槽取名整活，啥都能算。仅供娱乐。',
  ogTitle: '图一乐 · 老黄历与 AI 运势',
  ogDescription: '每日宜忌 · 十二时辰吉凶 · AI 算命脱口秀，仅供娱乐。',
  ogUrl: pageUrl('/toys/tuyile'),
})

// 微信内：默认卡片指向图一乐首页；query 变化（单页内切换视图）后在新 URL 上重新签名
function applyDefaultWechatCard() {
  setTuyileWechatCard({
    title: '图一乐 · 老黄历与 AI 运势',
    desc: '大罗天仙在线摆摊，搬家跳槽取名整活，啥都能算，图个乐呵～',
    link: pageUrl('/toys/tuyile'),
    imgUrl: `${siteUrl}/android-chrome-512x512.png`,
  })
}
onMounted(applyDefaultWechatCard)
watch(() => route.fullPath, () => {
  nextTick(applyDefaultWechatCard)
})
</script>
