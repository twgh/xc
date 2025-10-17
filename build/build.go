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
	var noTrimpath bool
	var buildmode string
	var tags string
	var printCommands bool
	var forceRebuild bool
	var dryRun bool
	var parallel int
	var msan bool
	var asan bool

	var cmd = &cobra.Command{
		Use:   "build",
		Short: `执行 go build -ldflags="-s -w -H windowsgui" -trimpath`,
		Long: `执行 go build -ldflags="-s -w -H windowsgui" -trimpath 命令来构建项目，支持添加常用的 go build 参数。

示例:
  xc build                    # 等于执行 go build -ldflags="-s -w -H windowsgui" -trimpath
  xc build -o myapp.exe       # 指定输出文件名
  xc build -v                 # 打印编译包的名称
  xc build -work              # 打印临时工作目录的名称并不删除它
  xc build -race              # 启用数据竞争检测
  xc build --no-trimpath      # 不添加 -trimpath 参数
  xc build -x                 # 打印执行的命令
  xc build -a                 # 强制重新构建已经是最新的包
  xc build -n                 # 打印命令但不执行
  xc build -p 4               # 设置并行执行的程序数量为 4
  xc build -msan              # 启用与内存清理器的互操作
  xc build -asan              # 启用与地址清理器的互操作
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
			// 默认添加 -trimpath，除非用户指定了 --no-trimpath
			if !noTrimpath {
				buildArgs = append(buildArgs, "-trimpath")
			}
			if buildmode != "" {
				buildArgs = append(buildArgs, "-buildmode", buildmode)
			}
			if tags != "" {
				buildArgs = append(buildArgs, "-tags", tags)
			}
			if printCommands {
				buildArgs = append(buildArgs, "-x")
			}
			if forceRebuild {
				buildArgs = append(buildArgs, "-a")
			}
			if dryRun {
				buildArgs = append(buildArgs, "-n")
			}
			if parallel > 0 {
				buildArgs = append(buildArgs, "-p", fmt.Sprintf("%d", parallel))
			}
			if msan {
				buildArgs = append(buildArgs, "-msan")
			}
			if asan {
				buildArgs = append(buildArgs, "-asan")
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
			if dryRun {
				fmt.Println("注意: 这是 dry-run 模式，不会实际执行命令")
			}
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
	cmd.Flags().BoolVar(&noTrimpath, "no-trimpath", false, "不添加 -trimpath 参数")
	cmd.Flags().StringVar(&buildmode, "buildmode", "", "构建模式")
	cmd.Flags().StringVar(&tags, "tags", "", "构建标签列表")
	cmd.Flags().BoolVarP(&printCommands, "print-commands", "x", false, "打印执行的命令")
	cmd.Flags().BoolVarP(&forceRebuild, "force-rebuild", "a", false, "强制重新构建已经是最新的包")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "打印命令但不执行")
	cmd.Flags().IntVarP(&parallel, "parallel", "p", 0, "设置并行执行的程序数量")
	cmd.Flags().BoolVar(&msan, "msan", false, "启用与内存清理器的互操作")
	cmd.Flags().BoolVar(&asan, "asan", false, "启用与地址清理器的互操作")

	return cmd
}
