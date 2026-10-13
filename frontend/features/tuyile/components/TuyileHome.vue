<template>
  <div class="mx-auto max-w-3xl px-4 pb-16 pt-20">
    <div class="mb-10 text-center">
      <div class="mb-4 inline-flex items-center gap-2 rounded-full bg-rose-500/10 px-4 py-1.5 text-sm font-medium text-rose-600 dark:bg-rose-500/20 dark:text-rose-300">
        🎲 仅供娱乐 · 图个乐呵
      </div>
      <h1 class="tuyile-serif mb-3 text-4xl font-black tracking-widest text-gray-900 dark:text-gray-100 sm:text-5xl">
        图一乐
      </h1>
      <p class="text-sm leading-relaxed text-gray-500 dark:text-gray-400">
        翻翻老黄历，再找图一乐唠唠今年运势
      </p>
    </div>

    <FortuneBanner class="mb-5" @click="emit('chat')" />

    <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <button
        v-for="entry in entries"
        :key="entry.key"
        type="button"
        :aria-disabled="!entry.ready"
        class="group relative overflow-hidden rounded-xl border bg-white p-5 text-left transition-all duration-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-400 dark:bg-gray-800"
        :class="entry.ready
          ? 'border-rose-100 shadow-sm hover:-translate-y-0.5 hover:shadow-md dark:border-gray-700 dark:hover:border-rose-500/40'
          : 'cursor-not-allowed border-gray-100 opacity-70 dark:border-gray-700'"
        @click="handleClick(entry)"
      >
        <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl text-2xl" :class="entry.iconBg">
          {{ entry.icon }}
        </div>
        <h2 class="tuyile-serif mb-1 text-lg font-bold text-gray-900 dark:text-gray-100" :class="entry.ready ? 'group-hover:text-rose-600 dark:group-hover:text-rose-400' : ''">
          {{ entry.title }}
        </h2>
        <p class="text-sm leading-relaxed text-gray-500 dark:text-gray-400">{{ entry.desc }}</p>
        <span
          v-if="!entry.ready"
          class="absolute right-3 top-3 rounded-full bg-gray-100 px-2 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-gray-700 dark:text-gray-400"
        >
          即将上线
        </span>
        <span
          v-else
          class="absolute right-3 top-3 rounded-full bg-rose-500/90 px-2 py-0.5 text-[10px] font-bold text-white"
        >
          NEW
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { showInfo } from '~/utils/ui'
import FortuneBanner from './FortuneBanner.vue'

interface Entry {
  key: 'horoscope' | 'bazi' | 'almanac'
  title: string
  desc: string
  icon: string
  iconBg: string
  ready: boolean
}

const emit = defineEmits<{
  (e: 'open', key: 'almanac'): void
  (e: 'chat'): void
}>()

const entries: Entry[] = [
  { key: 'horoscope', title: '星座运势', desc: '十二星座今日运势', icon: '✨', iconBg: 'bg-indigo-50 dark:bg-indigo-500/10', ready: false },
  { key: 'bazi', title: '生辰八字', desc: '排盘看五行喜忌', icon: '☯️', iconBg: 'bg-amber-50 dark:bg-amber-500/10', ready: false },
  { key: 'almanac', title: '老黄历', desc: '每日宜忌 · 时辰吉凶 · 聊聊运势', icon: '📜', iconBg: 'bg-rose-50 dark:bg-rose-500/10', ready: true },
]

function handleClick(entry: Entry) {
  if (entry.key === 'almanac') {
    emit('open', 'almanac')
    return
  }
  showInfo(`「${entry.title}」即将上线，先去翻翻老黄历吧～`)
}
</script>
