<template>
  <header class="sticky top-0 z-50 border-b border-gray-100 bg-white/80 backdrop-blur-lg">
    <div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
      <!-- Logo -->
      <NuxtLink to="/" class="flex items-center gap-2.5 group">
        <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-primary-500 text-white font-bold text-lg transition-transform duration-200 group-hover:scale-105">
          A
        </div>
        <span class="text-lg font-bold text-gray-900 hidden sm:block">fatwill</span>
      </NuxtLink>

      <!-- 桌面端导航 -->
      <nav class="hidden md:flex items-center gap-1">
        <NuxtLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="rounded-lg px-4 py-2 text-sm font-medium text-gray-600 transition-all duration-200 hover:bg-gray-100 hover:text-gray-900"
          active-class="!bg-primary-50 !text-primary-600"
        >
          {{ item.label }}
        </NuxtLink>
      </nav>

      <!-- 右侧操作区 -->
      <div class="flex items-center gap-3">
        <!-- 搜索按钮 -->
        <button
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-all duration-200 hover:bg-gray-100 hover:text-gray-700"
          aria-label="搜索"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>

        <!-- 移动端菜单按钮 -->
        <button
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-all duration-200 hover:bg-gray-100 hover:text-gray-700 md:hidden"
          aria-label="菜单"
          @click="isMobileMenuOpen = !isMobileMenuOpen"
        >
          <svg v-if="!isMobileMenuOpen" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 移动端导航菜单 -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="isMobileMenuOpen" class="border-t border-gray-100 bg-white px-4 py-3 md:hidden">
        <nav class="flex flex-col gap-1">
          <NuxtLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="rounded-lg px-4 py-2.5 text-sm font-medium text-gray-600 transition-all duration-200 hover:bg-gray-100 hover:text-gray-900"
            active-class="!bg-primary-50 !text-primary-600"
            @click="isMobileMenuOpen = false"
          >
            {{ item.label }}
          </NuxtLink>
        </nav>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import type { NavItem } from '~/types'

const isMobileMenuOpen = ref(false)

const navItems: NavItem[] = [
  { label: '首页', to: '/' },
  { label: '文章', to: '/articles' },
  { label: '分类', to: '/categories' },
  { label: '关于', to: '/about' },
]
</script>
