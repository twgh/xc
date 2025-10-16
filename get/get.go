package get

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// NewCommand 创建 get 命令
func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "执行 go get -u github.com/twgh/xcgui",
		Long: `执行 go get -u github.com/twgh/xcgui 命令来获取或更新 xcgui 包，支持追加 go get 的其它参数。

示例:
  xc get          # 等于执行 go get -u github.com/twgh/xcgui`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否安装了 Go
			if _, err := exec.LookPath("go"); err != nil {
				fmt.Println("错误: 未找到 Go 命令，请确保已安装 Go 并将其添加到 PATH 环境变量中")
				os.Exit(1)
			}

			// 构建命令参数
			getArgs := []string{"get", "-u", "github.com/twgh/xcgui"}

			// 添加用户传递的其他参数
			getArgs = append(getArgs, args...)

			// 构建命令
			getCmd := exec.Command("go", getArgs...)

			// 设置输出
			getCmd.Stdout = os.Stdout
			getCmd.Stderr = os.Stderr

			// 执行命令
			fmt.Printf("执行: go %v\n", getArgs)
			if err := getCmd.Run(); err != nil {
				fmt.Printf("执行失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("xcgui 包获取/更新完成!")
		},
	}

	return cmd
}
