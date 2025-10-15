package zipdownload

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twgh/xc/internal/downloader"
	"github.com/twgh/xc/internal/utils"
)

// 代理配置
type ProxyConfig struct {
	Name string
	URL  string
}

var proxies = []ProxyConfig{
	{"direct", ""}, // 直接下载
	{"ghfast", "https://ghfast.top/"},
	{"llkk", "https://gh.llkk.cc/"},
}

// NewCommand 创建 ZIP 下载命令
func NewCommand() *cobra.Command {
	var proxyName string
	var repoName string

	var cmd = &cobra.Command{
		Use:   "zipdownload",
		Short: "下载并解压 xcgui 和 example 仓库的源码 ZIP",
		Long: `下载并解压 xcgui 和 example 仓库的源码 ZIP, 支持使用代理加速下载, 支持指定 xcgui 或 example 仓库来下载。

示例:
  xc zipdownload                          # 使用直连下载
  xc zipdownload -p llkk             	  # 使用 llkk 代理下载
  xc zipdownload -p ghfast           	  # 使用 ghfast 代理下载
  xc zipdownload -n xcgui            	  # 下载 xcgui 仓库
  xc zipdownload -n example          	  # 下载 example 仓库`,
		Run: func(cmd *cobra.Command, args []string) {
			// 验证代理选择
			var selectedProxy ProxyConfig
			found := false
			for _, proxy := range proxies {
				if proxy.Name == proxyName {
					selectedProxy = proxy
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("错误: 无效的代理选项 '%s'\n", proxyName)
				fmt.Println("可用选项: direct, ghfast, llkk")
				os.Exit(1)
			}

			fmt.Printf("使用代理: %s\n", selectedProxy.Name)
			if selectedProxy.Name != "direct" {
				fmt.Printf("代理地址: %s\n", selectedProxy.URL)
			}

			// 定义默认要下载的仓库
			repos := []struct {
				url      string
				filename string
				finalDir string
			}{
				{
					url:      getRepoBranchUrl("twgh/xcgui"),
					filename: "xcgui.zip",
					finalDir: "xcgui",
				},
				{
					url:      getRepoBranchUrl("twgh/xcgui-example"),
					filename: "xcgui-example.zip",
					finalDir: "xcgui-example",
				},
			}

			// 如果指定了仓库名称，则只下载指定的仓库
			if repoName != "" {
				switch repoName {
				case "xcgui":
					repos = repos[:1]
				case "example":
					repos = repos[1:2]
				default:
					fmt.Printf("错误: 无效的仓库名 '%s'\n", repoName)
					fmt.Println("可用选项: xcgui, example")
					os.Exit(1)
				}
			}

			// 创建临时目录
			tempDir, err := os.MkdirTemp("", "github_downloads")
			if err != nil {
				fmt.Printf("创建临时目录失败: %v\n", err)
				return
			}
			defer os.RemoveAll(tempDir) // 程序结束时清理临时目录

			fmt.Println("临时目录:", tempDir)

			// 创建下载器
			dl := downloader.NewHTTPDownloader()

			// 处理每个仓库
			for _, repo := range repos {
				fmt.Printf("\n处理仓库: %s\n", repo.finalDir)
				fmt.Println("==============================")

				// 构建下载URL
				downloadURL := buildDownloadURL(repo.url, selectedProxy)
				fmt.Printf("下载地址: %s\n", downloadURL)

				// 下载文件
				zipPath := filepath.Join(tempDir, repo.filename)
				if err := downloadFile(dl, downloadURL, zipPath); err != nil {
					fmt.Printf("下载失败: %v\n", err)
					continue
				}

				// 解压文件
				extractedDir, err := unzip(zipPath, tempDir)
				if err != nil {
					fmt.Printf("解压失败: %v\n", err)
					continue
				}

				// 重命名文件夹（移除 -main 后缀）
				if err := renameDir(extractedDir, repo.finalDir); err != nil {
					fmt.Printf("重命名失败: %v\n", err)
					continue
				}

				fmt.Printf("成功处理: %s\n", repo.finalDir)
			}

			fmt.Println("\n所有操作完成!")
		},
	}

	cmd.Flags().StringVarP(&proxyName, "proxy", "p", "direct", "下载代理: direct, ghfast, llkk")
	cmd.Flags().StringVarP(&repoName, "name", "n", "", "指定要下载的仓库: xcgui, example")

	return cmd
}

// 获取指定仓库指定分支的源码 ZIP URL.
//
// repo: 用户名/仓库名, 如: twgh/xcgui.
//
// branch: 分支名, 不填则默认为 main.
func getRepoBranchUrl(repo string, branch ...string) string {
	branchName := "main"
	if len(branch) > 0 {
		branchName = branch[0]
		if branchName == "" {
			branchName = "main"
		}
	}
	return fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.zip", repo, branchName)
}

// 构建下载URL
func buildDownloadURL(originalURL string, proxy ProxyConfig) string {
	if proxy.Name == "direct" {
		return originalURL
	}
	return proxy.URL + originalURL
}

// downloadFile 下载文件
func downloadFile(dl downloader.Downloader, url, filepath string) error {
	fmt.Printf("下载中...\n")

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 下载数据
	body, err := dl.Download(url)
	if err != nil {
		return err
	}
	defer body.Close()

	// 获取内容长度
	if seeker, ok := body.(interface {
		Seek(int64, int) (int64, error)
	}); ok {
		size, _ := seeker.Seek(0, 2) // 移到末尾
		seeker.Seek(0, 0)            // 移到开头
		if size > 0 {
			fmt.Printf("文件大小: %.2f MB\n", float64(size)/1024/1024)
		}
	}

	// 复制内容到文件
	_, err = io.Copy(out, body)
	if err != nil {
		return err
	}

	fmt.Printf("下载完成: %s\n", filepath)
	return nil
}

// unzip 解压文件
func unzip(src, dest string) (string, error) {
	fmt.Printf("解压中: %s\n", src)

	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// 获取根目录名
	var rootDir string
	if len(r.File) > 0 {
		parts := strings.Split(r.File[0].Name, "/")
		if len(parts) > 0 {
			rootDir = parts[0]
		}
	}

	// 创建解压目录
	extractPath := filepath.Join(dest, "extracted")
	if err := utils.EnsureDirExists(extractPath); err != nil {
		return "", err
	}

	// 解压所有文件
	for _, f := range r.File {
		// 处理文件名
		fpath := filepath.Join(extractPath, f.Name)

		// 创建目录
		if f.FileInfo().IsDir() {
			utils.EnsureDirExists(fpath)
			continue
		}

		// 创建文件目录
		if err := utils.EnsureDirExists(filepath.Dir(fpath)); err != nil {
			return "", err
		}

		// 创建目标文件
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		// 打开源文件
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return "", err
		}

		// 复制内容
		_, err = io.Copy(outFile, rc)

		// 关闭文件
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}

	fmt.Printf("解压完成: %s\n", extractPath)
	return filepath.Join(extractPath, rootDir), nil
}

// renameDir 重命名目录
func renameDir(oldPath, newName string) error {
	// 检查源目录是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("源目录不存在: %s", oldPath)
	}

	// 获取新路径
	newPath := filepath.Join(filepath.Dir(oldPath), newName)

	// 重命名目录
	fmt.Printf("重命名: %s -> %s\n", filepath.Base(oldPath), newName)
	if err := os.Rename(oldPath, newPath); err != nil {
		// 如果重命名失败，可能是跨磁盘，使用MoveDir工具函数
		fmt.Printf("直接重命名失败，尝试跨磁盘移动: %v\n", err)
	}

	// 移动目录到当前工作目录
	currentDir, err := utils.GetWorkingDir()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	finalPath := filepath.Join(currentDir, newName)
	fmt.Printf("移动目录到: %s\n", finalPath)

	// 使用工具函数移动目录（支持跨磁盘）
	return utils.MoveDir(newPath, finalPath)
}
