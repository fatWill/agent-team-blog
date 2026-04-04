/**
 * 移动端统一 UI 工具
 * 提供 Toast（成功/错误提示）和 Dialog（确认弹窗）
 * 风格参照 antd-mobile 设计语言
 */
import { createApp, defineComponent, h, ref, Transition } from 'vue'

// ==================== Toast ====================
let toastTimer: ReturnType<typeof setTimeout> | null = null
let toastContainer: HTMLElement | null = null

interface ToastOptions {
  content: string
  icon?: 'success' | 'error' | 'loading'
  duration?: number // 毫秒，默认 2000
}

function createToastIcon(icon: string) {
  if (icon === 'success') {
    return h('svg', { class: 'mobile-toast-icon', viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2.5' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M5 13l4 4L19 7' }),
    ])
  }
  if (icon === 'error') {
    return h('svg', { class: 'mobile-toast-icon', viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2.5' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M6 18L18 6M6 6l12 12' }),
    ])
  }
  if (icon === 'loading') {
    return h('svg', { class: 'mobile-toast-icon mobile-toast-spin', viewBox: '0 0 24 24', fill: 'none' }, [
      h('circle', { class: 'opacity-25', cx: '12', cy: '12', r: '10', stroke: 'currentColor', 'stroke-width': '4' }),
      h('path', { class: 'opacity-75', fill: 'currentColor', d: 'M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z' }),
    ])
  }
  return null
}

const ToastComponent = defineComponent({
  props: {
    content: { type: String, required: true },
    icon: { type: String, default: '' },
  },
  setup(props) {
    const visible = ref(false)
    setTimeout(() => { visible.value = true }, 10)

    return () =>
      h(Transition, { name: 'mobile-toast' }, {
        default: () =>
          visible.value
            ? h('div', { class: 'mobile-toast-mask' }, [
                h('div', { class: 'mobile-toast-content' }, [
                  props.icon ? createToastIcon(props.icon) : null,
                  h('span', { class: 'mobile-toast-text' }, props.content),
                ]),
              ])
            : null,
      })
  },
})

function ensureToastStyles() {
  if (document.getElementById('mobile-toast-styles')) return
  const style = document.createElement('style')
  style.id = 'mobile-toast-styles'
  style.textContent = `
    .mobile-toast-mask {
      position: fixed;
      inset: 0;
      z-index: 10000;
      display: flex;
      align-items: center;
      justify-content: center;
      pointer-events: none;
    }
    .mobile-toast-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      padding: 16px 24px;
      background: rgba(0,0,0,0.75);
      border-radius: 12px;
      color: #fff;
      backdrop-filter: blur(8px);
      max-width: 70vw;
      pointer-events: auto;
    }
    .mobile-toast-icon {
      width: 36px;
      height: 36px;
    }
    .mobile-toast-text {
      font-size: 14px;
      line-height: 1.4;
      text-align: center;
      word-break: break-word;
    }
    .mobile-toast-spin {
      animation: mobile-toast-rotate 1s linear infinite;
    }
    @keyframes mobile-toast-rotate {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }
    .mobile-toast-enter-active { transition: opacity 0.2s ease; }
    .mobile-toast-leave-active { transition: opacity 0.15s ease; }
    .mobile-toast-enter-from, .mobile-toast-leave-to { opacity: 0; }

    /* Dialog 样式 */
    .mobile-dialog-overlay {
      position: fixed;
      inset: 0;
      z-index: 10001;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(0,0,0,0.55);
    }
    .mobile-dialog-box {
      width: 280px;
      background: #fff;
      border-radius: 16px;
      overflow: hidden;
    }
    .mobile-dialog-body {
      padding: 28px 20px 16px;
      text-align: center;
    }
    .mobile-dialog-title {
      font-size: 17px;
      font-weight: 600;
      color: #333;
      margin-bottom: 8px;
    }
    .mobile-dialog-message {
      font-size: 14px;
      color: #666;
      line-height: 1.5;
    }
    .mobile-dialog-footer {
      display: flex;
      border-top: 1px solid #eee;
    }
    .mobile-dialog-btn {
      flex: 1;
      padding: 12px 0;
      font-size: 16px;
      font-weight: 500;
      border: none;
      background: transparent;
      cursor: pointer;
      -webkit-tap-highlight-color: transparent;
    }
    .mobile-dialog-btn:active {
      background: #f5f5f5;
    }
    .mobile-dialog-btn + .mobile-dialog-btn {
      border-left: 1px solid #eee;
    }
    .mobile-dialog-btn-cancel {
      color: #666;
    }
    .mobile-dialog-btn-ok {
      color: #1677ff;
    }
    .mobile-dialog-btn-danger {
      color: #ff4d4f;
    }
    /* loading 样式 */
    .mobile-loading-wrap {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 12px;
      padding: 20px 0;
    }
    .mobile-loading-spinner {
      width: 32px;
      height: 32px;
      color: #1677ff;
      animation: mobile-toast-rotate 0.8s linear infinite;
    }
    .mobile-loading-tip {
      font-size: 14px;
      color: #999;
    }
    @media (prefers-color-scheme: dark) {
      .mobile-dialog-box { background: #1f1f1f; }
      .mobile-dialog-title { color: #e5e5e5; }
      .mobile-dialog-message { color: #999; }
      .mobile-dialog-footer { border-top-color: #333; }
      .mobile-dialog-btn:active { background: #2a2a2a; }
      .mobile-dialog-btn + .mobile-dialog-btn { border-left-color: #333; }
      .mobile-dialog-btn-cancel { color: #999; }
    }
  `
  document.head.appendChild(style)
}

