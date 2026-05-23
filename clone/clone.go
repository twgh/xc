package clone

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twgh/xc/internal/utils"
)

// checkGitExists 检测 git 是否安装
func checkGitExists() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// NewCommand 创建 clone 命令
func NewCommand() *cobra.Command {
	var cloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "克隆 GitHub 仓库",
		Long: `克隆 GitHub 仓库，自动测试并选择可用的代理加速下载，如果都无法访问则中止克隆。

示例:
  xc clone https://github.com/twgh/xcgui          # 克隆 xcgui 仓库
  xc clone https://github.com/xxx/xxxx            # 克隆任意仓库`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否提供了仓库 URL
			if len(args) == 0 {
				fmt.Println("错误: 请提供要克隆的仓库 URL")
				fmt.Println("使用方法: xc clone <repository_url>")
				os.Exit(1)
			}

			repoURL := args[0]

			// 检查 git 是否安装
			if !checkGitExists() {
				fmt.Println("错误: 未检测到 git，请先安装 git")
				fmt.Println("下载地址: https://git-scm.com/download/win")
				os.Exit(1)
			}

			// 自动测试并选择可用代理
			selectedProxy, err := utils.SelectAvailableProxy()
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("使用代理: %s\n", selectedProxy.Name)
			if selectedProxy.Name != "direct" {
				fmt.Printf("代理地址: %s\n", selectedProxy.URL)
			}

			// 构建最终的 clone URL
			finalURL := utils.BuildDownloadURL(repoURL, selectedProxy)
			fmt.Printf("克隆地址: %s\n", finalURL)

			// 获取仓库名称作为目标目录
			repoName := filepath.Base(repoURL)
			repoName = strings.TrimSuffix(repoName, ".git")

			fmt.Printf("目标目录: %s\n", repoName)
			fmt.Println("==============================")

			// 执行 git clone
			fmt.Println("正在执行 git clone...")
			gitCmd := exec.Command("git", "clone", finalURL, repoName)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr

			if err := gitCmd.Run(); err != nil {
				fmt.Printf("克隆失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("\n克隆完成!")
		},
	}

	return cloneCmd
}
