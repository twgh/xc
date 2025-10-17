package get

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// NewCommand 创建 get 命令
func NewCommand() *cobra.Command {
	// 定义标志变量
	var patch bool
	var tool bool
	var printCommands bool

	var cmd = &cobra.Command{
		Use:   "get",
		Short: "执行 go get -u github.com/twgh/xcgui",
		Long: `执行 go get -u github.com/twgh/xcgui 命令来获取或更新 xcgui 包，支持添加常用的 go get 参数。

示例:
  xc get          # 等于执行 go get -u github.com/twgh/xcgui
  xc get -t       # 同时考虑测试依赖
  xc get --patch  # 使用补丁版本更新
  xc get -tool    # 为每个列出的包添加匹配的工具行到 go.mod
  xc get -x       # 打印执行的命令`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否安装了 Go
			if _, err := exec.LookPath("go"); err != nil {
				fmt.Println("错误: 未找到 Go 命令，请确保已安装 Go 并将其添加到 PATH 环境变量中")
				os.Exit(1)
			}

			// 构建命令参数
			getArgs := []string{"get", "-u"}

			// 添加补丁标志
			if patch {
				getArgs[len(getArgs)-1] = "-u=patch" // 替换 -u 为 -u=patch
			}

			// 添加其他标志
			if tool {
				getArgs = append(getArgs, "-tool")
			}
			if printCommands {
				getArgs = append(getArgs, "-x")
			}

			// 添加包名
			getArgs = append(getArgs, "github.com/twgh/xcgui")

			// 添加用户传递的其他参数
			getArgs = append(getArgs, args...)

			// 构建命令
			getCmd := exec.Command("go", getArgs...)

			// 设置输出
			getCmd.Stdout = os.Stdout
			getCmd.Stderr = os.Stderr

			// 执行命令
			fmt.Printf("执行: go %s\n", strings.Join(getArgs, " "))
			if err := getCmd.Run(); err != nil {
				fmt.Printf("执行失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("xcgui 包获取/更新完成!")
		},
	}

	// 添加常用的 go get 标志
	cmd.Flags().BoolVar(&patch, "patch", false, "更新依赖到最新的补丁版本")
	cmd.Flags().BoolVar(&tool, "tool", false, "为每个列出的包添加匹配的工具行到 go.mod")
	cmd.Flags().BoolVarP(&printCommands, "print-commands", "x", false, "打印执行的命令")

	return cmd
}
