<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="cat in categories"
      :key="cat.id"
      class="rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200"
      :class="
        activeSlug === cat.slug
          ? 'bg-primary-500 text-white shadow-sm shadow-primary-200'
          : 'bg-white text-gray-600 border border-gray-200 hover:border-primary-200 hover:text-primary-600 hover:bg-primary-50'
      "
      @click="$emit('select', cat.slug)"
    >
      {{ cat.name }}
      <span
        v-if="cat.articleCount"
        class="ml-1 text-xs"
        :class="activeSlug === cat.slug ? 'text-primary-200' : 'text-gray-400'"
      >
        {{ cat.articleCount }}
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
import type { Category } from '~/types'

defineProps<{
  categories: Category[]
  activeSlug: string
}>()

defineEmits<{
  select: [slug: string]
}>()
</script>
