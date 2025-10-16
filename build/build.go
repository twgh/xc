package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// NewCommand 创建 build 命令
func NewCommand() *cobra.Command {
	// 定义标志变量
	var output string
	var verbose bool
	var work bool
	var race bool
	var trimpath bool
	var buildmode string
	var tags string

	var cmd = &cobra.Command{
		Use:   "build",
		Short: `执行 go build -ldflags="-s -w -H windowsgui" -trimpath`,
		Long: `执行 go build -ldflags="-s -w -H windowsgui" -trimpath 命令来构建项目，支持常用的 go build 参数。

示例:
  xc build                    # 等于执行 go build -ldflags="-s -w -H windowsgui" -trimpath
  xc build -o myapp.exe       # 指定输出文件名
  xc build -v                 # 打印编译包的名称
  xc build -work              # 打印临时工作目录的名称并不删除它
  xc build -race              # 启用数据竞争检测
  xc build -tags "netgo"      # 设置构建标签`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否安装了 Go
			if _, err := exec.LookPath("go"); err != nil {
				fmt.Println("错误: 未找到 Go 命令，请确保已安装 Go 并将其添加到 PATH 环境变量中")
				os.Exit(1)
			}

			// 构建命令参数
			buildArgs := []string{"build", "-ldflags=-s -w -H windowsgui"}

			// 添加输出文件参数
			if output != "" {
				buildArgs = append(buildArgs, "-o", output)
			}

			// 添加其他标志
			if verbose {
				buildArgs = append(buildArgs, "-v")
			}
			if work {
				buildArgs = append(buildArgs, "-work")
			}
			if race {
				buildArgs = append(buildArgs, "-race")
			}
			// 只有当用户没有明确指定 --no-trimpath 时才添加 -trimpath
			if trimpath {
				buildArgs = append(buildArgs, "-trimpath")
			}
			if buildmode != "" {
				buildArgs = append(buildArgs, "-buildmode", buildmode)
			}
			if tags != "" {
				buildArgs = append(buildArgs, "-tags", tags)
			}

			// 添加用户传递的其他参数
			buildArgs = append(buildArgs, args...)

			// 构建命令
			buildCmd := exec.Command("go", buildArgs...)

			// 设置输出
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr

			// 执行命令
			fmt.Printf("执行: go %s\n", strings.Join(buildArgs, " "))
			if err := buildCmd.Run(); err != nil {
				fmt.Printf("构建失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("项目构建完成!")
		},
	}

	// 添加常用的 go build 标志
	cmd.Flags().StringVarP(&output, "output", "o", "", "指定输出文件名")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "打印编译包的名称")
	cmd.Flags().BoolVar(&work, "work", false, "打印临时工作目录的名称并不删除它")
	cmd.Flags().BoolVar(&race, "race", false, "启用数据竞争检测")
	cmd.Flags().BoolVar(&trimpath, "trimpath", true, "移除结果可执行文件中的所有文件系统路径")
	cmd.Flags().StringVar(&buildmode, "buildmode", "", "构建模式")
	cmd.Flags().StringVar(&tags, "tags", "", "构建标签列表")

	return cmd
}
