/**
 * 修复 pinia 2.3.x SSR hydration 报错：
 * TypeError: obj.hasOwnProperty is not a function
 *
 * 原因：pinia 的 shouldHydrate 使用 obj.hasOwnProperty() 而非
 * Object.prototype.hasOwnProperty.call(obj)，当 SSR payload 中包含
 * 无原型对象（如 mysql2 JSON 解析结果）时崩溃。
 *
 * 修复方式：在 payload 序列化前，用 JSON 往返清洗所有 useAsyncData 数据，
 * 确保所有对象都是纯 Object 实例。
 */
export default defineNuxtPlugin({
  name: 'pinia-hydration-fix',
  enforce: 'post',
  setup(nuxtApp) {
    // 仅服务端需要处理
    if (!import.meta.server) return

    // 在 payload 渲染前清洗数据，确保所有对象都有 hasOwnProperty 方法
    nuxtApp.hook('app:rendered', () => {
      const payload = nuxtApp.payload
      if (payload?.data) {
        // 清洗 useAsyncData 返回的数据
        for (const key of Object.keys(payload.data)) {
          try {
            payload.data[key] = JSON.parse(JSON.stringify(payload.data[key]))
          } catch {
            // 序列化失败的数据保持原样
          }
        }
      }
      // 清洗 pinia state
      if (payload?.pinia) {
        try {
          payload.pinia = JSON.parse(JSON.stringify(payload.pinia))
        } catch {
          // 序列化失败保持原样
        }
      }
    })
  },
})
