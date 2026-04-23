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

	// Arguments
	if len(c._arguments) > 0 {
		b.WriteString("\nArguments:\n")
		for _, arg := range c._arguments {
			b.WriteString(fmt.Sprintf("  %-20s %s\n", arg.name, arg.desc))
		}
	}

	// Options
	if len(c._options) > 0 {
		b.WriteString("\nOptions:\n")
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
			b.WriteString(fmt.Sprintf("  %-20s %s\n", flag, opt.desc))
		}
	}

	// Commands
	if len(c._subCommands) > 0 {
		b.WriteString("\nCommands:\n")
		for _, sub := range c._subCommands {
			nameCol := sub.name
			for _, arg := range sub._arguments {
				if arg.valueRequired {
					nameCol += fmt.Sprintf(" <%s>", arg.name)
				} else {
					nameCol += fmt.Sprintf(" [%s]", arg.name)
				}
			}
			b.WriteString(fmt.Sprintf("  %-20s %s\n", nameCol, sub.description))
		}
	}

	return b.String()
}
