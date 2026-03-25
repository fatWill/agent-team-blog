<template>
  <aside class="space-y-6">
    <!-- 博主信息卡片 -->
    <div class="card p-6 text-center">
      <img
        src="https://api.dicebear.com/9.x/adventurer/svg?seed=Felix"
        alt="博主头像"
        class="mx-auto mb-4 h-20 w-20 rounded-full bg-primary-50 p-1"
      />
      <h3 class="mb-1 text-base font-semibold text-gray-900">Agent Team</h3>
      <p class="mb-4 text-sm text-gray-500">全栈开发团队，热爱技术与开源</p>
      <div class="flex justify-center gap-6 text-center">
        <div>
          <div class="text-lg font-bold text-gray-900">{{ stats.articles }}</div>
          <div class="text-xs text-gray-400">文章</div>
        </div>
        <div>
          <div class="text-lg font-bold text-gray-900">{{ stats.categories }}</div>
          <div class="text-xs text-gray-400">分类</div>
        </div>
        <div>
          <div class="text-lg font-bold text-gray-900">{{ stats.tags }}</div>
          <div class="text-xs text-gray-400">标签</div>
        </div>
      </div>
    </div>

    <!-- 热门标签 -->
    <div class="card p-6">
      <h3 class="mb-4 text-sm font-semibold text-gray-900">热门标签</h3>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="tag in tags"
          :key="tag.id"
          class="cursor-pointer rounded-md bg-gray-50 px-3 py-1.5 text-xs font-medium text-gray-600 transition-all duration-200 hover:bg-primary-50 hover:text-primary-600"
        >
          # {{ tag.name }}
        </span>
      </div>
    </div>

    <!-- 推荐阅读 -->
    <div class="card p-6">
      <h3 class="mb-4 text-sm font-semibold text-gray-900">推荐阅读</h3>
      <ul class="space-y-3">
        <li v-for="(article, index) in recommendedArticles" :key="article.id">
          <NuxtLink
            :to="`/articles/${article.id}`"
            class="group flex items-start gap-3"
          >
            <span class="flex h-6 w-6 shrink-0 items-center justify-center rounded bg-primary-50 text-xs font-bold text-primary-500">
              {{ index + 1 }}
            </span>
            <span class="text-sm text-gray-600 leading-relaxed transition-colors duration-200 group-hover:text-primary-600 line-clamp-2">
              {{ article.title }}
            </span>
          </NuxtLink>
        </li>
      </ul>
    </div>
  </aside>
</template>

<script setup lang="ts">
import type { Article, Tag } from '~/types'

defineProps<{
  tags: Tag[]
  recommendedArticles: Article[]
  stats: {
    articles: number
    categories: number
    tags: number
  }
}>()
</script>
