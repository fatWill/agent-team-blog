/**
 * 设备检测 composable re-export（向后兼容层）
 * 实际实现已移至 shared/composables/useDevice.ts
 * 保留在 composables/ 目录以确保 Nuxt 自动导入正常工作
 */
export { useDevice, isMobileDevice } from '~/shared/composables/useDevice'
