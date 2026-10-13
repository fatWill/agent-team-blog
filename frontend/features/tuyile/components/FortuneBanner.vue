<template>
  <button
    type="button"
    class="tuyile-banner group relative block w-full overflow-hidden rounded-2xl px-4 py-3.5 text-left text-white shadow-md shadow-rose-500/20 transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg hover:shadow-rose-500/30 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-300 active:translate-y-0 sm:px-5 sm:py-4"
    aria-label="聊一下你的运势，进入图一乐运势聊天"
    @click="emit('click')"
  >
    <span class="tuyile-banner-shine" aria-hidden="true" />
    <span class="relative flex items-center gap-3">
      <span class="tuyile-banner-orb flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-white/20 text-2xl backdrop-blur-sm sm:h-12 sm:w-12">
        🔮
      </span>
      <span class="min-w-0 flex-1">
        <span class="tuyile-serif block text-base font-bold leading-tight tracking-wide sm:text-lg">
          聊一下你的运势
        </span>
        <span class="mt-0.5 block truncate text-xs text-white/85 sm:text-sm">
          {{ slogan }}
        </span>
      </span>
      <span class="flex shrink-0 items-center gap-1 rounded-full bg-white/95 px-3 py-1.5 text-xs font-bold text-rose-600 shadow-sm transition-transform duration-200 group-hover:translate-x-0.5 sm:text-sm">
        开聊<span aria-hidden="true">→</span>
      </span>
    </span>
  </button>
</template>

<script setup lang="ts">
const emit = defineEmits<{ (e: 'click'): void }>()

const SLOGANS = [
  '大罗天仙坐镇 · 搬家跳槽表白取名，啥都能算',
  '本仙在线摆摊，问啥答啥，图个乐呵',
  '想知道今天几点放屁最吉利？来问就对了',
]
// SSR 与客户端保持一致：按日期取一条，避免 hydration 不一致
const slogan = SLOGANS[new Date().getDate() % SLOGANS.length]
</script>

<style scoped>
.tuyile-banner {
  background: linear-gradient(120deg, #e11d48 0%, #f43f5e 35%, #f97316 75%, #f59e0b 100%);
  background-size: 200% 100%;
  animation: tuyile-banner-flow 8s ease-in-out infinite alternate;
}
.tuyile-banner-shine {
  position: absolute;
  inset: 0;
  background: linear-gradient(105deg, transparent 40%, rgba(255, 255, 255, 0.28) 50%, transparent 60%);
  transform: translateX(-100%);
  animation: tuyile-banner-sweep 4.5s ease-in-out 1s infinite;
}
.tuyile-banner-orb {
  animation: tuyile-banner-float 3s ease-in-out infinite;
}
@keyframes tuyile-banner-flow {
  from { background-position: 0% 50%; }
  to { background-position: 100% 50%; }
}
@keyframes tuyile-banner-sweep {
  0% { transform: translateX(-100%); }
  35%, 100% { transform: translateX(100%); }
}
@keyframes tuyile-banner-float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}
@media (prefers-reduced-motion: reduce) {
  .tuyile-banner,
  .tuyile-banner-shine,
  .tuyile-banner-orb {
    animation: none;
  }
}
</style>
