package nginx

import (
	"log"
	"net/http"
	"time"
)

// purgeClient 复用 HTTP 客户端，避免每次请求都创建连接
var purgeClient = &http.Client{
	Timeout: 3 * time.Second,
	// 不跟随重定向
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// PurgeCache 异步清除 Nginx 缓存
// 通过向 Nginx 的 ngx_cache_purge 模块发送 PURGE 请求来清除指定路径的缓存。
// 调用方式为异步非阻塞（内部启动 goroutine），不影响主流程。
// 清除失败只打 warn 日志，不返回错误。
//
// paths 参数为需要清除缓存的 URL 路径列表，例如 ["/", "/articles/abc-123"]
// 实际发送的请求为 PURGE http://127.0.0.1/purge{path}
func PurgeCache(paths []string) {
	if len(paths) == 0 {
		return
	}

	go func() {
		for _, path := range paths {
			purgeURL := "http://127.0.0.1/purge" + path
			req, err := http.NewRequest("PURGE", purgeURL, nil)
			if err != nil {
				log.Printf("⚠️ Nginx 缓存清除请求构建失败 [%s]: %v", purgeURL, err)
				continue
			}

			resp, err := purgeClient.Do(req)
			if err != nil {
				log.Printf("⚠️ Nginx 缓存清除请求失败 [%s]: %v", purgeURL, err)
				continue
			}
			resp.Body.Close()

			// ngx_cache_purge 成功返回 200，缓存不存在返回 404（均视为正常）
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
				log.Printf("⚠️ Nginx 缓存清除异常 [%s]: status=%d", purgeURL, resp.StatusCode)
			}
		}
	}()
}
