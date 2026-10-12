# 玩具「图一乐」接口文档

前端路由：`/toys/tuyile`。三个接口均**公开、无需登录**。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/toys/tuyile/almanac` | 单日老黄历（完整字段） |
| GET | `/api/toys/tuyile/almanac/month` | 月度老黄历（日历视图精简字段） |
| POST | `/api/toys/tuyile/chat` | AG-UI 协议运势聊天（SSE 流式） |

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

遵循 [AG-UI HTTP + SSE 绑定](https://docs.ag-ui.com/spec/1.0/basic/transports/http-sse)：POST 一个 `RunAgentInput` JSON，响应为 `text/event-stream`，每个 `data:` 行是**一个**完整事件 JSON，换行符为 LF。可直接用 `@ag-ui/client` 的 `HttpAgent` 对接。

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
| messages | Message[] | ✅ | 完整对话历史。可传 `[]`，此时图一乐主动打招呼 |
| messages[].id | string | ❌ | 消息 ID |
| messages[].role | string | ✅ | 只处理 `user` / `assistant`，`system`/`developer`/`tool` 等会被丢弃 |
| messages[].content | string 或 ContentPart[] | ✅ | 文本；ContentPart 只取 `type=text` |
| tools / state / context / forwardedProps | any | ❌ | 可传，服务端忽略 |

服务端无状态、不落库，**前端需自行保存并每次带上完整 messages**（服务端只保留最近 20 条）。

```json
{
  "threadId": "thr-1",
  "runId": "run-2",
  "messages": [
    { "id": "u1", "role": "user", "content": "你好" },
    { "id": "a1", "role": "assistant", "content": "嘿嘿，欢迎光临图一乐卦摊～先报上大名呗？" },
    { "id": "u2", "role": "user", "content": "我叫阿乐" },
    { "id": "a2", "role": "assistant", "content": "阿乐你好！生日是哪天？公历还是农历？时辰知道的话也可以说～" },
    { "id": "u3", "role": "user", "content": "农历1990年五月初三，晚上9点" },
    { "id": "a3", "role": "assistant", "content": "收到！想看哪一年的运势？" },
    { "id": "u4", "role": "user", "content": "2027年运势" }
  ]
}
```

### 成功事件流 (200, text/event-stream)
```
data: {"type":"RUN_STARTED","timestamp":1791810113882,"threadId":"thr-1","runId":"run-2"}

data: {"type":"TEXT_MESSAGE_START","timestamp":1791810113890,"messageId":"msg_50655da1-e172-420e-87f9-55be70c1e6d5","role":"assistant"}

data: {"type":"TEXT_MESSAGE_CONTENT","timestamp":1791810113890,"messageId":"msg_50655da1-e172-420e-87f9-55be70c1e6d5","delta":"🌟 整体运势："}

data: {"type":"TEXT_MESSAGE_CONTENT","timestamp":1791810114043,"messageId":"msg_50655da1-e172-420e-87f9-55be70c1e6d5","delta":"流年六合，顺风顺水。"}

data: {"type":"TEXT_MESSAGE_CONTENT","timestamp":1791810114193,"messageId":"msg_50655da1-e172-420e-87f9-55be70c1e6d5","delta":"\n🎲 以上仅供娱乐，图一乐～"}

data: {"type":"TEXT_MESSAGE_END","timestamp":1791810114346,"messageId":"msg_50655da1-e172-420e-87f9-55be70c1e6d5"}

data: {"type":"RUN_FINISHED","timestamp":1791810114346,"threadId":"thr-1","runId":"run-2"}
```

首 token 较慢时，服务端每 15 秒发送一行 SSE 注释 `: keep-alive`。按 AG-UI 规范，消费方会忽略这类注释。

### 运行内错误（200 + RUN_ERROR，流随即关闭）
```
data: {"type":"RUN_STARTED","timestamp":1791810096322,"threadId":"thr-1","runId":"run-1"}

