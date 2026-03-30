<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 transition-colors dark:bg-gray-900">
    <div class="w-full max-w-sm">
      <!-- Logo -->
      <div class="mb-8 text-center">
        <NuxtLink to="/" class="inline-block">
          <img src="/avatar.png" alt="fatwillzeng" class="mx-auto h-16 w-16 rounded-full ring-2 ring-gray-200 dark:ring-gray-700" />
        </NuxtLink>
        <h1 class="mt-4 text-xl font-bold text-gray-900 dark:text-gray-100">登录</h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">登录后可进入管理后台</p>
      </div>

      <!-- 登录表单 -->
      <form class="space-y-4" @submit.prevent="handleLogin">
        <div>
          <label for="username" class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            账号
          </label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            autocomplete="username"
            required
            class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
            placeholder="请输入账号"
          />
        </div>

        <div>
          <label for="password" class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            密码
          </label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            required
            class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
            placeholder="请输入密码"
          />
        </div>

        <!-- 错误提示 -->
        <p v-if="errorMsg" class="text-sm text-red-500">
          {{ errorMsg }}
        </p>

        <button
          type="submit"
          :disabled="submitting"
          class="w-full rounded-lg bg-primary-500 px-4 py-2.5 text-sm font-medium text-white transition-all duration-200 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-300 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ submitting ? '登录中...' : '登录' }}
        </button>
      </form>

      <!-- 返回首页 -->
      <div class="mt-6 text-center">
        <NuxtLink to="/" class="text-sm text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300">
          ← 返回首页
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { apiLogin } from '~/utils/api'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: '',
})
const submitting = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  errorMsg.value = ''
  submitting.value = true
  try {
    await apiLogin({
      username: form.username,
      password: form.password,
    })
    authStore.setLoggedIn(true)
    router.push('/admin')
  } catch (err: unknown) {
    const fetchErr = err as { data?: { statusMessage?: string } }
    errorMsg.value = fetchErr?.data?.statusMessage || '登录失败，请重试'
  } finally {
    submitting.value = false
  }
}

// 初始化主题
const { initTheme } = useTheme()
onMounted(() => {
  initTheme()
})

useHead({
  title: '登录 - fatwillzeng',
})
</script>
