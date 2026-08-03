/**
 * Sandbox iframe 兼容插件（最高优先级，所有 client plugin 之前执行）
 *
 * 当页面被嵌入 sandbox="allow-scripts"（无 allow-same-origin）的 iframe 时，
 * document.cookie / localStorage / sessionStorage 的访问会抛 SecurityError。
 * 此插件在最早期检测并用无害的 polyfill 替换，避免后续代码崩溃。
 */
export default defineNuxtPlugin({
  name: 'sandbox-safe',
  enforce: 'pre',
  setup() {
    // 仅客户端执行（文件名带 .client.ts 已保证，这里 double check）
    if (import.meta.server) return

    // ---- 兜底 document.cookie ----
    try {
      void document.cookie
    } catch {
      Object.defineProperty(Document.prototype, 'cookie', {
        get() { return '' },
        set() { /* noop */ },
        configurable: true,
      })
    }

    // ---- 兜底 localStorage / sessionStorage ----
    function makeMemoryStorage(): Storage {
      let store: Record<string, string> = {}
      return {
        getItem(k: string) { return k in store ? store[k] : null },
        setItem(k: string, v: string) { store[k] = String(v) },
        removeItem(k: string) { delete store[k] },
        clear() { store = {} },
        key(i: number) { return Object.keys(store)[i] || null },
        get length() { return Object.keys(store).length },
      } as Storage
    }

    ;(['localStorage', 'sessionStorage'] as const).forEach((name) => {
      try {
        void (window[name] as Storage).length
      } catch {
        Object.defineProperty(window, name, {
          value: makeMemoryStorage(),
          configurable: true,
        })
      }
    })
  },
})
