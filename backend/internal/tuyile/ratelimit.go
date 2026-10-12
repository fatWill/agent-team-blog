package tuyile

import (
	"sync"
	"time"

	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
)

// ipLimiter 内存双窗口限流：每分钟固定窗口 + 北京时间自然日计数。
//
// 只在单进程内生效（与现有 middleware.RateLimit 一致，不依赖 Redis），
// 2C2G 单实例部署足够；条目每 10 分钟清理一次，内存占用与活跃 IP 数线性相关。
type ipLimiter struct {
	mu        sync.Mutex
	perMin    int
	perDay    int
	entries   map[string]*ipEntry
	lastSweep time.Time
	now       func() time.Time
}

type ipEntry struct {
	minStart time.Time
	minCount int
	day      string
	dayCount int
}

func newIPLimiter(perMin, perDay int) *ipLimiter {
	return &ipLimiter{perMin: perMin, perDay: perDay, entries: map[string]*ipEntry{}, now: time.Now}
}

// allow 返回是否放行；不放行时 reason 为 "minute" 或 "day"
func (l *ipLimiter) allow(ip string) (bool, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	day := now.In(algo.CST).Format(algo.DateLayout)
	if now.Sub(l.lastSweep) > 10*time.Minute {
		for k, e := range l.entries {
			if e.day != day && now.Sub(e.minStart) > time.Minute {
				delete(l.entries, k)
			}
		}
		l.lastSweep = now
	}
	e, ok := l.entries[ip]
	if !ok {
		e = &ipEntry{minStart: now, day: day}
		l.entries[ip] = e
	}
	if e.day != day {
		e.day, e.dayCount = day, 0
	}
	if now.Sub(e.minStart) >= time.Minute {
		e.minStart, e.minCount = now, 0
	}
	if e.dayCount >= l.perDay {
		return false, "day"
	}
	if e.minCount >= l.perMin {
		return false, "minute"
	}
	e.minCount++
	e.dayCount++
	return true, ""
}
