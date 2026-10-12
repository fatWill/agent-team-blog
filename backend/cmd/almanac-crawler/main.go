// almanac-crawler 老黄历爬虫：抓取 huangli.com 日页，经 lunar-go 交叉校验后幂等写入 SQLite almanac_days 表。
//
// 用法：
//
//	DB_PATH=./data/almanac.db go run ./cmd/almanac-crawler --start 2024-01-01 --end 2027-12-31
//
// 特性：
//   - 遵守 robots.txt（启动时拉取并校验目标路径）
//   - 全局限速，默认 1.2s/次（不超过 1 req/s）
//   - 单日失败指数退避重试；最终失败时用算法兜底写入（source=lunar-go），下次运行会自动重抓
//   - 断点续爬：已以 merged 来源入库的日期直接跳过（--force 强制重抓）
//   - 按 date 唯一键 upsert，可重复执行
//   - --import-db 导入模式：把本地生成的库按 date upsert 到目标库（线上推荐用法）
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/internal/almanac"
	"github.com/fatWill/agent-team-blog/backend/models"
	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
)

const userAgent = "Mozilla/5.0 (compatible; fatwill-almanac-bot/1.0; +https://fatwill.cn)"

type stats struct {
	total, skipped, crawled, fallback, mismatchDays, mismatchFields int
}

func main() {
	start := flag.String("start", "2024-01-01", "开始日期 YYYY-MM-DD")
	end := flag.String("end", "2027-12-31", "结束日期 YYYY-MM-DD（含）")
	interval := flag.Duration("interval", 1200*time.Millisecond, "请求最小间隔（不得小于 1s）")
	retries := flag.Int("retries", 3, "单日最大重试次数")
	force := flag.Bool("force", false, "忽略断点，强制重抓已入库日期")
	offline := flag.Bool("offline", false, "离线模式：不访问网络，全部用 lunar-go 生成")
	dbPath := flag.String("db", "", "SQLite 路径（缺省读 DB_PATH 环境变量）")
	importDB := flag.String("import-db", "", "导入模式：从该 SQLite 文件读取 almanac_days 并按 date upsert 到 --db（不访问网络）")
	flag.Parse()

	if *interval < time.Second {
		log.Fatalf("interval 不得小于 1s（限速约束），当前 %s", *interval)
	}
	st, err := algo.ParseDate(*start)
	if err != nil {
		log.Fatalf("start 非法: %v", err)
	}
	et, err := algo.ParseDate(*end)
	if err != nil {
		log.Fatalf("end 非法: %v", err)
	}
	if et.Before(st) {
		log.Fatalf("end 早于 start")
	}

	cfg := config.Load()
	if *dbPath != "" {
		cfg.DB.Path = *dbPath
	}
	if err := db.Init(&cfg.DB); err != nil {
		log.Fatalf("初始化 SQLite 失败: %v", err)
	}
	log.Printf("db=%s range=%s~%s interval=%s offline=%v", cfg.DB.Path, *start, *end, *interval, *offline)
	repo := almanac.NewRepo(db.DB)

	if *importDB != "" {
		if err := importFrom(*importDB, repo, *start, *end); err != nil {
			log.Fatalf("导入失败: %v", err)
		}
		report(repo, *start, *end, st, et, stats{})
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client := &http.Client{Timeout: 15 * time.Second}
	allowed := true
	if !*offline {
		allowed, err = robotsAllows(ctx, client, "https://www.huangli.com", "/huangli/")
		if err != nil {
			log.Printf("⚠️ robots.txt 获取失败，保守起见切换离线模式: %v", err)
			allowed = false
		}
		if !allowed {
			log.Printf("⚠️ robots.txt 不允许抓取 /huangli/，仅用算法生成")
		}
	}
	useSite := !*offline && allowed

	done, err := repo.SourceMap(*start, *end)
	if err != nil {
		log.Fatalf("读取断点失败: %v", err)
	}

	var s stats
	var lastReq time.Time
	for t := st; !t.After(et); t = t.AddDate(0, 0, 1) {
		if ctx.Err() != nil {
			log.Printf("收到中断信号，已安全停止；下次运行将从断点继续")
			break
		}
		s.total++
		date := t.Format(algo.DateLayout)
		src, exists := done[date]
		if exists && !*force && (src == models.AlmanacSourceMerged || !useSite) {
			s.skipped++
			continue
		}

		calc := algo.Compute(t)
		day := calc
		if useSite {
			site, err := fetchWithRetry(ctx, client, date, *retries, *interval, &lastReq)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					break
				}
				log.Printf("❌ %s 抓取失败，算法兜底: %v", date, err)
				s.fallback++
			} else {
				var mm []string
				day, mm = algo.Merge(calc, site)
				s.crawled++
				if len(mm) > 0 {
					s.mismatchDays++
					s.mismatchFields += len(mm)
					log.Printf("⚠️ %s 交叉校验不一致: %s", date, strings.Join(mm, "; "))
				}
			}
		}
		if err := repo.Upsert(day); err != nil {
			log.Fatalf("写库失败 %s: %v", date, err)
		}
		if s.total%30 == 0 {
			log.Printf("进度 %s total=%d crawled=%d skipped=%d fallback=%d mismatchDays=%d", date, s.total, s.crawled, s.skipped, s.fallback, s.mismatchDays)
		}
	}

	report(repo, *start, *end, st, et, s)
}

