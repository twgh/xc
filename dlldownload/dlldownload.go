package dlldownload

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// NewCommand 创建 DLL 下载命令
func NewCommand() *cobra.Command {
	var (
		output  string
		version string
		bit     uint
	)

	var cmd = &cobra.Command{
		Use:   "dlldownload",
		Short: "下载 xcgui.dll 文件",
		Long: `下载 xcgui.dll 文件，默认会下载最新版本的 64 位 DLL 文件。可指定版本号、位数、输出文件名及位置。

示例:
  xc dlldownload                          # 下载最新版本的 64 位 DLL
  xc dlldownload -v 3.3.5.0               # 下载指定版本的 DLL
  xc dlldownload -b 32 -o my.dll          # 下载 32 位 DLL 并保存为 my.dll
  xc dlldownload -v 3.3.5.0 -b 32 -o test.dll  # 下载指定版本的 32 位 DLL 并保存为 test.dll`,
		Run: func(cmd *cobra.Command, args []string) {
			// 获取最新版本号
			version = strings.TrimSpace(version)
			if version == "" {
				latest, err := getLatestVersion()
				if err != nil {
					fmt.Println("获取最新版本失败:", err)
					return
				}
				version = latest
			}

			// 删首尾空
			output = strings.TrimSpace(output)
			if output == "" {
				output = "xcgui.dll"
			}

			// 判断位数, 得到下载地址
			if bit == 32 || bit == 86 {
				bit = 32
			} else {
				bit = 64
			}

			addr := ""
			if bit == 32 {
				addr = fmt.Sprintf("https://cnb.cool/twgh521/xcguidll/-/releases/download/%s/xcgui-32.dll", version)
			} else {
				addr = fmt.Sprintf("https://cnb.cool/twgh521/xcguidll/-/releases/download/%s/xcgui.dll", version)
			}

			// 开始下载 dll
			fmt.Printf("开始下载 xcgui.dll\n版本: %s\n位数: %d-bit\n输出文件: %s\n", version, bit, output)
			quit := make(chan bool)
			go func() {
				for i := 0; i < 1500; i++ { // 超过 300 秒就判定为下载失败
					select {
					case <-quit:
						return
					default:
						fmt.Print(".")
						time.Sleep(time.Millisecond * 200)
					}
				}
				fmt.Println("\n下载失败: 超时")
				os.Exit(1)
			}()

			data, err := getDll(addr)
			if err != nil {
				quit <- true
				fmt.Println("\n下载失败:", err.Error())
				return
			}

			if len(data) < 1.5*1024*1024 { // 小于 1.5M 肯定下载失败了
				quit <- true
				fmt.Println("\n下载失败: 文件大小异常")
				return
			}

			err = os.WriteFile(output, data, 0777)
			if err != nil {
				quit <- true
				fmt.Println("\n保存文件失败:", err)
				return
			}

			// 计算 data 的 md5
			md5Str := fmt.Sprintf("%x", md5.Sum(data))
			fmt.Printf("\nMD5: %s\n", strings.ToUpper(md5Str))

			quit <- true
			fmt.Println("下载成功!")
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "xcgui.dll", "输出文件名")
	cmd.Flags().StringVarP(&version, "version", "v", "", "xcgui.dll 的版本号, 例如: 3.3.5.0")
	cmd.Flags().UintVarP(&bit, "bit", "b", 64, "xcgui.dll 的位数 (32 或 64)")

	return cmd
}

// 从指定网址下载 dll
func getDll(addr string) ([]byte, error) {
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

// 获取 dll 的最新版本号
func getLatestVersion() (string, error) {
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
