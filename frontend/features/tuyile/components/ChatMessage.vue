<template>
  <div class="flex gap-2" :class="role === 'user' ? 'flex-row-reverse' : ''">
    <span v-if="role === 'assistant'" class="tuyile-avatar">乐</span>
    <div
      class="max-w-[85%] rounded-2xl px-3.5 py-2.5 text-sm leading-relaxed"
      :class="role === 'user'
        ? 'whitespace-pre-wrap break-words rounded-tr-sm bg-rose-600 text-white'
        : 'tuyile-md break-words rounded-tl-sm bg-gray-50 text-gray-800 dark:bg-gray-900/50 dark:text-gray-100'"
    >
      <!-- user：纯文本插值，Vue 自动转义 -->
      <template v-if="role === 'user'">{{ content }}</template>
      <!-- assistant：markdown-it html:false，原始 HTML 被转义，可安全 v-html -->
      <div v-else v-html="html" />
      <span v-if="streaming" class="tuyile-caret" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChatRole } from '../types'
import { renderMarkdown } from '../utils'

const props = defineProps<{ role: ChatRole, content: string, streaming?: boolean }>()

const html = computed(() => (props.role === 'assistant' ? renderMarkdown(props.content) : ''))
</script>
