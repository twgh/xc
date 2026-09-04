package res

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed assets/go-winres.exe
var winresExe []byte

//go:embed assets/winres.json
var winresJSON []byte

//go:embed assets/icon.ico
var iconIco []byte

// winresExeName 释放到临时目录时的文件名
const winresExeName = "go-winres.exe"

// ensureWinresExe 将内置的 go-winres.exe 释放到系统临时目录,
// 仅当临时目录中不存在该文件时才进行释放, 返回释放后的完整路径。
func ensureWinresExe() (string, error) {
	tmpDir := os.TempDir()
	exePath := filepath.Join(tmpDir, winresExeName)

	// 如果临时目录中已存在该文件, 则直接复用, 不再重复释放
	if _, err := os.Stat(exePath); err == nil {
		return exePath, nil
	}

	if err := os.WriteFile(exePath, winresExe, 0o755); err != nil {
		return "", fmt.Errorf("释放 %s 失败: %w", winresExeName, err)
	}
	return exePath, nil
}

// NewCommand 创建 res 命令, 转发 go-winres 的子命令
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "res",
		Short: "执行 go-winres 命令（给程序添加 Windows 资源 / 版本信息 / 程序清单）",
		Long: `执行 go-winres 命令（给程序添加 Windows 资源 / 版本信息 / 程序清单）, 相当于在命令行中直接调用 go-winres。

go-winres 常用子命令:
  init     在当前目录创建初始的 ./winres/winres.json
  make     根据配置生成供 "go build" 使用的 syso 文件
  simply   以简化方式生成 syso 文件
  extract  从可执行文件中提取全部资源
  patch    替换可执行文件（exe / dll）中的资源

示例:
  xc res help            # 查看 go-winres 帮助
  xc res simply            # 生成简化的 syso 文件
  xc res make              # 根据配置生成 syso 文件
  xc res init              # 初始化 winres.json 配置
  xc res extract app.exe   # 从 app.exe 提取资源`,
		// 关闭 flag 解析, 让所有参数原样转发给 go-winres
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			exePath, err := ensureWinresExe()
			if err != nil {
				return err
			}

			// 将剩余参数原样转发给 go-winres 执行
			c := exec.Command(exePath, args...)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			if err := c.Run(); err != nil {
				// 若 go-winres 自身返回非零退出码, 透传其退出码
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					os.Exit(exitErr.ExitCode())
				}
				return err
			}

			// 若执行的是 init 子命令, 则覆盖 winres/winres.json 并删除图标文件
			if len(args) > 0 && args[0] == "init" {
				if err := applyInitOverlay(); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

// applyInitOverlay 在 go-winres init 之后执行:
// 将内置的 winres.json 覆盖到 winres/ 目录, 并删除 go-winres 生成的 icon.png 和 icon16.png。
func applyInitOverlay() error {
	winresDir := "winres"

	if _, err := os.Stat(winresDir); err != nil {
		return fmt.Errorf("未找到 %s 目录, 可能 go-winres init 未能成功执行: %w", winresDir, err)
	}

	// 覆盖 winres/winres.json
	targetJSON := filepath.Join(winresDir, "winres.json")
	if err := os.WriteFile(targetJSON, winresJSON, 0o644); err != nil {
		return fmt.Errorf("覆盖 %s 失败: %w", targetJSON, err)
	}
	// fmt.Printf("已覆盖 %s\n", targetJSON)

	// 删除 go-winres init 生成的图标文件
	for _, name := range []string{"icon.png", "icon16.png"} {
		p := filepath.Join(winresDir, name)
		if _, err := os.Stat(p); err == nil {
			if err := os.Remove(p); err != nil {
				return fmt.Errorf("删除 %s 失败: %w", p, err)
			}
			// fmt.Printf("已删除 %s\n", p)
		} else {
			fmt.Printf("未找到 %s, 跳过删除\n", p)
		}
	}

	// 释放内置的 icon.ico（winres.json 中 RT_GROUP_ICON 引用了它）
	targetIco := filepath.Join(winresDir, "icon.ico")
	if err := os.WriteFile(targetIco, iconIco, 0o644); err != nil {
		return fmt.Errorf("写入 %s 失败: %w", targetIco, err)
	}
	fmt.Println("提示: ico 中最好包含 16，24，32，48，64，256 尺寸的图片以达到最佳显示效果")
	fmt.Println("提示: 继续执行 xc res make 可生成 .syso 文件")
	return nil
}