export const MobileToast = {
  show(options: ToastOptions) {
    MobileToast.clear()
    ensureToastStyles()

    toastContainer = document.createElement('div')
    document.body.appendChild(toastContainer)

    const app = createApp(ToastComponent, {
      content: options.content,
      icon: options.icon || '',
    })
    app.mount(toastContainer)

    const duration = options.duration ?? 2000
    if (duration > 0) {
      toastTimer = setTimeout(() => {
        MobileToast.clear()
      }, duration)
    }
  },
  success(content: string) {
    MobileToast.show({ content, icon: 'success' })
  },
  error(content: string) {
    MobileToast.show({ content, icon: 'error' })
  },
  loading(content: string) {
    MobileToast.show({ content, icon: 'loading', duration: 0 })
  },
  clear() {
    if (toastTimer) { clearTimeout(toastTimer); toastTimer = null }
    if (toastContainer) {
      toastContainer.remove()
      toastContainer = null
    }
  },
}

// ==================== Dialog ====================
interface DialogOptions {
  title?: string
  content: string
  okText?: string
  cancelText?: string
  danger?: boolean
  onOk?: () => void | Promise<void>
  onCancel?: () => void
}

export const MobileDialog = {
  confirm(options: DialogOptions) {
    ensureToastStyles()

    const container = document.createElement('div')
    document.body.appendChild(container)

    function destroy() {
      container.remove()
    }

    const DialogComp = defineComponent({
      setup() {
        const loading = ref(false)

        async function handleOk() {
          if (options.onOk) {
            loading.value = true
            try {
              await options.onOk()
            } finally {
              loading.value = false
            }
          }
          destroy()
        }

        function handleCancel() {
          options.onCancel?.()
          destroy()
        }

        return () =>
          h('div', { class: 'mobile-dialog-overlay', onClick: (e: MouseEvent) => { if (e.target === e.currentTarget) handleCancel() } }, [
            h('div', { class: 'mobile-dialog-box' }, [
              h('div', { class: 'mobile-dialog-body' }, [
                options.title ? h('div', { class: 'mobile-dialog-title' }, options.title) : null,
                h('div', { class: 'mobile-dialog-message' }, options.content),
              ]),
              h('div', { class: 'mobile-dialog-footer' }, [
                h('button', {
                  class: 'mobile-dialog-btn mobile-dialog-btn-cancel',
                  onClick: handleCancel,
                  disabled: loading.value,
                }, options.cancelText || '取消'),
                h('button', {
                  class: `mobile-dialog-btn ${options.danger ? 'mobile-dialog-btn-danger' : 'mobile-dialog-btn-ok'}`,
                  onClick: handleOk,
                  disabled: loading.value,
                }, loading.value ? '处理中...' : (options.okText || '确定')),
              ]),
            ]),
          ])
      },
    })

    createApp(DialogComp).mount(container)
  },

  alert(options: Omit<DialogOptions, 'cancelText' | 'onCancel'>) {
    ensureToastStyles()

    const container = document.createElement('div')
    document.body.appendChild(container)

    function destroy() {
      container.remove()
    }

    const AlertComp = defineComponent({
      setup() {
        function handleOk() {
          options.onOk?.()
          destroy()
        }

        return () =>
          h('div', { class: 'mobile-dialog-overlay' }, [
            h('div', { class: 'mobile-dialog-box' }, [
              h('div', { class: 'mobile-dialog-body' }, [
                options.title ? h('div', { class: 'mobile-dialog-title' }, options.title) : null,
                h('div', { class: 'mobile-dialog-message' }, options.content),
              ]),
              h('div', { class: 'mobile-dialog-footer' }, [
                h('button', {
                  class: 'mobile-dialog-btn mobile-dialog-btn-ok',
                  onClick: handleOk,
                }, options.okText || '我知道了'),
              ]),
            ]),
          ])
      },
    })

    createApp(AlertComp).mount(container)
  },
}
