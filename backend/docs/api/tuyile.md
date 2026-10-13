# 玩具「图一乐」接口文档

前端路由：`/toys/tuyile`、`/toys/tuyile/share/<hash>`。所有接口均**公开、无需登录**。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/toys/tuyile/almanac` | 单日老黄历（完整字段） |
| GET | `/api/toys/tuyile/almanac/month` | 月度老黄历（日历视图精简字段） |
| POST | `/api/toys/tuyile/chat` | AG-UI 协议运势聊天（SSE 流式） |
| POST | `/api/toys/tuyile/share` | 创建对话分享快照（v2.17.0） |
| GET | `/api/toys/tuyile/share/:hash` | 读取对话分享（只读，v2.17.0） |

数据来源：`almanac_days` 表（huangli.com 抓取 + `6tail/lunar-go` 算法交叉校验）。库里没有的日期（如 2028 年以后）由算法现算返回，`source=lunar-go`。支持范围为 1900-01-01 至 2100-12-31。

---

## [GET] /api/toys/tuyile/almanac - 单日黄历

### 请求参数（Query）
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| date | string | ❌ | 北京时间今天 | `YYYY-MM-DD` |

### 响应头
`Cache-Control: public, max-age=600`（不传 date 时为 `max-age=60`，避免跨零点后还拿到旧数据）

### 成功响应 (200)
```json
{
  "date": "2024-01-01",
  "week": "星期一",
  "lunarYear": 2023,
  "lunarMonth": 11,
  "lunarDay": 20,
  "isLeapMonth": false,
  "lunarMonthCn": "冬月",
  "lunarDayCn": "二十",
  "lunarDate": "冬月二十",
  "lunarYearCn": "二〇二三",
  "zodiacYear": "兔",
  "ganZhiYear": "癸卯",
  "ganZhiMonth": "甲子",
  "ganZhiDay": "甲子",
  "jieQi": "",
  "currentJieQi": "冬至",
  "festivals": ["元旦节"],
  "yi": ["交易", "进人口", "祭祀", "沐浴", "捕捉", "入殓", "除服", "成服", "安葬"],
  "ji": ["斋醮", "入宅", "修造", "动土", "破土"],
  "chongZodiac": "马",
  "chongGanZhi": "戊午",
  "chongAge": 46,
  "chongDesc": "冲马(戊午)46岁 煞南",
  "sha": "南",
  "naYinYear": "金箔金",
  "naYinMonth": "海中金",
  "naYinDay": "海中金",
  "zhiShen": "金匮",
  "zhiShenType": "黄道",
  "jianChu": "建",
  "xiu": "毕月乌",
  "xiuLuck": "吉",
  "xiuGong": "西方白虎",
  "pengZu": "甲不开仓财物耗散 子不问卜自惹祸殃",
  "jiShen": ["天恩", "天赦", "月恩", "四相", "官日", "敬安", "金匮"],
  "xiongSha": ["月建", "小时", "土府", "月厌", "地火"],
  "taiShen": "占门碓 外东南",
  "xiShen": "东北",
  "caiShen": "东北",
  "fuShen": "正北",
  "yangGui": "西南",
  "yinGui": "东北",
  "liuYao": "赤口",
  "wuHou": "水泉动",
  "yueXiang": "更待月",
  "hours": [
    {
      "zhi": "子",
      "ganZhi": "甲子",
      "range": "23:00-00:59",
      "tianShen": "金匮",
      "luck": "吉",
      "chong": "马",
      "sha": "南",
      "yi": ["祈福", "订婚", "嫁娶", "开市", "安葬"],
      "ji": ["造船", "乘船"],
      "jiShen": ["金匮", "比肩", "大进", "福德"],
      "xiongSha": ["日建"]
    }
  ],
  "source": "huangli.com+lunar-go",
  "algoFields": ["lunarYear", "lunarYearCn", "zodiacYear", "jieQi", "currentJieQi", "festivals", "chongGanZhi", "chongAge", "naYinYear", "naYinMonth", "zhiShenType", "xiuGong", "xiShen", "caiShen", "fuShen", "yangGui", "yinGui", "liuYao", "wuHou", "yueXiang"],
  "mismatches": [],
  "updatedAt": "2026-10-12T21:00:01.123+08:00"
}
```
> `hours` 固定 12 项（子→亥），示例只列出第 1 项。

### 字段说明
| 字段 | 类型 | 说明 |
|------|------|------|
| date | string | 公历日期 `YYYY-MM-DD` |
| week | string | 星期，如 `星期一` |
| lunarYear / lunarMonth / lunarDay | int | 农历年/月/日（数字）。年以春节为界 |
| isLeapMonth | bool | 是否闰月（闰月时 `lunarMonth` 仍为正数） |
| lunarMonthCn / lunarDayCn | string | 农历月名（`冬月`、`闰六月`）/ 日名（`初一`、`廿八`） |
| lunarDate | string | `lunarMonthCn + lunarDayCn` |
| lunarYearCn | string | 农历年中文数字，如 `二〇二三` |
| zodiacYear | string | 生肖年（春节为界） |
| ganZhiYear / ganZhiMonth / ganZhiDay | string | 年/月/日柱。年柱以春节为界，月柱以节为界 |
| jieQi | string | 当天交接的节气，不交节气则为空 |
| currentJieQi | string | 当前所处的节气区间 |
| festivals | string[] | 节日（公历、农历及其他），可能为空数组 |
| yi / ji | string[] | 当日宜 / 忌 |
| chongZodiac | string | 冲的生肖 |
| chongGanZhi | string | 冲的干支 |
| chongAge | int | 被冲者当年虚岁（按农历年算） |
| chongDesc | string | 冲煞展示文案，如 `冲马(戊午)46岁 煞南` |
| sha | string | 煞方：东/南/西/北 |
| naYinYear / naYinMonth / naYinDay | string | 年/月/日纳音，其中日纳音即「五行」 |
| zhiShen / zhiShenType | string | 值神，以及 `黄道` / `黑道` |
| jianChu | string | 建除十二神：建除满平定执破危成收开闭 |
| xiu / xiuLuck / xiuGong | string | 二十八星宿（如 `毕月乌`）/ 吉凶 / 所属四象（如 `西方白虎`） |
| pengZu | string | 彭祖百忌（天干句 + 地支句，空格分隔） |
| jiShen / xiongSha | string[] | 吉神宜趋 / 凶神宜忌 |
| taiShen | string | 胎神占方 |
| xiShen / caiShen / fuShen | string | 喜神 / 财神 / 福神方位 |
| yangGui / yinGui | string | 阳贵神 / 阴贵神方位 |
| liuYao / wuHou / yueXiang | string | 六曜 / 物候 / 月相 |
| hours[] | object[] | 12 时辰吉凶，见下表 |
| source | string | `huangli.com+lunar-go`：站点抓取，经算法校验补齐；`lunar-go`：纯算法 |
| algoFields | string[] | 由算法补齐的字段名。纯算法记录为 `["*"]` |
| mismatches | string[] | 站点与算法不一致的项（确定性字段以算法为准），正常情况为空 |
| updatedAt | string | 记录更新时间（纯算法现算的记录为零值时间） |

`hours[]` 字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| zhi | string | 地支：子…亥 |
| ganZhi | string | 时柱干支 |
| range | string | 时间范围，子时为 `23:00-00:59` |
| tianShen | string | 时辰星神 |
| luck | string | `吉` / `凶` |
| chong | string | 冲生肖 |
| sha | string | 煞方 |
| yi / ji | string[] | 时宜 / 时忌（已去掉「无」） |
| jiShen / xiongSha | string[] | 时辰吉神 / 凶煞，可能为空数组 |

### 错误响应
| HTTP | statusMessage | 场景 |
|------|---------------|------|
| 400 | 日期格式应为 YYYY-MM-DD | date 格式非法 |
| 400 | 日期超出支持范围（1900-2100） | 年份越界 |
| 500 | 查询黄历失败 | 数据库异常 |

---

## [GET] /api/toys/tuyile/almanac/month - 月度黄历

### 请求参数（Query）
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| month | string | ❌ | 北京时间当月 | `YYYY-MM`（必须两位月份） |

### 成功响应 (200)
```json
{
  "month": "2024-02",
  "year": 2024,
  "days": 29,
  "list": [
    {
      "date": "2024-02-04",
      "week": "星期日",
      "lunarDate": "腊月廿五",
      "cellText": "立春",
      "ganZhiDay": "戊戌",
      "jieQi": "立春",
      "festivals": [],
      "zhiShenType": "黄道",
      "yiTop": ["祭祀", "祈福", "求嗣"],
      "jiTop": ["开光", "嫁娶"],
      "chongZodiac": "龙",
      "source": "huangli.com+lunar-go"
    },
    {
      "date": "2024-02-10",
      "week": "星期六",
      "lunarDate": "正月初一",
      "cellText": "正月",
      "ganZhiDay": "甲辰",
      "jieQi": "",
      "festivals": ["春节"],
      "zhiShenType": "黄道",
      "yiTop": ["嫁娶", "开光", "求嗣"],
      "jiTop": ["入宅", "移徙", "出火"],
      "chongZodiac": "狗",
      "source": "huangli.com+lunar-go"
    }
  ]
}
```
> `list` 按日期升序，包含该月每一天（示例只列 2 天）。

### 字段说明
| 字段 | 类型 | 说明 |
|------|------|------|
| month / year / days | string / int / int | 查询月份、年份、当月天数 |
| list[].cellText | string | 日格副标题：有节气显示节气，初一显示月名，其余显示农历日名 |
| list[].yiTop / jiTop | string[] | 宜 / 忌前 3 项 |
| list[].zhiShenType | string | `黄道` / `黑道`，可用于标记吉日 |
| 其余字段 | — | 含义同单日接口 |

节日优先展示：前端可在 `festivals` 非空时用 `festivals[0]` 覆盖 `cellText`。

### 错误响应
| HTTP | statusMessage | 场景 |
|------|---------------|------|
| 400 | month 格式应为 YYYY-MM | 格式非法或越界 |

---

## [POST] /api/toys/tuyile/chat - 运势聊天（AG-UI）

遵循 [AG-UI HTTP + SSE 绑定](https://docs.ag-ui.com/spec/1.0/basic/transports/http-sse)：POST 一个 `RunAgentInput` JSON，响应为 `text/event-stream`，每个 `data:` 行是**一个**完整事件 JSON（LF 换行）。可直接用 `@ag-ui/client` 的 `HttpAgent` 对接。

> 二期（v2.17.0）变更：上游切换为智谱 GLM（`glm-5.3-flash`）；新增 Reasoning + `CUSTOM tuyile.status` 状态事件；`query_almanac` 工具在服务端闭环；`messages: []` 开场稳定为图一乐人设；不再强制收集姓名 / 生日 / 年份。

### 请求头
| Header | 必填 | 说明 |
|--------|------|------|
| Content-Type | ✅ | `application/json` |
| Accept | 建议 | `text/event-stream` |

### 请求体（RunAgentInput）
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| threadId | string | 建议 | 会话 ID（≤128 字符），缺省或非法时服务端生成 |
| runId | string | 建议 | 本次运行 ID（≤128 字符），缺省时服务端生成 |
| messages | Message[] | ✅ | 完整对话历史。传 `[]` 表示开场，图一乐会主动打招呼 |
| messages[].role | string | ✅ | 只处理 `user` / `assistant`；`system` / `developer` / `tool` / `reasoning` 会被丢弃 |
| messages[].content | string 或 ContentPart[] | ✅ | 文本；ContentPart 只取 `type=text` |
| tools / state / context / forwardedProps | any | ❌ | 可传，服务端忽略 |

服务端无状态、不落库：**前端要自行保存并每次带上完整 messages**（只需 user/assistant 文本，**不要回传 reasoning 消息**）。服务端只保留最近 20 条。

```json
{
  "threadId": "thr-1",
  "runId": "run-3",
  "messages": [
    { "id": "a0", "role": "assistant", "content": "哟，来客了！本仙乃大罗天上下凡……先报个网名？" },
    { "id": "u1", "role": "user", "content": "我叫熬夜冠军" },
    { "id": "a1", "role": "assistant", "content": "「熬夜冠军」，这是把职业写在头像上了啊……" },
    { "id": "u2", "role": "user", "content": "下周哪天适合搬家" }
  ]
}
```

### 事件总览
| 事件 | 何时出现 | 前端处理建议 |
|------|----------|--------------|
| `RUN_STARTED` | 第一条 | 记录 threadId / runId |
| `REASONING_START` / `REASONING_END` | 思考阶段的 span，**收到 RUN_STARTED 后立刻打开**，第一个正文字到达前关闭 | 可忽略，或用来控制「思考中」状态条显隐 |
| `REASONING_MESSAGE_START/CONTENT/END` | span 内，每切换一个阶段一条，`delta` 是**状态摘要**（如「掐指一算中…」「翻黄历中…」），**不是思考原文** | 用 `@ag-ui/client` 时会生成 `role: "reasoning"` 消息，**不要渲染成气泡**，也不要回传给服务端 |
| `CUSTOM`（`name: "tuyile.status"`） | 和上面的阶段同步 | **推荐用它来驱动 UI**：直接展示 `value.text` |
| `TEXT_MESSAGE_START/CONTENT/END` | 正文，一个 run 只有一条 assistant 消息 | 逐字追加 `delta`（Markdown，用安全模式渲染） |
| `RUN_FINISHED` / `RUN_ERROR` | 最后一条，二选一 | 结束 loading |
| `: keep-alive`（SSE 注释） | 每 15 秒一次 | 忽略 |

`CUSTOM tuyile.status` 的 `value`：

| 字段 | 类型 | 说明 |
|------|------|------|
| stage | string | `thinking`：模型思考中；`querying`：服务端正在查黄历（工具调用） |
| text | string | 前端可直接展示的文案：`掐指一算中…` / `翻黄历中…` |
| tool | string | 仅 `querying` 时出现，固定为 `query_almanac` |
| range | string | 仅 `querying` 时出现，查询范围，如 `2028-02-02~2028-02-16` 或 `2026-10-14` |

为什么选 Reasoning 事件 + CUSTOM：AG-UI 1.0 用 `REASONING_*` 取代了旧的 `THINKING_*`，并明确允许「只给摘要」的可见性。我们只下发状态摘要，模型的思考原文不出服务端。`CUSTOM` 用来给前端一个结构化、稳定的状态信号，不依赖 SDK 对 reasoning 消息的处理方式。

**工具调用对前端透明**：`query_almanac` 在服务端闭环执行（模型调工具 → 服务端查库 → 结果回填 → 模型继续作答），**不透出 `TOOL_CALL_*` 事件**，只通过 `CUSTOM stage=querying` 告诉前端「翻黄历中…」。前端不需要实现任何工具。

### 成功事件流示例（含工具调用，本地实测 glm-5.3-flash，正文 delta 已省略中间部分）
```
data: {"type":"RUN_STARTED","timestamp":1791893082375,"threadId":"thr-demo","runId":"run-demo"}

