package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestProxyDownload 测试几种代理方式能否跑通
//
// 依次使用每个代理下载 https://github.com/twgh/xcgui/archive/main.zip，
// 验证下载成功后立即删除文件，避免残留。
func TestProxyDownload(t *testing.T) {
	const targetURL = "https://github.com/twgh/xcgui/archive/main.zip"

	// 单文件超时 60 秒
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	for _, proxy := range Proxies {
		proxyName := proxy.Name
		if proxyName == "direct" {
			proxyName = "直连"
		}
		t.Run(proxyName, func(t *testing.T) {
			downloadURL := BuildFileDownloadURL(targetURL, proxy)
			fmt.Printf("使用 %s 下载：%s\n", proxyName, downloadURL)

			req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
			if err != nil {
				t.Fatalf("%s 构建请求失败：%v", proxyName, err)
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-Downloader/1.0)")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s 下载请求失败：%v", proxyName, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				t.Fatalf("%s 下载返回非预期状态码：%d", proxyName, resp.StatusCode)
			}

			// 写出到临时文件
			fileName := fmt.Sprintf("xcgui_main_%s.zip", proxy.Name)
			filePath := filepath.Join(t.TempDir(), fileName)
			out, err := os.Create(filePath)
			if err != nil {
				t.Fatalf("%s 创建文件失败：%v", proxyName, err)
			}

			n, err := io.Copy(out, resp.Body)
			out.Close()
			if err != nil {
				t.Fatalf("%s 写入文件失败：%v", proxyName, err)
			}
			if n == 0 {
				t.Fatalf("%s 下载内容为空", proxyName)
			}

			fmt.Printf("%s 下载成功，大小：%d 字节，路径：%s\n", proxyName, n, filePath)

			// 完成后删除下载的文件
			if err := os.Remove(filePath); err != nil {
				t.Fatalf("%s 删除文件失败：%v", proxyName, err)
			}
			fmt.Printf("%s 下载文件已删除\n", proxyName)
		})
	}
}
