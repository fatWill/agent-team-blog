# 图一乐 - 组件清单

## TuyileHome

- **路径**: `features/tuyile/components/TuyileHome.vue`
- **职责**: 3 块入口（星座运势、生辰八字、老黄历），前两块显示「即将上线」，点击后弹 toast 提示
- **Emits**: `open('almanac')`

## AlmanacView

- **路径**: `features/tuyile/components/AlmanacView.vue`
- **职责**: 老黄历主视图，包括日期条（前一天/后一天、日历、回到今天）以及公历农历、宜忌、日神吉凶、彭祖百忌、吉神凶神、方位胎神、十二时辰和「聊一下你的运势」按钮
- **Props**: `date: string`（YYYY-MM-DD）、`today: string`
- **Emits**: `update:date(string)`、`chat()`

## AlmanacCalendar

- **路径**: `features/tuyile/components/AlmanacCalendar.vue`
- **职责**: 月视图（周一起始），日格副标题优先显示 `festivals[0]`，没有节日时显示 `cellText`；黄道日带金点；按月缓存
- **Props**: `modelValue: string`、`today: string`
- **Emits**: `update:modelValue(string)`

## AlmanacHours

- **路径**: `features/tuyile/components/AlmanacHours.vue`
- **职责**: 十二时辰吉凶网格，标出当前时辰（仅当天），点击展开该时辰的宜忌、冲煞和神煞
- **Props**: `hours: AlmanacHour[]`、`currentIndex: number`（-1 表示不标）

## FortuneChat

- **路径**: `features/tuyile/components/FortuneChat.vue`
- **职责**: AG-UI 运势聊天窗口，支持流式 Markdown、500 字限制计数、停止生成、错误和限流提示，并显示「仅供娱乐」
- **Emits**: `back()`
