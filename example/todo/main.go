package main

import (
	"fmt"
	"os"

	commander "github.com/DoYoungDo/commander-go"
)

func main() {
	app := commander.New("todo").
		Description("待办项").
		Version("0.0.1").
		Arguments("[todo...]", "添加 待办项", nil)

	// add
	app.Command("add", "添加 待办项").
		Arguments("<todo...>", "待办项", nil).
		Options("-d, --done", "添加时完成", false).
		Action(func(ctx *commander.Context) {
			fmt.Println("add todo:", ctx.Arg("todo").ToString())
		})

	// rm
	app.Command("rm", "删除 待办项").
		Arguments("<index...>", "索引序号", nil).
		Action(func(ctx *commander.Context) {
			fmt.Println("rm todo:", ctx.Arg("index").ToString())
		})

	// mod
	app.Command("mod", "修改 待办项").
		Arguments("<index>", "索引序号", nil).
		Arguments("[todo]", "待办内容", nil).
		Options("-a, --append", "在原内容上追加", false).
		Options("-i, --insert", "在原内容上头插", false).
		Options("-d, --done [done]", "修改为完成", nil).
		Options("-p, --priority <priority>", "设置优先级，取值1-5", nil).
		Action(func(ctx *commander.Context) {
			fmt.Println("mod todo index:", ctx.Arg("index").ToString())
		})

	// list
	app.Command("list", "显示 待办项").
		Arguments("[range]", "显示范围", nil).
		Options("-d, --done [done]", "只显示完成的", nil).
		Options("-c, --count", "显示数量", false).
		Action(func(ctx *commander.Context) {
			fmt.Println("list todos")
		})

	// done
	app.Command("done", "完成 待办项").
		Arguments("<index...>", "索引序号", nil).
		Action(func(ctx *commander.Context) {
			fmt.Println("done todo:", ctx.Arg("index").ToString())
		})

	// mv
	app.Command("mv", "移动 待办项").
		Arguments("<index>", "待移动的待办项索引序号", nil).
		Arguments("<distindex>", "目标索引序号", nil).
		Action(func(ctx *commander.Context) {
			fmt.Println("mv todo:", ctx.Arg("index").ToString(), "->", ctx.Arg("distindex").ToString())
		})

	// find
	app.Command("find", "查找 待办项").
		Arguments("<todo...>", "查找内容", nil).
		Options("-c, --caseSensitive", "区分大小写", false).
		Options("-s, --single", "匹配单个条件", false).
		Options("-d, --done [done]", "匹配完成待办", nil).
		Action(func(ctx *commander.Context) {
			fmt.Println("find todo:", ctx.Arg("todo").ToString())
		})

	// clear
	app.Command("clear", "清空 待办项").
		Action(func(ctx *commander.Context) {
			fmt.Println("clear all todos")
		})

	// 根命令 action：todo xxx 等同于 todo add xxx
	app.Action(func(ctx *commander.Context) {
		todo := ctx.Arg("todo")
		if todo.IsString() && todo.ToString() != "" {
			fmt.Println("add todo:", todo.ToString())
		}
	})

	if err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
