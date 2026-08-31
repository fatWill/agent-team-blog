package wechat

import (
	"strings"
)

// siteURL 站点根 URL，由 main 在启动时通过 SetSiteURL 注入（默认值保持线上现状）
var siteURL = "https://fatwill.cloud"

// SetSiteURL 设置站点根 URL（用于拼接文章原文链接），空值则保留默认值
func SetSiteURL(u string) {
	if u = strings.TrimRight(strings.TrimSpace(u), "/"); u != "" {
		siteURL = u
	}
}

// SiteURL 返回当前生效的站点根 URL（不带结尾斜杠）
func SiteURL() string {
	return siteURL
}

// articleURL 拼接文章在博客站点上的绝对访问地址
func articleURL(articleID string) string {
	return siteURL + "/articles/" + articleID
}
