# xc

<p>
	<a href="https://github.com/twgh/xc/releases"><img src="https://img.shields.io/badge/release-0.0.7-blue" alt="release"></a>
	<a href="https://golang.org"> <img src="https://img.shields.io/badge/golang-≥1.18-blue" alt="golang"></a>
	<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-brightgreen" alt="License"></a>
</p>
## 介绍

xc 是一个 xcgui 助手类型的命令行工具, 也包含一些其它的功能, 比如: 从 GitHub 仓库克隆代码、下载文件等(自动使用代理)。

## 功能

- 给项目添加 xcgui
- 编译程序
- 给程序添加 Windows 资源 / 版本信息 / 程序清单
- 下载 xcgui 和 xcgui-example 仓库的源码
- 下载 go-xcgui-dev 技能
- 下载 xcgui.dll 文件
- 克隆任意 GitHub 仓库
- 从链接下载文件, 加速下载 Github文件

## 安装

```
go install github.com/twgh/xc@latest
```

成功则 `%GOPATH%\bin` 目录中会出现一个 `xc.exe`

## 使用方法

### 查看帮助信息

```bash
xc -h
```

输出:

```
xc 是一个 xcgui 助手类型的命令行工具, 功能包括给项目添加 xcgui、编译程序、给程序添加 Windows 资源 / 版本信息 / 程序清单、下载 xcgui 和 xcgui-example 仓库的源码、下载 xcgui.dll 文件, 克隆任意 GitHub 仓库, 从链接下载文件, 下载 go-xcgui-dev 技能。

使用方法:
  xc [command]

可用命令:
  get           执行 go get -u github.com/twgh/xcgui
  build         执行 go build -ldflags="-s -w -H windowsgui" -trimpath
  zipdl         下载并解压 xcgui 和 xcgui-example 仓库的源码 ZIP，自动选择可用代理
  clone         克隆任意 GitHub 仓库，自动选择可用代理
  dl            从链接下载文件，如果是 GitHub 文件会自动选择可用代理
  dlldl         下载 xcgui.dll 文件
  getskill      下载 go-xcgui-dev 技能
  res           执行 go-winres 命令（添加 Windows 资源 / 版本信息 / 程序清单）
  version       显示版本信息
  help          显示命令帮助信息

使用 "xc [command] --help" 获取更多关于某个命令的信息。
```
