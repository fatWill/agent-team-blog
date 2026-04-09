/**
 * Web Vitals 性能指标采集 & 上报插件（仅客户端执行）
 *
 * 采集指标：FCP、LCP、TTFB、DOM Ready、Load Time、JS 加载耗时
 * 首次页面加载：采集全部 Navigation Timing + FCP/LCP
 * SPA 路由切换：仅上报 loadTime（路由切换耗时）+ page
 */
export default defineNuxtPlugin(() => {
  const router = useRouter()

  let fcpValue: number | null = null
  let lcpValue: number | null = null
  let hasReportedInitial = false

  // ---- 设备类型判断 ----
  function getDeviceType(): string {
    const ua = navigator.userAgent
    if (/Mobi|Android|iPhone/i.test(ua)) return 'mobile'
    if (/iPad|Tablet/i.test(ua)) return 'tablet'
    return 'desktop'
  }

  // ---- 静默上报 ----
  async function report(data: Record<string, unknown>) {
    try {
      await $fetch('/api/perf/report', { method: 'POST', body: data })
    } catch {
      // 静默失败，不影响用户体验
    }
  }

  // ---- 监听 FCP ----
  if ('PerformanceObserver' in window) {
    try {
      const fcpObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          if (entry.name === 'first-contentful-paint') {
            fcpValue = Math.round(entry.startTime * 10) / 10
          }
        }
      })
      fcpObserver.observe({ type: 'paint', buffered: true })
    } catch {
      // 浏览器不支持，忽略
    }

    // ---- 监听 LCP ----
    try {
      const lcpObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries()
        // LCP 会持续更新，取最后一个
        if (entries.length > 0) {
          const last = entries[entries.length - 1]
          lcpValue = Math.round(last.startTime * 10) / 10
        }
      })
      lcpObserver.observe({ type: 'largest-contentful-paint', buffered: true })
    } catch {
      // 浏览器不支持，忽略
    }
  }

  // ---- 首次页面加载采集 ----
  function collectInitialMetrics() {
    if (hasReportedInitial) return
    hasReportedInitial = true

    const navEntries = performance.getEntriesByType('navigation') as PerformanceNavigationTiming[]
    const nav = navEntries[0]

    let ttfb: number | null = null
    let domReady: number | null = null
    let loadTime: number | null = null

    if (nav) {
      ttfb = Math.round(nav.responseStart * 10) / 10
      domReady = Math.round(nav.domContentLoadedEventEnd * 10) / 10
      loadTime = Math.round(nav.loadEventEnd * 10) / 10
    }

    // JS 资源加载：取 script 类型中 duration 最大值
    let jsLoad: number | null = null
    const resources = performance.getEntriesByType('resource') as PerformanceResourceTiming[]
    const scripts = resources.filter(r => r.initiatorType === 'script')
    if (scripts.length > 0) {
      jsLoad = Math.round(Math.max(...scripts.map(s => s.duration)) * 10) / 10
    }

    report({
      page: router.currentRoute.value.path,
      fcp: fcpValue,
      lcp: lcpValue,
      ttfb,
      domReady,
      loadTime,
      jsLoad,
      deviceType: getDeviceType(),
    })
  }

  // window.load + 3秒延迟，等待 LCP 稳定
  if (document.readyState === 'complete') {
    setTimeout(collectInitialMetrics, 3000)
  } else {
    window.addEventListener('load', () => {
      setTimeout(collectInitialMetrics, 3000)
    })
  }

  // ---- SPA 路由切换：记录路由耗时 ----
  let routeStartTime = 0

  router.beforeEach(() => {
    routeStartTime = Date.now()
  })

  router.afterEach((to) => {
    // 跳过首次加载（已由 collectInitialMetrics 处理）
    if (!hasReportedInitial) return

    // 延迟 3 秒上报，等待页面渲染完成
    const startTime = routeStartTime
    setTimeout(() => {
      const routeLoadTime = startTime > 0 ? Date.now() - startTime - 3000 : null
      // routeLoadTime 减去延迟的 3 秒，得到实际渲染耗时的近似值
      // 如果结果为负数说明页面很快就渲染完了，设为一个较小的正值
      const adjustedLoadTime = routeLoadTime !== null
        ? Math.max(routeLoadTime, 0)
        : null

      report({
        page: to.path,
        fcp: null,
        lcp: null,
        ttfb: null,
        domReady: null,
        loadTime: adjustedLoadTime,
        jsLoad: null,
        deviceType: getDeviceType(),
      })
    }, 3000)
  })
})
