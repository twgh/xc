# xc

<p>
	<a href="https://github.com/twgh/xc/releases"><img src="https://img.shields.io/badge/release-0.0.1-blue" alt="release"></a>
	<a href="https://golang.org"> <img src="https://img.shields.io/badge/golang-≥1.18-blue" alt="golang"></a>
	<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-brightgreen" alt="License"></a>
</p>

## 介绍

xc 是一个 xcgui 助手类型的命令行工具, 功能包括给项目添加 xcgui、编译程序、下载 xcgui 和 example 仓库的源码 ZIP、下载 xcgui.dll 文件等。

## 安装

```
go install -ldflags="-s -w" -trimpath github.com/twgh/xc/cmd@latest
```

成功则 `%GOPATH%\bin` 目录中会出现一个 `xc.exe`

## 使用方法

### 查看帮助信息

```bash
xc --help
```

输出:

```
xc 是一个 xcgui 助手类型的命令行工具, 功能包括给项目添加 xcgui、编译程序、
下载 xcgui 和 example 仓库的源码 ZIP、下载 xcgui.dll 文件等。

使用方法:
  xc [command]

可用命令:
  get           执行 go get -u github.com/twgh/xcgui
  build         执行 go build -ldflags="-s -w -H windowsgui" -trimpath
  zipdownload   下载并解压 xcgui 和 example 仓库的源码 ZIP
  dlldownload   下载 xcgui.dll 文件
  version       显示版本信息
  help          显示命令帮助信息

使用 "xc [command] --help" 获取更多关于某个命令的信息。
```

