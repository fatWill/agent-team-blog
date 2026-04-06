package ipgeo

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	searcher *xdb.Searcher
	once     sync.Once
	initErr  error
)

// Init 初始化 IP 地理位置解析器（加载整个 xdb 到内存，线程安全）
func Init(dbPath string) error {
	once.Do(func() {
		cBuff, err := xdb.LoadContentFromFile(dbPath)
		if err != nil {
			initErr = fmt.Errorf("加载 ip2region.xdb 失败: %w", err)
			return
		}

		searcher, err = xdb.NewWithBuffer(xdb.IPv4, cBuff)
		if err != nil {
			initErr = fmt.Errorf("创建 ip2region searcher 失败: %w", err)
			return
		}

		log.Printf("✅ ip2region 初始化成功 (%s)", dbPath)
	})
	return initErr
}

// GeoInfo IP 地理位置信息
type GeoInfo struct {
	Country  string `json:"country"`
	Region   string `json:"region"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

// Search 解析 IP 地址的地理位置
// ip2region v3 xdb 返回格式: "国家|省份|城市|ISP|国家代码"
// 例如: "中国|广东省|深圳市|联通|CN"
func Search(ip string) *GeoInfo {
	if searcher == nil {
		return &GeoInfo{}
	}

	// 跳过内网 IP
	if isPrivateIP(ip) {
		return &GeoInfo{Country: "内网", Province: "内网IP", City: "内网IP"}
	}

	region, err := searcher.Search(ip)
	if err != nil {
		return &GeoInfo{}
	}

	parts := strings.Split(region, "|")
	if len(parts) < 4 {
		return &GeoInfo{}
	}

	info := &GeoInfo{
		Country:  cleanField(parts[0]),
		Province: cleanField(parts[1]),
		City:     cleanField(parts[2]),
		ISP:      cleanField(parts[3]),
	}
	if len(parts) >= 5 {
		info.Region = cleanField(parts[4]) // 国家代码
	}

	return info
}

// cleanField 清理 "0" 值为空字符串
func cleanField(s string) string {
	s = strings.TrimSpace(s)
	if s == "0" {
		return ""
	}
	return s
}

// isPrivateIP 判断是否为内网 IP
func isPrivateIP(ip string) bool {
	return strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.16.") ||
		strings.HasPrefix(ip, "172.17.") ||
		strings.HasPrefix(ip, "172.18.") ||
		strings.HasPrefix(ip, "172.19.") ||
		strings.HasPrefix(ip, "172.2") ||
		strings.HasPrefix(ip, "172.30.") ||
		strings.HasPrefix(ip, "172.31.") ||
		strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "127.") ||
		ip == "::1" ||
		ip == ""
}
