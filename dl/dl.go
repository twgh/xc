package dl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twgh/xc/internal/utils"
)

// getFilenameFromURL 从 URL 获取文件名
func getFilenameFromURL(url string) string {
	// 尝试从 URL 中提取文件名
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		// 如果文件名包含查询参数，进一步处理
		if strings.Contains(filename, "?") {
			filename = strings.Split(filename, "?")[0]
		}
		if filename != "" && strings.Contains(filename, ".") {
			return filename
		}
	}
	return "download_file"
}

// NewCommand 创建 dl 命令
func NewCommand() *cobra.Command {
	var outputName string

	var cmd = &cobra.Command{
		Use:   "dl",
		Short: "下载文件",
		Long: `下载文件，自动测试并选择可用的代理加速下载，如果都无法访问则中止下载。

示例:
  xc dl https://github.com/twgh/xcgui/archive/main.zip               		# 下载分支源码
  xc dl https://github.com/twgh/xcgui/archive/refs/tags/v1.4.0.zip    		# 下载 Release 源码
  xc dl https://github.com/twgh/xcgui/releases/download/v1.3.393/xcgui.dll   	# 下载 Release 附件
  xc dl https://raw.githubusercontent.com/twgh/xc/main/README.md -o readme.md   # 下载 Raw 文件并指定输出文件名`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否提供了 URL
			if len(args) == 0 {
				fmt.Println("错误: 请提供要下载的 URL")
				fmt.Println("使用方法: xc dl <url> [-o output_filename]")
				os.Exit(1)
			}

			url := args[0]

			// 确定输出文件名
			filename := outputName
			if filename == "" {
				filename = getFilenameFromURL(url)
			}

			// 获取当前工作目录
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("获取当前目录失败: %v\n", err)
				os.Exit(1)
			}
			filepath := filepath.Join(cwd, filename)

			fmt.Printf("保存路径: %s\n", filepath)
			fmt.Println("==============================")

			// 如果 URL 中不包含 github，则直接使用直连下载
			if !strings.Contains(url, "github") {
				fmt.Printf("下载地址: %s\n", url)
				if err := utils.DownloadFile(url, filepath); err != nil {
					fmt.Printf("下载失败: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("\n下载完成!")
				return
			}

			// 遍历所有代理尝试下载，直到成功
			success := false
			for i, proxy := range utils.Proxies {
				proxyName := proxy.Name
				if proxyName == "direct" {
					proxyName = "直连"
				}
				fmt.Printf("使用代理: %s\n", proxyName)
				if proxy.Name != "direct" {
					fmt.Printf("代理地址: %s\n", proxy.URL)
				}

				// 构建最终的下载 URL（用户输入的是原始 URL，自动拼接代理前缀）
				finalURL := utils.BuildFileDownloadURL(url, proxy)
				fmt.Printf("下载地址: %s\n", finalURL)

				// 下载文件
				if err := utils.DownloadFile(finalURL, filepath); err != nil {
					fmt.Printf("下载失败: %v\n", err)
					if i < len(utils.Proxies)-1 {
						fmt.Println("尝试下一个代理...")
					}
					continue
				}
				success = true
				break
			}

			if !success {
				fmt.Println("所有下载方式均失败")
				os.Exit(1)
			}

			fmt.Println("\n下载完成!")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "指定输出文件名")

	return cmd
}