// fetchWithRetry 限速 + 指数退避重试
func fetchWithRetry(ctx context.Context, client *http.Client, date string, retries int, interval time.Duration, lastReq *time.Time) (*algo.SiteDay, error) {
	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<attempt) * time.Second
			if err := sleepCtx(ctx, backoff); err != nil {
				return nil, err
			}
		}
		if wait := interval - time.Since(*lastReq); wait > 0 {
			if err := sleepCtx(ctx, wait); err != nil {
				return nil, err
			}
		}
		*lastReq = time.Now()
		site, err := fetchOnce(ctx, client, date)
		if err == nil {
			return site, nil
		}
		lastErr = err
		var perm permanentError
		if errors.As(err, &perm) {
			break
		}
		log.Printf("重试 %s (%d/%d): %v", date, attempt+1, retries, err)
	}
	return nil, lastErr
}

type permanentError struct{ error }

func fetchOnce(ctx context.Context, client *http.Client, date string) (*algo.SiteDay, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, algo.SiteURL(date), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, permanentError{fmt.Errorf("HTTP 404")}
	case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500:
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	case resp.StatusCode != http.StatusOK:
		return nil, permanentError{fmt.Errorf("HTTP %d", resp.StatusCode)}
	}
	site, err := algo.ParseSite(io.LimitReader(resp.Body, 2<<20), date)
	if err != nil {
		return nil, permanentError{err}
	}
	return site, nil
}

// robotsAllows 极简 robots.txt 解析：只看 User-agent: * 组的 Allow/Disallow 前缀，最长匹配优先
func robotsAllows(ctx context.Context, client *http.Client, base, path string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/robots.txt", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return true, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("robots.txt HTTP %d", resp.StatusCode)
	}
	return parseRobots(resp.Body, path), nil
}

func parseRobots(r io.Reader, path string) bool {
	sc := bufio.NewScanner(r)
	inStar := false
	bestLen, allow := -1, true
	for sc.Scan() {
		line := strings.TrimSpace(strings.SplitN(sc.Text(), "#", 2)[0])
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k, v = strings.ToLower(strings.TrimSpace(k)), strings.TrimSpace(v)
		switch k {
		case "user-agent":
			inStar = v == "*"
		case "allow", "disallow":
			if !inStar || v == "" || !strings.HasPrefix(path, v) {
				continue
			}
			if len(v) > bestLen {
				bestLen, allow = len(v), k == "allow"
			}
		}
	}
	return allow
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func report(repo *almanac.Repo, start, end string, st, et time.Time, s stats) {
	m, err := repo.SourceMap(start, end)
	if err != nil {
		log.Printf("统计失败: %v", err)
		return
	}
	expect := int(et.Sub(st).Hours()/24) + 1
	merged, algoOnly := 0, 0
	for _, src := range m {
		if src == models.AlmanacSourceMerged {
			merged++
		} else {
			algoOnly++
		}
	}
	var mmDays int64
	db.DB.Model(&models.AlmanacDay{}).Where("date >= ? AND date <= ? AND mismatches != '[]'", start, end).Count(&mmDays)
	log.Printf("========== 汇总 ==========")
	log.Printf("本次处理 %d 天：抓取 %d，跳过(断点) %d，抓取失败算法兜底 %d，本次不一致 %d 天 / %d 项",
		s.total, s.crawled, s.skipped, s.fallback, s.mismatchDays, s.mismatchFields)
	log.Printf("库内覆盖 %d/%d 天（缺失 %d）：站点+算法 %d，纯算法 %d；累计交叉校验不一致 %d 天",
		len(m), expect, expect-len(m), merged, algoOnly, mmDays)
}
