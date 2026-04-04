/**
 * 统一 UI 工具
 * 自动根据设备类型（移动端/PC端）选择对应的 UI 组件
 * - PC 端：ant-design-vue 的 message / Modal
 * - 移动端：自定义轻量级 MobileToast / MobileDialog
 */
import { message, Modal } from 'ant-design-vue'
import { isMobileDevice } from '~/composables/useDevice'
import { MobileToast, MobileDialog } from '~/shared/utils/mobileUI'

/** 成功提示 */
export function showSuccess(content: string) {
  if (isMobileDevice()) {
    MobileToast.success(content)
  } else {
    message.success(content)
  }
}

/** 错误提示 */
export function showError(content: string) {
  if (isMobileDevice()) {
    MobileToast.error(content)
  } else {
    message.error(content)
  }
}

/** 普通信息提示 */
export function showInfo(content: string) {
  if (isMobileDevice()) {
    MobileToast.show({ content })
  } else {
    message.info(content)
  }
}

/** 确认弹窗（删除等危险操作） */
export function showConfirm(options: {
  title: string
  content: string
  okText?: string
  cancelText?: string
  danger?: boolean
  onOk: () => void | Promise<void>
}) {
  if (isMobileDevice()) {
    MobileDialog.confirm({
      title: options.title,
      content: options.content,
      okText: options.okText || '确定',
      cancelText: options.cancelText || '取消',
      danger: options.danger,
      onOk: options.onOk,
    })
  } else {
    Modal.confirm({
      title: options.title,
      content: options.content,
      okText: options.okText || '确定',
      okType: options.danger ? 'danger' : 'primary',
      cancelText: options.cancelText || '取消',
      onOk: options.onOk,
    })
  }
}
