package utils

import (
	"fmt"
	"io"
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
	if testURL == "" {
		testURL = "https://github.com"
	}

	// 优先尝试 HEAD
	req, err := http.NewRequest(http.MethodHead, testURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-SiteChecker/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()

	// HEAD 成功则直接返回
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}

	// HEAD 返回 403/405 等，降级为 GET 重试
	reqGet, err := http.NewRequest(http.MethodGet, testURL, nil)
	if err != nil {
		return false
	}
	reqGet.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-SiteChecker/1.0)")

	respGet, err := client.Do(reqGet)
	if err != nil {
		return false
	}
	defer respGet.Body.Close()

	// 读取并丢弃 Body，否则底层 TCP 连接不会被复用
	io.Copy(io.Discard, respGet.Body)

	return respGet.StatusCode >= 200 && respGet.StatusCode < 400
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
