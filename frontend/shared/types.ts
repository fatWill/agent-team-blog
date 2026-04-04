/** 导航项 */
export interface NavItem {
  label: string
  to: string
}

/** Tab 项 */
export interface TabItem {
  key: string
  label: string
}

/** 博主个人资料 */
export interface Profile {
  avatar: string
  bio: string
}
