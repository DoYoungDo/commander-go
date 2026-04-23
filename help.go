package commandergo

import (
	"fmt"
	"strings"
)

func (c *Command) helpText() string {
	var b strings.Builder

	// Usage line
	usage := "Usage: " + c.name
	if len(c._options) > 0 {
		usage += " [options]"
	}
	if len(c._subCommands) > 0 {
		usage += " [command]"
	}
	for _, arg := range c._arguments {
		if arg.valueRequired {
			if arg.multiValue {
				usage += fmt.Sprintf(" <%s...>", arg.name)
			} else {
				usage += fmt.Sprintf(" <%s>", arg.name)
			}
		} else {
			if arg.multiValue {
				usage += fmt.Sprintf(" [%s...]", arg.name)
			} else {
				usage += fmt.Sprintf(" [%s]", arg.name)
			}
		}
	}
	b.WriteString(usage + "\n")

	if c.description != "" {
		b.WriteString("\n" + c.description + "\n")
	}

	// 动态列宽对齐辅助
	type row struct{ left, right string }
	writeRows := func(title string, rows []row) {
		if len(rows) == 0 {
			return
		}
		maxLen := 0
		for _, r := range rows {
			if len(r.left) > maxLen {
				maxLen = len(r.left)
			}
		}
		b.WriteString("\n" + title + "\n")
		for _, r := range rows {
			b.WriteString(fmt.Sprintf("  %-*s  %s\n", maxLen, r.left, r.right))
		}
	}

	// Arguments
	var argRows []row
	for _, arg := range c._arguments {
		argRows = append(argRows, row{arg.name, arg.desc})
	}
	writeRows("Arguments:", argRows)

	// Options
	var optRows []row
	for _, opt := range c._options {
		flag := "--" + opt.name
		if opt.alias != "" {
			flag = "-" + opt.alias + ", " + flag
		}
		if opt.valueName != "" {
			if opt.valueRequired {
				flag += fmt.Sprintf(" <%s>", opt.valueName)
			} else {
				flag += fmt.Sprintf(" [%s]", opt.valueName)
			}
		}
		optRows = append(optRows, row{flag, opt.desc})
	}
	writeRows("Options:", optRows)

	// Commands
	var cmdRows []row
	for _, sub := range c._subCommands {
		nameCol := sub.name
		for _, arg := range sub._arguments {
			if arg.valueRequired {
				nameCol += fmt.Sprintf(" <%s>", arg.name)
			} else {
				nameCol += fmt.Sprintf(" [%s]", arg.name)
			}
		}
		cmdRows = append(cmdRows, row{nameCol, sub.description})
	}
	writeRows("Commands:", cmdRows)

	return b.String()
}
