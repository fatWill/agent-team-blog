package nginx

import (
	"log"
	"os/exec"
)

// purgeScript 缓存清除脚本路径
const purgeScript = "/usr/local/bin/purge-nginx-cache.sh"

// PurgeCache 异步清除 Nginx 缓存
// 通过调用服务器上的 purge-nginx-cache.sh 脚本清除指定路径的缓存。
// 调用方式为异步非阻塞（内部启动 goroutine），不影响主流程。
// 清除失败只打 warn 日志，不返回错误。
//
// paths 参数为需要清除缓存的 URL 路径列表，例如 ["/", "/articles/abc-123"]
//   - "/" 会转为 --home 参数（清除首页缓存）
//   - 其他路径直接作为脚本参数传递
func PurgeCache(paths []string) {
	if len(paths) == 0 {
		return
	}

	go func() {
		for _, path := range paths {
			var arg string
			if path == "/" {
				arg = "--home"
			} else {
				arg = path
			}

			out, err := exec.Command(purgeScript, arg).CombinedOutput()
			if err != nil {
				log.Printf("⚠️ Nginx 缓存清除失败 [%s %s]: %v, output: %s", purgeScript, arg, err, string(out))
			}
		}
	}()
}