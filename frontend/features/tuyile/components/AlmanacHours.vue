<template>
  <section class="rounded-xl border border-rose-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-5">
    <h3 class="tuyile-serif mb-3 flex items-center gap-2 text-base font-bold text-gray-800 dark:text-gray-100">
      <span class="h-4 w-1 rounded-full bg-rose-500" />十二时辰吉凶
    </h3>
    <div class="grid grid-cols-3 gap-2 sm:grid-cols-4 lg:grid-cols-6">
      <button
        v-for="(h, i) in hours"
        :key="h.zhi"
        type="button"
        class="relative rounded-lg border p-2 text-center transition-all duration-150 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-400"
        :class="[
          i === activeIndex ? 'ring-2 ring-rose-400' : '',
          h.luck === '吉'
            ? 'border-emerald-100 bg-emerald-50/60 hover:bg-emerald-50 dark:border-emerald-500/20 dark:bg-emerald-500/10'
            : 'border-gray-100 bg-gray-50 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900/40',
        ]"
        :aria-expanded="i === activeIndex"
        @click="activeIndex = activeIndex === i ? -1 : i"
      >
        <div class="tuyile-serif text-lg font-bold text-gray-800 dark:text-gray-100">{{ h.zhi }}时</div>
        <div class="text-[11px] text-gray-400 dark:text-gray-500">{{ h.range }}</div>
        <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ h.ganZhi }} · {{ h.tianShen }}</div>
        <span
          class="absolute right-1 top-1 rounded px-1 text-[10px] font-bold"
          :class="h.luck === '吉' ? 'bg-emerald-500 text-white' : 'bg-gray-400 text-white dark:bg-gray-600'"
        >
          {{ h.luck }}
        </span>
        <span v-if="i === currentIndex" class="absolute left-1 top-1 rounded bg-rose-500 px-1 text-[10px] text-white">此刻</span>
      </button>
    </div>

    <Transition name="tuyile-fade" mode="out-in">
      <div
        v-if="active"
        :key="active.zhi"
        class="mt-3 space-y-2 rounded-lg bg-rose-50/50 p-3 text-sm leading-relaxed text-gray-600 dark:bg-gray-900/40 dark:text-gray-300"
      >
        <div class="font-semibold text-gray-800 dark:text-gray-100">
          {{ active.zhi }}时（{{ active.range }}）· {{ active.ganZhi }} · {{ active.tianShen }} · {{ active.luck }}
        </div>
        <div>冲{{ active.chong }} 煞{{ active.sha }}</div>
        <div><span class="font-semibold text-emerald-600 dark:text-emerald-400">宜：</span>{{ join(active.yi) }}</div>
        <div><span class="font-semibold text-rose-600 dark:text-rose-400">忌：</span>{{ join(active.ji) }}</div>
        <div><span class="text-gray-500">吉神：</span>{{ join(active.jiShen) }}</div>
        <div><span class="text-gray-500">凶煞：</span>{{ join(active.xiongSha) }}</div>
      </div>
    </Transition>
  </section>
</template>

<script setup lang="ts">
import type { AlmanacHour } from '../types'

const props = defineProps<{ hours: AlmanacHour[], currentIndex: number }>()

const activeIndex = ref(props.currentIndex)
watch(() => props.hours, () => {
  activeIndex.value = props.currentIndex
})

const active = computed(() => props.hours[activeIndex.value])

function join(list: string[]): string {
  return list.length ? list.join(' ') : '—'
}
</script>
