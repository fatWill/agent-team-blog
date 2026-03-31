/**
 * 主题切换组合式函数
 * - 初始主题：服务端从 cookie 读取（SSR 零闪烁）
 * - 切换时：写 cookie + 调用 /api/theme 存入 Redis
 */
export function useTheme() {
  // 服务端/客户端均可读取 cookie
  const colorModeCookie = useCookie<'light' | 'dark'>('color-mode', {
    default: () => 'light',
    maxAge: 60 * 60 * 24 * 30,
    path: '/',
  })

  const colorMode = useState<'light' | 'dark'>('color-mode', () => colorModeCookie.value)

  /** 应用主题到 DOM（仅客户端） */
  function applyTheme(theme: 'light' | 'dark') {
    if (import.meta.server) return
    const html = document.documentElement
    if (theme === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  /** 初始化主题（从 cookie 同步，客户端调用） */
  function initTheme() {
    if (import.meta.server) return
    // cookie 中已有服务端设置的值，直接应用
    applyTheme(colorMode.value)
  }

  /** 切换主题 */
  async function toggleTheme() {
    const next: 'dark' | 'light' = colorMode.value === 'dark' ? 'light' : 'dark'
    colorMode.value = next
    colorModeCookie.value = next

    // 应用到 DOM
    applyTheme(next)

    // 异步存入 Redis（不阻塞 UI）
    try {
      await $fetch('/api/theme', { method: 'POST', body: { theme: next } })
    } catch {
      // 失败静默，不影响切换体验
    }
  }

  const isDark = computed(() => colorMode.value === 'dark')

  return {
    colorMode,
    isDark,
    toggleTheme,
    initTheme,
  }
}
