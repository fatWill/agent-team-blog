/**
 * 主题切换组合式函数
 * 支持 dark/light 模式，默认跟随系统，持久化到 localStorage
 */
export function useTheme() {
  const colorMode = useState<'light' | 'dark'>('color-mode', () => 'light')

  /** 初始化主题（仅客户端） */
  function initTheme() {
    if (import.meta.server) return

    const stored = localStorage.getItem('color-mode')
    if (stored === 'dark' || stored === 'light') {
      colorMode.value = stored
    } else {
      // 跟随系统
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      colorMode.value = prefersDark ? 'dark' : 'light'
    }
    applyTheme()
  }

  /** 切换主题 */
  function toggleTheme() {
    colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'
    if (import.meta.client) {
      localStorage.setItem('color-mode', colorMode.value)
    }
    applyTheme()
  }

  /** 应用主题到 DOM */
  function applyTheme() {
    if (import.meta.server) return
    const html = document.documentElement
    if (colorMode.value === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
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
