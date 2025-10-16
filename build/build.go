package build

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// NewCommand 创建 build 命令
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "build",
		Short: `执行 go build -ldflags="-s -w -H windowsgui" -trimpath`,
		Long: `执行 go build -ldflags="-s -w -H windowsgui" -trimpath 命令来构建项目，支持追加 go build 的其它参数。

示例:
  xc build                    # 等于执行 go build -ldflags="-s -w -H windowsgui" -trimpath
  xc build -o myapp.exe       # 指定输出文件名
  xc build -o dist/app.exe    # 指定输出路径和文件名`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否安装了 Go
			if _, err := exec.LookPath("go"); err != nil {
				fmt.Println("错误: 未找到 Go 命令，请确保已安装 Go 并将其添加到 PATH 环境变量中")
				os.Exit(1)
			}

			// 构建命令参数
			buildArgs := []string{"build", "-ldflags=-s -w -H windowsgui", "-trimpath"}

			// 添加用户传递的其他参数
			buildArgs = append(buildArgs, args...)

			// 构建命令
			buildCmd := exec.Command("go", buildArgs...)

			// 设置输出
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr

			// 执行命令
			fmt.Printf("执行: go %v\n", buildArgs)
			if err := buildCmd.Run(); err != nil {
				fmt.Printf("构建失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("项目构建完成!")
		},
	}

	return cmd
}
