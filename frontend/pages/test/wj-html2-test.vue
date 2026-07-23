<template>
  <div class="min-h-screen bg-white p-6 font-sans text-gray-800">
    <!-- 顶部状态栏 -->
    <div class="mb-4 flex items-center justify-between">
      <div class="flex items-center gap-4 text-sm">
        <span class="rounded bg-blue-100 px-2 py-1 text-blue-700">
          当前步骤：{{ stepLabel }}
        </span>
        <span
          :class="isInIframe ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'"
          class="rounded px-2 py-1"
        >
          {{ isInIframe ? '✅ iframe 环境' : '⚠️ 非 iframe 环境' }}
        </span>
      </div>
      <button
        class="rounded bg-gray-200 px-3 py-1 text-sm text-gray-600 transition-colors hover:bg-gray-300"
        @click="handleReset"
      >
        重置
      </button>
    </div>

    <!-- 独立测试模式警告 -->
    <div
      v-if="!isInIframe"
      class="mb-4 rounded border border-orange-300 bg-orange-50 p-3 text-sm text-orange-700"
    >
      ⚠️ 当前非 iframe 环境，postMessage 无接收方。所有消息将发送到自身 window.parent（即自身）。
    </div>

    <!-- 上一次发送的 payload -->
    <div v-if="lastPayload" class="mb-4 rounded border border-gray-200 bg-gray-50 p-3">
      <div class="mb-1 text-xs text-gray-500">上次发送的 payload：</div>
      <pre class="text-sm text-gray-800">{{ JSON.stringify(lastPayload, null, 2) }}</pre>
    </div>

    <!-- 主交互区域 -->
    <div class="mb-6 flex min-h-[200px] items-center justify-center">
      <!-- 第一步：选肉类 -->
      <div v-if="currentStep === 1" class="text-center">
        <p class="mb-4 text-lg font-medium text-gray-600">请选择一种肉类</p>
        <div class="flex gap-6">
          <button
            class="rounded-lg bg-red-600 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-red-700 hover:shadow-md active:scale-95"
            @click="handleSelect('肉类', '牛肉')"
          >
            🥩 牛肉
          </button>
          <button
            class="rounded-lg bg-rose-500 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-rose-600 hover:shadow-md active:scale-95"
            @click="handleSelect('肉类', '猪肉')"
          >
            🐖 猪肉
          </button>
        </div>
      </div>

      <!-- 第二步：选蔬菜 -->
      <div v-else-if="currentStep === 2" class="text-center">
        <p class="mb-4 text-lg font-medium text-gray-600">请选择一种蔬菜</p>
        <div class="flex gap-6">
          <button
            class="rounded-lg bg-orange-500 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-orange-600 hover:shadow-md active:scale-95"
            @click="handleSelect('蔬菜', '胡萝卜')"
          >
            🥕 胡萝卜
          </button>
          <button
            class="rounded-lg bg-green-600 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-green-700 hover:shadow-md active:scale-95"
            @click="handleSelect('蔬菜', '白菜')"
          >
            🥬 白菜
          </button>
        </div>
      </div>

      <!-- 第三步：选水果 -->
      <div v-else-if="currentStep === 3" class="text-center">
        <p class="mb-4 text-lg font-medium text-gray-600">请选择一种水果</p>
        <div class="flex gap-6">
          <button
            class="rounded-lg bg-red-500 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-red-600 hover:shadow-md active:scale-95"
            @click="handleSelect('水果', '苹果')"
          >
            🍎 苹果
          </button>
          <button
            class="rounded-lg bg-yellow-400 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-yellow-500 hover:shadow-md active:scale-95"
            @click="handleSelect('水果', '香蕉')"
          >
            🍌 香蕉
          </button>
        </div>
      </div>

      <!-- 第四步：选水果（第二轮） -->
      <div v-else-if="currentStep === 4" class="text-center">
        <p class="mb-4 text-lg font-medium text-gray-600">请选择一种水果（第二轮）</p>
        <div class="flex gap-6">
          <button
            class="rounded-lg bg-lime-500 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-lime-600 hover:shadow-md active:scale-95"
            @click="handleSelect('水果', '梨子')"
          >
            🍐 梨子
          </button>
          <button
            class="rounded-lg bg-amber-500 px-10 py-4 text-xl font-bold text-white shadow transition-all hover:bg-amber-600 hover:shadow-md active:scale-95"
            @click="handleSelect('水果', '柚子')"
          >
            🍊 柚子
          </button>
        </div>
      </div>

      <!-- 完成状态 -->
      <div v-else class="text-center">
        <div class="mb-2 text-4xl">✅</div>
        <p class="text-xl font-bold text-green-600">作答已完成</p>
        <p class="mt-2 text-sm text-gray-500">所有变量已上报，结束关键词已发送</p>
      </div>
    </div>

    <!-- 日志面板 -->
    <div class="rounded border border-gray-200 bg-gray-50">
      <div class="flex items-center justify-between border-b border-gray-200 px-4 py-2">
        <span class="text-sm font-medium text-gray-600">📋 操作日志</span>
        <span class="text-xs text-gray-400">共 {{ logs.length }} 条</span>
      </div>
      <div class="max-h-[300px] overflow-y-auto p-3 font-mono text-xs">
        <div v-if="logs.length === 0" class="py-4 text-center text-gray-400">暂无日志</div>
        <div
          v-for="(log, index) in logs"
          :key="index"
          class="mb-1 rounded px-2 py-1"
          :class="logColorClass(log.type)"
        >
          <span class="text-gray-400">{{ log.time }}</span>
          <span class="ml-2 font-medium">{{ log.label }}</span>
          <span class="ml-2">{{ log.message }}</span>
          <pre v-if="log.payload" class="ml-4 mt-1 whitespace-pre-wrap text-gray-600">{{ log.payload }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

useHead({
  title: 'wj-html2-test | iframe 通信测试（4步版）',
})

interface LogEntry {
  time: string
  type: 'info' | 'send' | 'warn' | 'error'
  label: string
  message: string
  payload?: string
}

type CategoryKey = '肉类' | '蔬菜' | '水果'

const currentStep = ref<1 | 2 | 3 | 4 | 5>(1)
const logs = ref<LogEntry[]>([])
const lastPayload = ref<Record<string, unknown> | null>(null)
const isInIframe = ref(false)
const stepStartAt = ref<number>(0)

const stepLabel = computed(() => {
  if (currentStep.value === 1) return '第1步 - 选择肉类'
  if (currentStep.value === 2) return '第2步 - 选择蔬菜'
  if (currentStep.value === 3) return '第3步 - 选择水果'
  if (currentStep.value === 4) return '第4步 - 选择水果（第二轮）'
  return '已完成'
})

function getTimestamp(): string {
  const now = new Date()
  return `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}.${String(now.getMilliseconds()).padStart(3, '0')}`
}

function addLog(type: LogEntry['type'], label: string, message: string, payload?: unknown) {
  logs.value.push({
    time: getTimestamp(),
    type,
    label,
    message,
    payload: payload ? JSON.stringify(payload, null, 2) : undefined,
  })
}

function logColorClass(type: LogEntry['type']): string {
  switch (type) {
    case 'info': return 'bg-blue-50 text-blue-700'
    case 'send': return 'bg-green-50 text-green-700'
    case 'warn': return 'bg-orange-50 text-orange-700'
    case 'error': return 'bg-red-50 text-red-700'
    default: return ''
  }
}

function generateRandomToken(length = 10): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

function sendPostMessage(payload: Record<string, unknown>) {
  console.log('[wj-html2-test] 发送 postMessage，payload:', JSON.stringify(payload), 'targetOrigin: *', '时间:', getTimestamp())
  window.parent.postMessage(payload, '*')
  lastPayload.value = payload
  addLog('send', '[postMessage]', '已发送', payload)
}

function handleSelect(key: CategoryKey, value: string) {
  const now = Date.now()
  const difftime = now - stepStartAt.value
  console.log(`[wj-html2-test] 用户选择 ${key} = ${value}（耗时 ${difftime}ms）`)
  addLog('info', '[用户操作]', `选择了 ${key}: ${value}（耗时 ${difftime}ms）`)

  // 上报当前步骤的变量（含 difftime）
  const payload: Record<string, unknown> = { [key]: value, difftime }
  console.log(`[wj-html2-test] 准备发送第${currentStep.value}步 payload:`, JSON.stringify(payload))
  sendPostMessage(payload)

  if (currentStep.value < 4) {
    // 非最后一步，递进
    currentStep.value = (currentStep.value + 1) as 1 | 2 | 3 | 4 | 5
    stepStartAt.value = Date.now()
    const nextLabels: Record<number, string> = {
      2: '进入第二步 - 选择蔬菜',
      3: '进入第三步 - 选择水果',
      4: '进入第四步 - 选择水果（第二轮）',
    }
    const msg = nextLabels[currentStep.value] || ''
    addLog('info', '[流程]', msg)
    console.log(`[wj-html2-test] ${msg}`)
  } else {
    // 第 4 步：发送结束关键词
    const token = generateRandomToken()
    console.log('[wj-html2-test] 生成结束关键词随机字符串:', token)
    addLog('info', '[生成]', `结束关键词携带随机字符串 = ${token}`)

    const endPayload = { '__WJ_IFRAME_QUESTION_END__': token }
    console.log('[wj-html2-test] 准备发送结束 payload:', JSON.stringify(endPayload))
    sendPostMessage(endPayload)

    currentStep.value = 5
    addLog('info', '[流程]', '作答完毕，已发送结束关键词')
    console.log('[wj-html2-test] 已发送结束关键词，本题作答完毕')
  }
}

function handleReset() {
  currentStep.value = 1
  logs.value = []
  lastPayload.value = null
  stepStartAt.value = Date.now()
  console.log('[wj-html2-test] 页面已重置，回到第一步')
  addLog('info', '[重置]', '页面已重置，等待用户选择第一步')
}

// 页面初始化
onMounted(() => {
  isInIframe.value = window.parent !== window
  stepStartAt.value = Date.now()

  if (!isInIframe.value) {
    console.warn('[wj-html2-test] 当前不是在 iframe 中运行，postMessage 将发送到自身 window')
    addLog('warn', '[环境检测]', '非 iframe 环境，postMessage 无接收方')
  } else {
    console.info('[wj-html2-test] 检测到 iframe 环境，postMessage 将发送到 parent')
    addLog('info', '[环境检测]', 'iframe 环境就绪')
  }

  console.log('[wj-html2-test] 页面初始化完成，等待用户选择第一步')
  addLog('info', '[初始化]', '页面加载完成，等待用户操作')
})
</script>
