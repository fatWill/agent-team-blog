package tuyile

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/fatWill/agent-team-blog/backend/pkg/middleware"
)

const (
	// ShareHashLen 分享 hash 长度（base62，62^12 ≈ 3.2e21，不可枚举）
	ShareHashLen = 12
	// ShareMaxMessages 单个分享最多消息条数
	ShareMaxMessages = 50
	// ShareMaxMessageRunes 单条消息最大字符数
	ShareMaxMessageRunes = 4000
	// ShareMaxBodyBytes 请求体上限（同时约束总体积）
	ShareMaxBodyBytes = 128 << 10
	base62            = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	shareLimiter *ipLimiter
	shareSiteURL = "https://fatwill.cn"
	reShareHash  = regexp.MustCompile(`^[0-9A-Za-z]{12}$`)
	// randReader 可在测试中替换
	randReader = rand.Reader
)

func initShare(cfg config.TuyileConfig) {
	perMin, perDay := cfg.ShareRatePerMinute, cfg.ShareRatePerDay
	if perMin <= 0 {
		perMin = 5
	}
	if perDay <= 0 {
		perDay = 30
	}
	shareLimiter = newIPLimiter(perMin, perDay)
	if u := strings.TrimRight(strings.TrimSpace(cfg.SiteURL), "/"); u != "" {
		shareSiteURL = u
	}
}

// newShareHash crypto/rand 生成 base62 随机串（拒绝采样，无取模偏差）
func newShareHash() (string, error) {
	b := make([]byte, ShareHashLen)
	max := big.NewInt(int64(len(base62)))
	for i := range b {
		n, err := rand.Int(randReader, max)
		if err != nil {
			return "", err
		}
		b[i] = base62[n.Int64()]
	}
	return string(b), nil
}

type shareInput struct {
	Messages []struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	} `json:"messages"`
}

// validateShare 校验并清洗：只保留 user/assistant 纯文本，至少 1 条 assistant
func validateShare(in shareInput) ([]models.TuyileShareMessage, error) {
	if len(in.Messages) == 0 {
		return nil, errors.New("messages 不能为空")
	}
	if len(in.Messages) > ShareMaxMessages {
		return nil, errors.New("最多分享 50 条消息")
	}
	out := make([]models.TuyileShareMessage, 0, len(in.Messages))
	hasAssistant := false
	for _, m := range in.Messages {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role != "user" && role != "assistant" {
			continue
		}
		text := strings.TrimSpace(inputMessage{Content: m.Content}.text())
		if text == "" {
			continue
		}
		if !utf8.ValidString(text) {
			return nil, errors.New("消息包含非法字符")
		}
		if utf8.RuneCountInString(text) > ShareMaxMessageRunes {
			return nil, errors.New("单条消息最多 4000 字")
		}
		if role == "assistant" {
			hasAssistant = true
		}
		out = append(out, models.TuyileShareMessage{Role: role, Content: text})
	}
	if !hasAssistant {
		return nil, errors.New("至少需要 1 条图一乐的回复才能分享")
	}
	return out, nil
}

// CreateShare POST /api/toys/tuyile/share — 用户主动分享时保存对话快照
func CreateShare(c *gin.Context) {
	if shareLimiter == nil {
		reject(c, http.StatusServiceUnavailable, "服务未初始化")
		return
	}
	if ok, reason := shareLimiter.allow(middleware.GetClientIP(c)); !ok {
		c.Header("Retry-After", ternary(reason == "day", "3600", "60"))
		reject(c, http.StatusTooManyRequests, ternary(reason == "day", "今天分享得够多啦，明天再来～", "分享太频繁啦，歇一分钟再来～"))
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, ShareMaxBodyBytes)
	var in shareInput
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			reject(c, http.StatusRequestEntityTooLarge, "分享内容超过 128KB")
			return
		}
		reject(c, http.StatusBadRequest, "请求体不是合法 JSON")
		return
	}
	msgs, err := validateShare(in)
	if err != nil {
		reject(c, http.StatusBadRequest, err.Error())
		return
	}

	// 唯一索引兜底：极小概率碰撞时重试
	var share models.TuyileShare
	for i := 0; i < 3; i++ {
		hash, herr := newShareHash()
		if herr != nil {
			reject(c, http.StatusInternalServerError, "生成分享链接失败")
			return
		}
		share = models.TuyileShare{Hash: hash, Messages: msgs, CreatedAt: nowFn()}
		if err = db.DB.Create(&share).Error; err == nil {
			break
		}
		if !strings.Contains(strings.ToLower(err.Error()), "unique") {
			break
		}
	}
	if err != nil {
		logf("tuyile.share.create.failed err=%v", err)
		reject(c, http.StatusInternalServerError, "保存分享失败")
		return
	}
	// 日志不记录内容，只记录条数
	logf("tuyile.share.create.success hash=%s msgs=%d", share.Hash, len(msgs))
	c.JSON(http.StatusOK, gin.H{"hash": share.Hash, "url": shareSiteURL + "/toys/tuyile/share/" + share.Hash})
}

// GetShare GET /api/toys/tuyile/share/:hash — 只读获取分享快照
func GetShare(c *gin.Context) {
	hash := c.Param("hash")
	if !reShareHash.MatchString(hash) {
		reject(c, http.StatusNotFound, "分享不存在或已失效")
		return
	}
	var list []models.TuyileShare
	if err := db.DB.Where("hash = ?", hash).Limit(1).Find(&list).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		reject(c, http.StatusInternalServerError, "读取分享失败")
		return
	}
	if len(list) == 0 {
		reject(c, http.StatusNotFound, "分享不存在或已失效")
		return
	}
	// 快照不可变：允许 CDN/浏览器缓存
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, list[0])
}
