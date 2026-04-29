package commandergo

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	reLongOpt  = regexp.MustCompile(`^--([a-zA-Z][a-zA-Z-]*)(?:=(.+))?$`)
	reShortOpt = regexp.MustCompile(`^-([a-zA-Z]+)(?:=(.+))?$`)
	reCommand  = regexp.MustCompile(`^([a-zA-Z][a-zA-Z\d]*)$`)
)

func parseValue(s string) Varaint {
	if i, err := strconv.Atoi(s); err == nil {
		return Varaint{value: i}
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return Varaint{value: f}
	}
	if s == "true" {
		return Varaint{value: true}
	}
	if s == "false" {
		return Varaint{value: false}
	}
	return Varaint{value: s}
}

func (c *Command) findOptionByName(name string) *Option {
	if opt, ok := c._options.get(name); ok {
		return opt
	}
	return nil
}

func (c *Command) findOptionByAlias(alias string) *Option {
	for _, opt := range c._options {
		if opt.alias == alias {
			return opt
		}
	}
	return nil
}

func (c *Command) parse(args []string) error {
	ctx := newContext(c)
	i := 0

	for i < len(args) {
		token := args[i]

		// --name 或 --name=value
		if m := reLongOpt.FindStringSubmatch(token); m != nil {
			name, inlineVal := m[1], m[2]
			opt := c.findOptionByName(name)
			if opt == nil {
				fmt.Fprintf(os.Stderr, "warning: unknown option: --%s\n", name)
				i++
				continue
			}
			i++
			if opt.valueName != "" {
				val := inlineVal
				if val == "" {
					if opt.valueRequired {
						if i >= len(args) {
							if !opt.defaultValue.IsEmpty() {
								ctx.parsedOpts[name] = opt.defaultValue
								continue
							}
							return fmt.Errorf("option --%s requires a value", name)
						}
						val = args[i]
						i++
					} else {
						if !opt.defaultValue.IsEmpty() {
							ctx.parsedOpts[name] = opt.defaultValue
							continue
						}
					}
				}
				if val != "" {
					ctx.parsedOpts[name] = parseValue(val)
				}
			} else {
				ctx.parsedOpts[name] = Varaint{value: true}
			}
			continue
		}

		// -f 或 -f=value 或 -abc
		if m := reShortOpt.FindStringSubmatch(token); m != nil {
			aliases, inlineVal := m[1], m[2]
			i++
			// 多别名合并：前 n-1 个只能是布尔
			for j := 0; j < len(aliases)-1; j++ {
				a := string(aliases[j])
				opt := c.findOptionByAlias(a)
				if opt == nil {
					fmt.Fprintf(os.Stderr, "warning: unknown option: -%s\n", a)
					continue
				}
				ctx.parsedOpts[opt.name] = Varaint{value: true}
			}
			// 最后一个可带值
			last := string(aliases[len(aliases)-1])
			opt := c.findOptionByAlias(last)
			if opt == nil {
				fmt.Fprintf(os.Stderr, "warning: unknown option: -%s\n", last)
				continue
			}
			if opt.valueName != "" {
				val := inlineVal
				if val == "" && opt.valueRequired {
					if i >= len(args) {
						return fmt.Errorf("option -%s requires a value", last)
					}
					val = args[i]
					i++
				}
				if val != "" {
					ctx.parsedOpts[opt.name] = parseValue(val)
				}
			} else {
				ctx.parsedOpts[opt.name] = Varaint{value: true}
			}
			continue
		}

		// 普通 token：子命令或位置参数
		if m := reCommand.FindStringSubmatch(token); m != nil {
			if sub, ok := c._subCommands.get(m[1]); ok {
				return sub.parse(args[i+1:])
			}
		}
		// 位置参数：按 _arguments 顺序填入
		ctx.parsedArgs = append(ctx.parsedArgs, parseValue(token))
		i++
	}

	// help / version 优先，不触发 action
	if _, ok := ctx.parsedOpts["help"]; ok {
		fmt.Print(c.helpText())
		return nil
	}
	if _, ok := ctx.parsedOpts["version"]; ok {
		fmt.Println(c.version)
		return nil
	}

	// 检查必填 argument
	for i, arg := range c._arguments {
		if arg.valueRequired {
			if i >= len(ctx.Args()) {
				return fmt.Errorf("argument <%s> is required", arg.name)
			}
		}
	}

	if c.actionFn != nil {
		c.actionFn(ctx)
	}
	return nil
}