data: {"type":"REASONING_START","timestamp":1791893082375,"messageId":"rsn_1d2eb32d-a735-4693-9492-f3f1edbff50b"}

data: {"type":"REASONING_MESSAGE_START","timestamp":1791893082375,"messageId":"rsm_4f1e25a3-aa18-4269-aa8c-e6082d7d23e6","role":"reasoning"}

data: {"type":"REASONING_MESSAGE_CONTENT","timestamp":1791893082375,"messageId":"rsm_4f1e25a3-aa18-4269-aa8c-e6082d7d23e6","delta":"掐指一算中…"}

data: {"type":"CUSTOM","timestamp":1791893082375,"name":"tuyile.status","value":{"stage":"thinking","text":"掐指一算中…"}}

data: {"type":"REASONING_MESSAGE_END","timestamp":1791893087917,"messageId":"rsm_4f1e25a3-aa18-4269-aa8c-e6082d7d23e6"}

data: {"type":"REASONING_MESSAGE_START","timestamp":1791893087918,"messageId":"rsm_1c82cbd4-6504-4302-9556-fceb4151ca2f","role":"reasoning"}

data: {"type":"REASONING_MESSAGE_CONTENT","timestamp":1791893087918,"messageId":"rsm_1c82cbd4-6504-4302-9556-fceb4151ca2f","delta":"翻黄历中…"}

