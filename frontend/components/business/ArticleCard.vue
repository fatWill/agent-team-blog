<template>
  <NuxtLink :to="`/articles/${article.id}`" class="card-hover group block">
    <!-- 封面图 -->
    <div class="relative aspect-[16/9] overflow-hidden rounded-t-xl bg-gray-100">
      <img
        :src="article.coverImage"
        :alt="article.title"
        class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        loading="lazy"
      />
      <!-- 分类标签 -->
      <span class="absolute left-3 top-3 rounded-md bg-white/90 px-2.5 py-1 text-xs font-medium text-gray-700 backdrop-blur-sm">
        {{ article.category.name }}
      </span>
    </div>

    <!-- 内容区 -->
    <div class="p-5">
      <!-- 标签 -->
      <div class="mb-3 flex flex-wrap gap-1.5">
        <span
          v-for="tag in article.tags.slice(0, 3)"
          :key="tag.id"
          class="rounded-md bg-primary-50 px-2 py-0.5 text-xs font-medium text-primary-600"
        >
          {{ tag.name }}
        </span>
      </div>

      <!-- 标题 -->
      <h3 class="mb-2 line-clamp-2 text-lg font-semibold text-gray-900 transition-colors duration-200 group-hover:text-primary-600">
        {{ article.title }}
      </h3>

      <!-- 摘要 -->
      <p class="mb-4 line-clamp-2 text-sm leading-relaxed text-gray-500">
        {{ article.summary }}
      </p>

      <!-- 底部信息 -->
      <div class="flex items-center justify-between border-t border-gray-50 pt-4">
        <div class="flex items-center gap-2">
          <img
            :src="article.author.avatar"
            :alt="article.author.name"
            class="h-6 w-6 rounded-full bg-gray-100"
          />
          <span class="text-xs font-medium text-gray-600">{{ article.author.name }}</span>
        </div>
        <div class="flex items-center gap-3 text-xs text-gray-400">
          <span>{{ article.createdAt }}</span>
          <span>·</span>
          <span>{{ article.readingTime }} 分钟</span>
        </div>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import type { Article } from '~/types'

defineProps<{
  article: Article
}>()
</script>
