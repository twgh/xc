package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/twgh/xc/internal/downloader"
)

// EnsureDirExists 确保目录存在，如果不存在则创建
func EnsureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// MoveDir 跨磁盘移动目录到目标位置
func MoveDir(src, dst string) error {
	// 如果目标目录已存在，先删除
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("删除已存在的目录失败: %v", err)
		}
	}

	// 尝试直接重命名（同一磁盘内）
	err := os.Rename(src, dst)
	if err == nil {
		// 成功直接重命名
		return nil
	}

	// 如果重命名失败，可能是跨磁盘，使用复制和删除的方式
	return moveDirCrossDisk(src, dst)
}

// moveDirCrossDisk 跨磁盘移动目录
func moveDirCrossDisk(src, dst string) error {
	// 复制目录
	if err := copyDir(src, dst); err != nil {
		return fmt.Errorf("复制目录失败: %v", err)
	}

	// 删除源目录
	if err := os.RemoveAll(src); err != nil {
		return fmt.Errorf("删除源目录失败: %v", err)
	}

	return nil
}

// copyDir 复制目录
func copyDir(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 复制每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 同步文件
	return dstFile.Sync()
}

// GetWorkingDir 获取当前工作目录
func GetWorkingDir() (string, error) {
	return os.Getwd()
}

// DownloadFile 下载文件
func DownloadFile(url, filepath string) error {
	fmt.Println("正在下载...")

	// 创建下载器
	dl := downloader.NewHTTPDownloader()

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
		size, _ := seeker.Seek(0, 2)
		seeker.Seek(0, 0)
		if size > 0 {
			fmt.Printf("文件大小: %.2f MB\n", float64(size)/1024/1024)
		}
	}

	// 复制内容到文件
	_, err = io.Copy(out, body)
	if err != nil {
		return err
	}

	fmt.Printf("文件已保存: %s\n", filepath)
	return nil
}

// Unzip 解压 ZIP 文件到目标目录，返回解压后的根目录路径
func Unzip(src, dest string) (string, error) {
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
	if err := EnsureDirExists(extractPath); err != nil {
		return "", err
	}

	// 解压所有文件
	for _, f := range r.File {
		fpath := filepath.Join(extractPath, f.Name)

		// 创建目录
		if f.FileInfo().IsDir() {
			EnsureDirExists(fpath)
			continue
		}

		// 创建文件目录
		if err := EnsureDirExists(filepath.Dir(fpath)); err != nil {
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

// RenameDir 重命名目录并移动到当前工作目录
func RenameDir(oldPath, newName string) error {
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("源目录不存在: %s", oldPath)
	}

	newPath := filepath.Join(filepath.Dir(oldPath), newName)

	fmt.Printf("重命名: %s -> %s\n", filepath.Base(oldPath), newName)
	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Printf("直接重命名失败，尝试跨磁盘移动: %v\n", err)
	}

	currentDir, err := GetWorkingDir()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	finalPath := filepath.Join(currentDir, newName)
	fmt.Printf("移动目录到: %s\n", finalPath)

	return MoveDir(newPath, finalPath)
}

// MergeDir 将 src 目录的内容合并到 dst 目录。
// 同名文件会被覆盖，不会删除 dst 中独有的文件。
func MergeDir(src, dst string) error {
	return copyDir(src, dst)
}
