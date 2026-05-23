package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ProxyConfig 代理配置
type ProxyConfig struct {
	Name string
	URL  string
}

// Proxies 代理列表（按优先级排序）
var Proxies = []ProxyConfig{
	{"ghfast", "https://ghfast.top/"},
	{"llkk", "https://gh.llkk.cc/"},
	{"direct", ""}, // 直接下载
}

// IsSiteReachable 测试网站是否可达
func IsSiteReachable(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	testURL := url
	if url == "" {
		testURL = "https://github.com"
	}

	resp, err := client.Head(testURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// SelectAvailableProxy 自动测试并选择可用的代理
func SelectAvailableProxy() (ProxyConfig, error) {
	fmt.Println("正在测试代理可用性...")
	fmt.Println("==============================")

	for _, proxy := range Proxies {
		proxyName := proxy.Name
		if proxyName == "direct" {
			proxyName = "github.com (直连)"
		}

		fmt.Printf("测试 %s ... ", proxyName)
		if IsSiteReachable(proxy.URL) {
			fmt.Println("✓ 可达")
			fmt.Println("==============================")
			return proxy, nil
		}
		fmt.Println("✗ 不可达")
	}

	fmt.Println("==============================")
	return ProxyConfig{}, fmt.Errorf("所有代理都无法访问，请检查网络连接")
}

// BuildDownloadURL 构建下载 URL（用于 git clone）
func BuildDownloadURL(repoURL string, proxy ProxyConfig) string {
	// 移除 .git 后缀（如果有）
	repoURL = strings.TrimSuffix(repoURL, ".git")

	if proxy.Name == "direct" {
		return repoURL
	}
	return proxy.URL + repoURL
}

// BuildFileDownloadURL 构建文件下载 URL
func BuildFileDownloadURL(originalURL string, proxy ProxyConfig) string {
	if proxy.Name == "direct" {
		return originalURL
	}
	return proxy.URL + originalURL
}