data: {"type":"RUN_ERROR","timestamp":1791810096322,"message":"图一乐还没接上大模型（服务端未配置 LLM_API_KEY），暂时算不了，稍后再来～","code":"LLM_NOT_CONFIGURED"}
```
如果错误发生在输出过程中，会先发 `TEXT_MESSAGE_END`，再发 `RUN_ERROR`。

| code | 场景 |
|------|------|
| LLM_NOT_CONFIGURED | 未配置 `LLM_API_KEY` / `LLM_BASE_URL` / `LLM_MODEL`，或使用了内网地址 |
| LLM_UNAVAILABLE | 上游连接失败 |
| LLM_TIMEOUT | 上游超时（`LLM_TIMEOUT_SECONDS`） |
| LLM_AUTH_FAILED | 上游返回 401 / 403（key 错误） |
| LLM_RATE_LIMITED | 上游返回 429 |
| LLM_UPSTREAM_ERROR | 上游返回其他错误 |
| LLM_STREAM_BROKEN | 上游流中断 |
| LLM_EMPTY_RESPONSE | 上游没有返回任何内容 |
| INTERNAL_ERROR | 其他内部错误 |

### 运行前拒绝（HTTP 错误，无事件流）
按 AG-UI 绑定约定，运行开始前的拒绝返回普通 JSON：`{"error":true,"statusCode":<code>,"statusMessage":"..."}`

| HTTP | statusMessage | 场景 |
|------|---------------|------|
| 400 | 请求体不是合法的 RunAgentInput JSON | JSON 非法，或请求体超过 64KB |
| 400 | 单条消息最多 500 字 | 最后一条 user 消息超长 |
| 400 | 最后一条消息必须是 user | messages 非空且末尾不是 user |
| 429 | 算得太勤快啦，歇一分钟再来～ | 单 IP 超过每分钟限额，带 `Retry-After: 60` |
| 429 | 今天的卦已经算满啦，明天再来图一乐～ | 单 IP 超过每日限额（北京时间自然日），带 `Retry-After: 3600` |
| 503 | 图一乐小摊排队中，稍后再试～ | 全局并发流已满，带 `Retry-After: 5` |

### 防刷与隐私
- **IP 限流**：默认每分钟 10 次、每天 100 次。IP 取 `X-Forwarded-For` 首段 → `X-Real-IP` → 连接地址。计数存在进程内存，重启后清零。
- **全局并发**：同时最多 8 条流（`TUYILE_MAX_CONCURRENCY`），保护 2C2G 服务器。
- **长度限制**：最后一条 user 消息不超过 500 字；历史只保留最近 20 条（约 10 轮）；历史 user 消息截断到 500 字，assistant 消息截断到 1200 字。
- **不落库**：姓名、生日、对话内容都不写数据库；日志只记录 runId、消息条数、输出字数、是否识别到生日、预测年份、耗时、错误码。

### 对话逻辑（System Prompt 要点）
人设「图一乐」，轻松幽默。引导顺序为：**姓名 → 出生年月日（区分公历/农历，时辰可选）→ 预测年份**。三项收齐前只追问，不给结论；收齐后输出 整体 / 事业 / 财运 / 感情 / 健康 / 宜忌建议，结尾附「🎲 以上仅供娱乐，图一乐～」。拒绝医疗、投资、法律、彩票类建议。

每次请求，服务端从 user 消息中抽取以下信息并注入知识库：
- 今天的黄历；
- 用户提到的日期（最多 3 个），支持 `2026-10-01`、`2026/10/1`、`2026年10月1日`、`农历1990年五月初三`、`农历2025年闰6月初一`；
- 出生日：换算为公历，并给出八字（立春口径）、纳音、生肖（春节口径），以及出生当天的黄历；
- 流年：干支、生肖、纳音，以及与用户生肖的冲合刑害破关系。

库里没有的日期由 `lunar-go` 现算。

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
