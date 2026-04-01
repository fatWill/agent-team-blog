import { message } from 'ant-design-vue'
import 'ant-design-vue/es/message/style'
import 'ant-design-vue/es/modal/style'
import 'ant-design-vue/es/spin/style'

export default defineNuxtPlugin(() => {
  // 配置 message 全局设置
  message.config({
    top: '60px',
    duration: 2,
    maxCount: 3,
  })
})
