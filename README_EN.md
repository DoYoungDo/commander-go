# Commander-Go

> [中文](README.md) | English

Commander-Go is a Go command-line parsing library inspired by [commander.js](https://github.com/tj/commander.js). It provides a clean, chainable API for defining commands, options, and arguments, with support for subcommands, default values, automatic help generation, and more.

## Features

- 🎯 Chainable API design, simple and easy to use
- 🔧 Support for options (boolean, single-value, optional-value)
- 📝 Support for arguments (required, optional, multi-value)
- 🌲 Support for subcommands and nested commands
- ⚙️ Support for default values and automatic type inference (int, float64, bool, string)
- 📖 Automatic help generation (dynamic column alignment, full command path)
- 🔍 Layered error handling (unknown option warning, missing required value error)
- 🎨 Support for option aliases and combinations (like `-vd`)

## Installation

```bash
go get github.com/DoYoungDo/commander-go
```

## Quick Start

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
        Description("An example command-line application").
        Options("-n, --name <name>", "Your name", "").
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

Usage example:

```bash
$ ./example -n Alice
Hello, Alice!

$ ./example -V
1.0.0

$ ./example --help
Usage: example [options]

An example command-line application

Options:
  -V, --version  output the version number
  -n, --name     Your name
  -h, --help     display help for command
```

## Core Concepts

### Option Syntax

| Format | Description |
|--------|-------------|
| `--option` | Boolean option |
| `-o, --option` | Boolean option with alias |
| `--option <value>` | Option with required value |
| `--option [value]` | Option with optional value |

### Argument Syntax

| Format | Description |
|--------|-------------|
| `<arg>` | Required argument |
| `[arg]` | Optional argument |
| `<arg...>` | Multi-value required argument |
| `[arg...]` | Multi-value optional argument |

## Detailed Usage

### 1. Defining Options

```go
commander.New("app").
    Options("-v, --verbose", "Verbose output", false).
    Options("-o, --output <file>", "Output file", "").
    Options("-n, --num <n>", "Number", 0).
    Action(func(ctx *commander.Context) {
        if ctx.Opt("verbose").ToBool() {
            fmt.Println("verbose mode")
        }
        fmt.Println("output:", ctx.Opt("output").ToString())
        fmt.Println("num:", ctx.Opt("num").ToInt())
    }).
    Parse(os.Args[1:])
```

### 2. Defining Arguments

```go
commander.New("copy").
    Arguments("<from>", "Source file", nil).
    Arguments("[to]", "Target file", nil).
    Action(func(ctx *commander.Context) {
        fmt.Println("from:", ctx.Arg("from").ToString())
        fmt.Println("to:", ctx.Arg("to").ToString())
    }).
    Parse(os.Args[1:])
```

### 3. Defining Subcommands

```go
app := commander.New("git").
    Version("1.0.0").
    Description("Git command-line tool")

app.Command("add", "Add files to staging area").
    Arguments("<files...>", "File list", nil).
    Options("-f, --force", "Force add", false).
    Action(func(ctx *commander.Context) {
        fmt.Println("add:", ctx.Arg("files").ToString())
    })

app.Command("commit", "Commit changes").
    Options("-m, --message <msg>", "Commit message", "").
    Action(func(ctx *commander.Context) {
        fmt.Println("commit:", ctx.Opt("message").ToString())
    })

app.Parse(os.Args[1:])
```

### 4. Option Combinations

Support for short option combinations (like `tar -xzvf`):

```bash
$ ./app -vd    # Enable both verbose and debug
```

### 5. Automatic Type Inference

Values are automatically converted to the appropriate type during Parse:

```bash
$ ./app --num 42      # → int
$ ./app --rate 3.14   # → float64
$ ./app --flag true   # → bool
$ ./app --name hello  # → string
```

### 6. Error Handling

```go
if err := app.Parse(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, "error:", err)
    os.Exit(1)
}
```

- **Unknown option**: prints warning to stderr, continues parsing, action still executes
- **Missing required option value**: returns error, stops parsing
- **Missing required argument**: returns error, stops parsing

## Complete Example

See [example/todo](example/todo/main.go) for a complete todo CLI:

```bash
$ cd example/todo
$ go run main.go --help
$ go run main.go add 买菜
$ go run main.go list
$ go run main.go add -h
```

## Testing

```bash
go test ./...
```

## Directory Structure

```
commander-go/
├── argument.go         # Argument definition and parsing
├── command.go          # Subcommand definition
├── commander_go.go     # Core Command struct and chainable API
├── context.go          # Runtime value storage (Context)
├── help.go             # Help text generation
├── options.go          # Option definition and parsing
├── parse.go            # Core command-line parsing logic
├── varaint.go          # Dynamic value type
├── example/
│   └── todo/           # Complete todo CLI example
└── *_test.go           # Test files
```

## License

MIT License
