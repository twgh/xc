package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
