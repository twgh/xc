package xc

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

// GetDll 从指定网址下载 dll
func GetDll(addr string) ([]byte, error) {
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("file not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	if bytes.Contains(body, []byte("NoSuchKey")) {
		return nil, errors.New("file not found")
	}
	return body, err
}

// GetLatestVersion 获取 dll 的最新版本号
func GetLatestVersion() (string, error) {
	res, err := http.Get("https://cnb.cool/twgh521/xcguidll/-/git/raw/main/version.txt?download=true")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(body))
	if version == "" {
		return "", errors.New("failed to get the latest version number")
	}
	return version, nil
}
