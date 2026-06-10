package getskill

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/twgh/xc/internal/utils"
)

// NewCommand 创建 getskill 命令
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "getskill",
		Short: "下载 go-xcgui-dev 技能",
		Long: `从 GitHub 下载 go-xcgui-dev 仓库的源码 ZIP 并解压到当前目录, 自动使用代理加速下载。

示例:
  xc getskill`,
		Run: func(cmd *cobra.Command, args []string) {
			url := "https://github.com/twgh/go-xcgui-dev/archive/main.zip"
			filename := "go-xcgui-dev.zip"
			finalDir := "go-xcgui-dev"

			// 创建临时目录
			tempDir, err := os.MkdirTemp("", "getskill_download")
			if err != nil {
				fmt.Printf("创建临时目录失败: %v\n", err)
				os.Exit(1)
			}
			defer os.RemoveAll(tempDir)

			fmt.Println("临时目录:", tempDir)

			// 准备 zipPath
			zipPath := filepath.Join(tempDir, filename)

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

				downloadURL := utils.BuildFileDownloadURL(url, proxy)
				fmt.Printf("下载地址: %s\n", downloadURL)

				if err := utils.DownloadFile(downloadURL, zipPath); err != nil {
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

			// 解压文件
			extractedDir, err := utils.Unzip(zipPath, tempDir)
			if err != nil {
				fmt.Printf("解压失败: %v\n", err)
				os.Exit(1)
			}

			// 判断目标目录是否已存在
			cwd, err := utils.GetWorkingDir()
			if err != nil {
				fmt.Printf("获取当前目录失败: %v\n", err)
				os.Exit(1)
			}
			targetDir := filepath.Join(cwd, finalDir)

			if _, err := os.Stat(targetDir); err == nil {
				// 目录已存在，合并覆盖，保留用户自己的文件
				fmt.Printf("目标目录已存在，合并覆盖: %s\n", targetDir)
				if err := utils.MergeDir(extractedDir, targetDir); err != nil {
					fmt.Printf("合并目录失败: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("合并完成!")
			} else {
				// 目录不存在，重命名并移动
				if err := utils.RenameDir(extractedDir, finalDir); err != nil {
					fmt.Printf("重命名失败: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Println("\n下载并安装完成!")
		},
	}

	return cmd
}
