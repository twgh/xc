package downloader

import (
	"io"
	"net/http"
)

// Downloader 定义下载器接口
type Downloader interface {
	Download(url string) (io.ReadCloser, error)
}

// HTTPDownloader HTTP下载器实现
type HTTPDownloader struct {
	client *http.Client
}

// NewHTTPDownloader 创建新的HTTP下载器
func NewHTTPDownloader() *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
}

// Download 从指定URL下载数据
func (d *HTTPDownloader) Download(url string) (io.ReadCloser, error) {
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 发起请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	return resp.Body, nil
}

// HTTPError HTTP错误类型
type HTTPError struct {
	StatusCode int
	Status     string
}

func (e *HTTPError) Error() string {
	return e.Status
}