data: {"type":"CUSTOM","timestamp":1791893087918,"name":"tuyile.status","value":{"stage":"querying","text":"翻黄历中…","tool":"query_almanac","range":"2028-02-02~2028-02-16"}}

data: {"type":"REASONING_MESSAGE_END","timestamp":1791893087934,"messageId":"rsm_1c82cbd4-6504-4302-9556-fceb4151ca2f"}

data: {"type":"REASONING_MESSAGE_START","timestamp":1791893087934,"messageId":"rsm_56433054-b6a4-4d27-bc21-1588b9bd9f78","role":"reasoning"}

data: {"type":"REASONING_MESSAGE_CONTENT","timestamp":1791893087934,"messageId":"rsm_56433054-b6a4-4d27-bc21-1588b9bd9f78","delta":"掐指一算中…"}

data: {"type":"CUSTOM","timestamp":1791893087934,"name":"tuyile.status","value":{"stage":"thinking","text":"掐指一算中…"}}

data: {"type":"REASONING_MESSAGE_END","timestamp":1791893091698,"messageId":"rsm_56433054-b6a4-4d27-bc21-1588b9bd9f78"}

data: {"type":"REASONING_END","timestamp":1791893091699,"messageId":"rsn_1d2eb32d-a735-4693-9492-f3f1edbff50b"}

