//go:generate goversioninfo
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/twgh/xc/build"
	"github.com/twgh/xc/dlldownload"
	"github.com/twgh/xc/get"
	"github.com/twgh/xc/zipdownload"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "xc",
		Short: "xc 是一个 xcgui 助手类型的命令行工具",
		Long: `xc 是一个 xcgui 助手类型的命令行工具, 功能包括给项目添加 xcgui、编译程序、下载 xcgui 和 example 仓库的源码 ZIP、下载 xcgui.dll 文件等。

使用方法:
  xc [command]

可用命令:
  get           执行 go get -u github.com/twgh/xcgui
  build         执行 go build -ldflags="-s -w -H windowsgui" -trimpath
  zipdownload   下载并解压 xcgui 和 example 仓库的源码 ZIP
  dlldownload   下载 xcgui.dll 文件
  version       显示版本信息
  help          显示命令帮助信息

使用 "xc [command] --help" 获取更多关于某个命令的信息。`,
	}

	// 添加版本命令
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("xc version 0.0.1")
		},
	}

	// 添加 DLL 下载命令
	rootCmd.AddCommand(dlldownload.NewCommand())

	// 添加 ZIP 下载命令
	rootCmd.AddCommand(zipdownload.NewCommand())

	// 添加 get 命令
	rootCmd.AddCommand(get.NewCommand())

	// 添加 build 命令
	rootCmd.AddCommand(build.NewCommand())

	// 添加版本命令
	rootCmd.AddCommand(versionCmd)

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
