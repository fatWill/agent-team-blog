import NProgress from 'nprogress'

export default defineNuxtPlugin((nuxtApp) => {
  // NProgress 配置
  NProgress.configure({
    showSpinner: false,   // 不显示右上角转圈
    speed: 400,
    minimum: 0.15,
    easing: 'ease',
  })

  // 路由开始 → 启动进度条
  nuxtApp.hook('page:start', () => {
    NProgress.start()
  })

  // 路由结束 → 完成进度条
  nuxtApp.hook('page:finish', () => {
    NProgress.done()
  })

  // 路由报错也要结束
  nuxtApp.hook('app:error', () => {
    NProgress.done()
  })
})