data: {"type":"TEXT_MESSAGE_START","timestamp":1791893091699,"messageId":"msg_95ece3ca-9131-402d-85eb-3af5ebc2a57b","role":"assistant"}

data: {"type":"TEXT_MESSAGE_CONTENT","timestamp":1791893091699,"messageId":"msg_95ece3ca-9131-402d-85eb-3af5ebc2a57b","delta":"施"}

data: {"type":"TEXT_MESSAGE_CONTENT","timestamp":1791893091699,"messageId":"msg_95ece3ca-9131-402d-85eb-3af5ebc2a57b","delta":"主"}

: keep-alive

data: {"type":"TEXT_MESSAGE_END","timestamp":1791893098636,"messageId":"msg_95ece3ca-9131-402d-85eb-3af5ebc2a57b"}

data: {"type":"RUN_FINISHED","timestamp":1791893098636,"threadId":"thr-demo","runId":"run-demo"}
```
不需要查库时（大多数问题的日期已在知识库注入里），流程是：`RUN_STARTED → REASONING_START → REASONING_MESSAGE_START → REASONING_MESSAGE_CONTENT → CUSTOM(thinking) → REASONING_MESSAGE_END → REASONING_END → TEXT_MESSAGE_START → TEXT_MESSAGE_CONTENT×N → TEXT_MESSAGE_END → RUN_FINISHED`。

**时延**（本地实测 glm-5.3-flash，`reasoning_effort=low`）：首个状态事件与 `RUN_STARTED` 同一毫秒发出（不等上游）；开场首字 0.65~1.9s；带黄历注入的问题首字 3.7~8.3s；一次工具调用总耗时约 17s。

### 运行内错误（200 + RUN_ERROR，流随即关闭）
出错前会先关闭所有已打开的 reasoning span / message 和 text message，再发 `RUN_ERROR`。

未配置 key：
```
data: {"type":"RUN_STARTED","timestamp":1791810096322,"threadId":"thr-1","runId":"run-1"}

