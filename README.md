# Commander-Go

> 中文 | [English](README_EN.md)

Commander-Go 是一个仿照 [commander.js](https://github.com/tj/commander.js) 实现的 Go 命令行解析库。它提供了简洁的链式 API，用于定义命令、选项和参数，支持子命令、默认值、自动生成帮助信息等功能。

## 特性

- 🎯 链式 API 设计，简洁易用
- 🔧 支持选项（布尔值、单值、可选值）
- 📝 支持参数（必需参数、可选参数、多值参数）
- 🌲 支持子命令和嵌套命令
- ⚙️ 支持默认值与值类型自动推断（int、float64、bool、string）
- 📖 自动生成帮助信息（动态列宽对齐，含完整命令路径）
- 🔍 分层错误处理（未知选项 warning，缺少必填值 error）
- 🎨 支持选项别名和组合（如 `-vd`）

## 安装

```bash
go get github.com/DoYoungDo/commander-go
```

## 快速开始

```go
package main

import (
    "fmt"
    "os"
    commander "github.com/DoYoungDo/commander-go"
)

func main() {
    err := commander.New("example").
        Version("1.0.0").
        Description("一个示例命令行应用").
        Options("-n, --name <name>", "你的名字", "").
        Action(func(ctx *commander.Context) {
            name := ctx.Opt("name").ToString()
            if name != "" {
                fmt.Println("Hello,", name)
            } else {
                fmt.Println("Hello, World!")
            }
        }).
        Parse(os.Args[1:])
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
}
```

运行示例：

```bash
$ ./example -n Alice
Hello, Alice!

$ ./example -V
1.0.0

$ ./example --help
Usage: example [options]

一个示例命令行应用

Options:
  -V, --version  output the version number
  -n, --name     你的名字
  -h, --help     display help for command
```

## 核心概念

### 选项语法

| 格式 | 说明 |
|------|------|
| `--option` | 布尔选项 |
| `-o, --option` | 带别名的布尔选项 |
| `--option <value>` | 必需值选项 |
| `--option [value]` | 可选值选项 |

### 参数语法

| 格式 | 说明 |
|------|------|
| `<arg>` | 必需参数 |
| `[arg]` | 可选参数 |
| `<arg...>` | 多值必需参数 |
| `[arg...]` | 多值可选参数 |

## 详细用法

### 1. 定义选项

```go
commander.New("app").
    Options("-v, --verbose", "详细输出", false).
    Options("-o, --output <file>", "输出文件", "").
    Options("-n, --num <n>", "数字", 0).
    Action(func(ctx *commander.Context) {
        if ctx.Opt("verbose").ToBool() {
            fmt.Println("verbose mode")
        }
        fmt.Println("output:", ctx.Opt("output").ToString())
        fmt.Println("num:", ctx.Opt("num").ToInt())
    }).
    Parse(os.Args[1:])
```

### 2. 定义参数

```go
commander.New("copy").
    Arguments("<from>", "源文件", nil).
    Arguments("[to]", "目标文件", nil).
    Action(func(ctx *commander.Context) {
        fmt.Println("from:", ctx.Arg("from").ToString())
        fmt.Println("to:", ctx.Arg("to").ToString())
    }).
    Parse(os.Args[1:])
```

### 3. 定义子命令

```go
app := commander.New("git").
    Version("1.0.0").
    Description("Git 命令行工具")

app.Command("add", "添加文件到暂存区").
    Arguments("<files...>", "文件列表", nil).
    Options("-f, --force", "强制添加", false).
    Action(func(ctx *commander.Context) {
        fmt.Println("add:", ctx.Arg("files").ToString())
    })

app.Command("commit", "提交更改").
    Options("-m, --message <msg>", "提交信息", "").
    Action(func(ctx *commander.Context) {
        fmt.Println("commit:", ctx.Opt("message").ToString())
    })

app.Parse(os.Args[1:])
```

### 4. 选项别名组合

支持短选项组合（类似 `tar -xzvf`）：

```bash
$ ./app -vd    # 同时启用 verbose 和 debug
```

### 5. 值类型自动推断

Parse 时自动将字符串转换为对应类型：

```bash
$ ./app --num 42      # → int
$ ./app --rate 3.14   # → float64
$ ./app --flag true   # → bool
$ ./app --name hello  # → string
```

### 6. 错误处理

```go
if err := app.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, "error:", err)
    os.Exit(1)
}
```

- **未知选项**：打印 warning 到 stderr，继续解析，不影响 action 执行
- **缺少必填选项值**：返回 error，终止解析
- **缺少必填参数**：返回 error，终止解析

## 完整示例

见 [example/todo](example/todo/main.go)，一个完整的 todo CLI：

```bash
$ cd example/todo
$ go run main.go --help
$ go run main.go add 买菜
$ go run main.go list
$ go run main.go add -h
```

## 测试

```bash
go test ./...
```

## 目录结构

```
commander-go/
├── argument.go         # 参数定义与解析
├── command.go          # 子命令定义
├── commander_go.go     # 核心 Command 结构与链式 API
├── context.go          # 运行时值存储（Context）
├── help.go             # 帮助文本生成
├── options.go          # 选项定义与解析
├── parse.go            # 命令行参数解析核心逻辑
├── varaint.go          # 动态值类型
├── example/
│   └── todo/           # 完整 todo CLI 示例
└── *_test.go           # 测试文件
```

## 许可证

MIT License