data: {"type":"RUN_ERROR","timestamp":1791810096322,"message":"图一乐还没接上大模型（服务端未配置 LLM_API_KEY），暂时算不了，稍后再来～","code":"LLM_NOT_CONFIGURED"}
```
超时（本地把超时调小后实测）：
```
data: {"type":"RUN_STARTED","timestamp":1791893112763,"threadId":"thr-to","runId":"run-to"}

data: {"type":"REASONING_START","timestamp":1791893112763,"messageId":"rsn_8c9ba956-b273-453f-ab8b-c65125b992ca"}

data: {"type":"REASONING_MESSAGE_START","timestamp":1791893112763,"messageId":"rsm_33f95db5-91aa-4f44-a442-a0d0d92bc653","role":"reasoning"}

data: {"type":"REASONING_MESSAGE_CONTENT","timestamp":1791893112763,"messageId":"rsm_33f95db5-91aa-4f44-a442-a0d0d92bc653","delta":"掐指一算中…"}

data: {"type":"CUSTOM","timestamp":1791893112763,"name":"tuyile.status","value":{"stage":"thinking","text":"掐指一算中…"}}

data: {"type":"REASONING_MESSAGE_END","timestamp":1791893113908,"messageId":"rsm_33f95db5-91aa-4f44-a442-a0d0d92bc653"}

data: {"type":"REASONING_END","timestamp":1791893113908,"messageId":"rsn_8c9ba956-b273-453f-ab8b-c65125b992ca"}

data: {"type":"RUN_ERROR","timestamp":1791893113908,"message":"图一乐掐指算得太久了（大模型响应超时），换个简单点的问法再试试～","code":"LLM_TIMEOUT"}
```

| code | 场景 |
|------|------|
| LLM_NOT_CONFIGURED | 未配置 `LLM_API_KEY` / `LLM_BASE_URL` / `LLM_MODEL`；使用了内网地址；`LLM_EXTRA_BODY` 不是合法 JSON |
| LLM_TIMEOUT | 等响应头超时（`LLM_HEADER_TIMEOUT_SECONDS`）、单轮总超时（`LLM_TIMEOUT_SECONDS`）、网络读超时、上游返回 504 / 408 |
| LLM_UNAVAILABLE | 连接失败（拒绝连接、DNS、TLS），**不含超时** |
| LLM_STREAM_BROKEN | 流读到一半被中断（不是超时） |
| LLM_AUTH_FAILED | 上游返回 401 / 403 |
| LLM_RATE_LIMITED | 上游返回 429 |
| LLM_UPSTREAM_ERROR | 上游返回其他错误（例如 400 参数错误） |
| LLM_EMPTY_RESPONSE | 上游没有返回正文 |
| INTERNAL_ERROR | 其他内部错误 |

### 运行前拒绝（HTTP 错误，无事件流）
`{"error":true,"statusCode":<code>,"statusMessage":"..."}`

| HTTP | statusMessage | 场景 |
|------|---------------|------|
| 400 | 请求体不是合法的 RunAgentInput JSON | JSON 非法，或请求体超过 64KB |
| 400 | 单条消息最多 500 字 | 最后一条 user 消息超长 |
| 400 | 最后一条消息必须是 user | messages 非空且最后一条不是 user |
| 429 | 算得太勤快啦，歇一分钟再来～ | 单 IP 每分钟超限（默认 10），带 `Retry-After: 60` |
| 429 | 今天的卦已经算满啦，明天再来图一乐～ | 单 IP 每日超限（默认 100，北京时间自然日），带 `Retry-After: 3600` |
| 503 | 图一乐小摊排队中，稍后再试～ | 全局并发流已满（默认 8），带 `Retry-After: 5` |

### 防刷与隐私
- IP 取 `X-Forwarded-For` 第一段 → `X-Real-IP` → 连接地址；计数存在进程内存里，重启后清零。
- 最后一条 user 消息不超过 500 字；历史只保留 20 条；历史中的 user 消息截断到 500 字、assistant 截断到 1200 字。
- 对话内容不落库。日志只记录 runId、消息条数、输出字数、思考字数、工具调用次数、首字耗时、是否识别到生日、预测年份、耗时和错误码，**不记录原文和思考内容**。

### 对话逻辑（System Prompt 要点）
- 人设是大罗天仙级算命先生兼职脱口秀：算命为主、段子点缀，不说客服腔。
- 什么都能算：择日（搬家、入职、开业、签约、表白、考试…）、取名测名、网名测性格，以及整活题（如「吉时屁」，会按方位和时辰一本正经地分析）。
- 会主动开话题（问网名或问最近想算啥）；**不强制收集**姓名 / 生日 / 年份，识别到了就作为参考注入。
- 干支、宜忌、时辰、方位等以知识库和工具返回为准；择日只推荐「宜」里明确包含该事项的日子。
- 给出结论的回复，结尾加一句俏皮的娱乐声明；涉及医疗 / 投资 / 法律时幽默地拒绝，并建议找专业人士。
- **开场**：`messages: []` 时服务端额外注入开场专用指令和一句占位发言，开场白稳定为图一乐人设（5/5 实测通过）。占位发言不会回显给前端。

### 黄历知识注入与工具（第 4 点修复）
**根因**：一期的 `buildKnowledge` 只注入「今天」和「用户写出的具体日期」。像「下周」「下个月」「明年三月」「明天」这类相对日期什么都没注入，模型只看到今天的数据，就回答「以后没数据」。

**二期方案**（三层叠加）：
1. **时间锚点**：每次都注入当前北京时间、本周和下周的起止日期、本月和下个月，模型换算相对日期有据可依。
2. **相对日期解析**（服务端规则，`reldate.go`）：支持今天/明天/后天/大后天、周X/下周X、本周/下周/下下周/周末、未来几天/半个月、本月/月底/下个月/下下个月、`YYYY年M月`、明年/后年 X 月、X 月上/中/下旬（没写年份时，已经过去的月份按明年算）。解析出的单日注入完整信息含 12 时辰；范围注入每日精简宜忌，单个范围最多 62 天。如果问句里有择日意图（搬家、开业、结婚、出行、签约…），会按「宜」预先筛选，与工具同一口径。另外始终注入「今天含 12 时辰」和「未来 7 天精简表」。
3. **`query_almanac` 工具**（服务端 function calling）：参数为 `start`、`end`（可选，闭区间，最多 62 天，超出截断）和 `keywords`（可选，口语词自动映射为黄历宜项，如 搬家→移徙/入宅、开业→开市）。单日返回完整信息加时辰，范围返回每日一行精简信息并标注 `✅命中`。每个 run 最多 3 轮工具调用。模型遇到注入里没有的日期（如 2028 年）会自己调用。已实测 glm-5.3-flash 支持 tools + stream。

兜底：设置 `LLM_DISABLE_TOOLS=1`（上游不支持 tools 时用）后不再下发工具，改为注入未来 **30** 天的精简表，相对日期解析照常生效。

**覆盖**：`almanac_days` 库内有 1461 天，即 2024-01-01 ~ 2027-12-31，全部为 `huangli.com+lunar-go`。库外日期（1900~2100）由 `lunar-go` 算法现算，`source=lunar-go`：干支、农历、冲煞、值神、星宿、时辰、方位等确定性字段精确，宜忌、吉神凶神用算法流派的值（和站点可能略有差异），不回写数据库。

---

## [POST] /api/toys/tuyile/share - 创建对话分享

只有用户主动点「分享」时才保存快照，**只读**，不能基于分享继续对话。

### 请求体
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| messages | {role, content}[] | ✅ | 1~50 条；只保留 `user` / `assistant`（其他 role 和空内容会被丢弃），**至少包含 1 条 assistant** |
| messages[].content | string 或 ContentPart[] | ✅ | 纯文本，单条不超过 4000 字（按字符算） |

请求体总大小不超过 128KB。

```json
{
  "messages": [
    { "role": "user", "content": "下周哪天适合搬家" },
    { "role": "assistant", "content": "**10月24日（周六）** 宜移徙……\n\n🎲 本卦仅供娱乐，图一乐～" }
  ]
}
```

### 成功响应 (200)
```json
{
  "hash": "gaYD0IYtJIhM",
  "url": "https://fatwill.cn/toys/tuyile/share/gaYD0IYtJIhM"
}
```
| 字段 | 说明 |
|------|------|
| hash | 12 位 base62（`[0-9A-Za-z]`），由 `crypto/rand` 生成（约 71 bit 熵），不可猜测 |
| url | 分享页地址，`SITE_URL + /toys/tuyile/share/<hash>` |

### 错误响应
| HTTP | statusMessage | 场景 |
|------|---------------|------|
| 400 | 请求体不是合法 JSON | JSON 非法 |
| 400 | messages 不能为空 | 空数组 |
| 400 | 最多分享 50 条消息 | 超过 50 条 |
| 400 | 单条消息最多 4000 字 | 单条超长 |
| 400 | 至少需要 1 条图一乐的回复才能分享 | 清洗后没有 assistant 消息 |
| 413 | 分享内容超过 128KB | 请求体超限 |
| 429 | 分享太频繁啦，歇一分钟再来～ | 单 IP 每分钟超过 5 次，带 `Retry-After: 60` |
| 429 | 今天分享得够多啦，明天再来～ | 单 IP 每天超过 30 次，带 `Retry-After: 3600` |
| 500 | 保存分享失败 / 生成分享链接失败 | 服务端异常 |

## [GET] /api/toys/tuyile/share/:hash - 读取分享

### 成功响应 (200)
响应头 `Cache-Control: public, max-age=3600`（快照不可变）。
```json
{
  "hash": "gaYD0IYtJIhM",
  "messages": [
    { "role": "user", "content": "下周哪天适合搬家<script>x</script>" },
    { "role": "assistant", "content": "**10月24日（周六）** 宜移徙\n\n🎲 本卦仅供娱乐，图一乐～" }
  ],
  "createdAt": "2026-10-13T20:04:34.173859+08:00"
}
```
`content` **原样存储为纯文本**（服务端不转义、也不渲染）。JSON 里 `<` `>` 会被编码为 `\u003c` `\u003e`，解析后还原为原文。**前端必须用安全模式渲染 Markdown**：禁用原始 HTML，或用 DOMPurify 之类的工具消毒，不能用 `v-html` 直接输出原文。

### 错误响应
| HTTP | statusMessage | 场景 |
|------|---------------|------|
| 404 | 分享不存在或已失效 | hash 不存在，或格式不是 12 位 base62 |
| 500 | 读取分享失败 | 数据库异常 |

---

## 微信分享卡片（复用 `GET /api/wechat/jssdk-config`）

不新增签名接口。现有接口按 JS-SDK 规则用 `jsapi_ticket + noncestr + timestamp + url` 做 SHA1 签名，**url 由调用方传入，不限路径**，所以 `/toys/tuyile` 和 `/toys/tuyile/share/<hash>` 都能签。前提是页面域名已在公众号后台的「JS 接口安全域名」里配置（`fatwill.cn`，和现有文章页同域，已在用）。

### 请求
`GET /api/wechat/jssdk-config?url=<encodeURIComponent(当前页完整 URL，去掉 # 及之后的部分)>`

| 参数 | 必填 | 说明 |
|------|------|------|
| url | ✅ | **必须和浏览器当前地址完全一致**：协议、域名、路径、query 都要一致，只去掉 `#` 之后的部分。分享页是 `https://fatwill.cn/toys/tuyile/share/gaYD0IYtJIhM`；聊天页有 query 时（如 `?view=chat`）要带上 |

### 响应 (200)
```json
{ "appId": "wx1234567890abcdef", "timestamp": 1791893200, "nonceStr": "Ab3dEf6hIj9kLm2n", "signature": "3f1c...e9" }
```
错误：400「缺少 url 参数」；500「微信公众号未配置」或「获取微信签名失败」。

### 前端调用要点
- 直接复用 `frontend/composables/useWechatShare.ts` 的 `initWechatShare({ title, desc, link, imgUrl })`，它已经按 `location.href.split('#')[0]` 签名，并调用 `updateAppMessageShareData` / `updateTimelineShareData`。
- **跨路径**：签名和 URL 一一对应。SPA 内从 `/toys/tuyile` 跳到 `/toys/tuyile/share/<hash>`（或反过来）后，要**在新路由上重新调用一次** `initWechatShare`，否则 `wx.config` 会报 invalid signature。iOS 微信对 SPA 的 pushState 有「以首次进入的 URL 为准」的老问题；如果 iOS 出现签名失败，可以对首次进入的 URL 再签一次作为回退。
- 分享卡片的 `link` 填分享页 URL（创建分享接口返回的 `url`），域名必须和安全域名一致；`imgUrl` 用绝对地址的 https 图片。

---

## 数据导入（运维）

- **表**：`almanac_days`（服务启动时 autoMigrate 自动建表，存量库也会自动补建）
- **生产库**：`DB_PATH`，默认 `/root/blog-data/blog.db`
- **推荐方式**：本地生成数据，上传后导入。服务器 2C2G，跑爬虫要 30 分钟；本地生成后导入，生产侧完全不用联网抓取
- **兜底**：即使表为空，接口也会用 lunar-go 现算返回（`source=lunar-go`），缺失宜忌 / 吉神等站点数据时用算法流派值补齐

```bash
# 1) 本地：抓取 + 交叉校验（≤1 req/s，约 30 分钟；中断后重跑同一命令即断点续爬）
cd backend
go run ./cmd/almanac-crawler --start 2024-01-01 --end 2027-12-31 --db ./data/almanac.db

# 2) 本地：交叉编译爬虫/导入工具
GOOS=linux GOARCH=amd64 go build -o almanac-crawler ./cmd/almanac-crawler

# 3) 上传（数据文件约 7MB）
scp -i ~/.ssh/id_ed25519_vps ./data/almanac.db almanac-crawler root@<server>:/root/blog-data/

# 4) 服务器：先部署 v2.16.0 新二进制并 restart（启动时自动建 almanac_days 表、备份库），再导入
#    导入为单事务 upsert by date，可重复执行；SQLite WAL + busy_timeout，服务运行中可以直接导入
cd /root/blog-data && ./almanac-crawler --import-db ./almanac.db --db /root/blog-data/blog.db
#    输出「库内覆盖 1461/1461 天（缺失 0）」即成功

# 5) 验证
curl -s 'http://127.0.0.1:8080/api/toys/tuyile/almanac?date=2024-02-10' | head -c 300
```

也可以直接在服务器上抓取：`./almanac-crawler --start ... --end ... --db /root/blog-data/blog.db`。爬虫遵守 robots，robots 不可达时自动降级为纯算法模式。

## 二期上线（v2.17.0）

无需导入数据，`tuyile_shares` 表在服务启动时自动建立。步骤：
1. 构建并替换二进制（流程同一期）。
2. 在 systemd `Environment` 中修改或新增以下配置：`LLM_BASE_URL=https://open.bigmodel.cn/api/paas/v4`、`LLM_MODEL=glm-5.3-flash`、`LLM_API_KEY=<GLM_API_KEY>`、`LLM_REASONING_EFFORT=low`；删除旧的 `LLM_TIMEOUT_SECONDS=90`，或改为 120。
3. `systemctl daemon-reload && systemctl restart blog-backend`。
4. 验证：`journalctl -u blog-backend` 中出现 `tuyile.llm.proto=HTTP/2.0`（首次聊天后打印）；`curl -sN -X POST http://127.0.0.1:8080/api/toys/tuyile/chat -H 'Content-Type: application/json' -d '{"messages":[]}'` 能收到 `CUSTOM` 与 `TEXT_MESSAGE_*`；`sqlite3 /root/blog-data/blog.db ".schema tuyile_shares"` 有输出。
5. Nginx：`/api/toys/tuyile/share*` 走现有的 `location /api/`，无需修改；chat 仍走一期新增的 `location = /api/toys/tuyile/chat`（关闭缓冲）。
