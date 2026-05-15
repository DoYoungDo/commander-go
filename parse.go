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

	FLAG_END = "--"
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

type parseOption struct {
	strict bool
}

func (c *Command) parse(args []string, parseOpt *parseOption) error {
	ctx := newContext(c)
	declareArgLen := len(c._arguments)

	flagEnd := false
	for i := 0; i < len(args); i++ {
		token := args[i]

		// 如果是--，就直接跳过
		if token == FLAG_END {
			flagEnd = true
			continue
		}

		// flag end后的所有token都作为参数
		if flagEnd {
			ctx.parsedArgs = append(ctx.parsedArgs, parseValue(token))
			continue
		}

		// 如果命中子命令，就直接交给子命令解析
		if m := reCommand.FindStringSubmatch(token); m != nil {
			if sub, ok := c._subCommands.get(m[1]); ok {
				return sub.parse(args[i+1:], parseOpt)
			}
		}

		// --name 或 --name=value
		if m := reLongOpt.FindStringSubmatch(token); m != nil {
			name, inlineVal := m[1], m[2]
			opt := c.findOptionByName(name)
			if opt == nil {
				if parseOpt.strict {
					return fmt.Errorf("unknown option: %v", token)
				} else {
					fmt.Fprintf(os.Stderr, "warning: unknown option: --%s\n", name)
					continue
				}
			}
			// option不需要值的情况：不接收cli值，也不接收默认值
			if opt.valueName == "" {
				ctx.parsedOpts[name] = Varaint{value: true}
				if inlineVal != "" && parseOpt.strict {
					return fmt.Errorf("option %v does not accept a value", token)
				}
				continue
			}

			// 如果有直接的值，就直接解析
			if inlineVal != "" {
				ctx.parsedOpts[name] = parseValue(inlineVal)
				continue
			}

			// 如果cli没有传值，尝试解析后一个token做为值,token尝试失败再尝试使用默认值
			if nextI := i + 1; nextI < len(args) {
				nextToken := args[nextI]
				// 检验下一个token
				if nextToken != FLAG_END &&
					!reLongOpt.MatchString(nextToken) &&
					!reShortOpt.MatchString(nextToken) {
					ctx.parsedOpts[name] = parseValue(nextToken)
					i = nextI
					continue
				}
			}

			// 没有下一个token,或者下一个token检验失败，尝试使用默认值
			if !opt.defaultValue.IsEmpty() {
				ctx.parsedOpts[name] = opt.defaultValue
				continue
			}

			// 以上都失败，如果这个值是必选的，直接结束
			if opt.valueRequired {
				return fmt.Errorf("option %v requires a value", token)
			} else {
				ctx.parsedOpts[name] = Varaint{value: true}
				continue
			}
		}

		// -f 或 -f=value 或 -abc
		if m := reShortOpt.FindStringSubmatch(token); m != nil {
			aliases, inlineVal := m[1], m[2]
			// 多别名合并：前 n-1 个只能是布尔
			for j := 0; j < len(aliases)-1; j++ {
				a := string(aliases[j])
				opt := c.findOptionByAlias(a)
				// 严格模式，如果别名不存在，就直接结束
				if opt == nil {
					if parseOpt.strict {
						return fmt.Errorf("unknown option alias: -%s", a)
					} else {
						fmt.Fprintf(os.Stderr, "warning: unknown option: -%s\n", a)
						continue
					}
				}

				// option不需要值的情况：不接收cli值，也不接收默认值
				if opt.valueName == "" {
					ctx.parsedOpts[opt.name] = Varaint{value: true}
					continue
				}

				// 尝试使用默认值
				if !opt.defaultValue.IsEmpty() {
					ctx.parsedOpts[opt.name] = opt.defaultValue
					continue
				}

				// 以上都失败，如果这个值是必选的，直接结束
				if opt.valueRequired {
					return fmt.Errorf("option %v requires a value", token)
				} else {
					ctx.parsedOpts[opt.name] = Varaint{value: true}
					continue
				}
			}
			// 最后一个可带值
			last := string(aliases[len(aliases)-1])
			opt := c.findOptionByAlias(last)
			// 严格模式，如果别名不存在，就直接结束
			if opt == nil {
				if parseOpt.strict {
					return fmt.Errorf("unknown option alias: -%s", last)
				} else {
					fmt.Fprintf(os.Stderr, "warning: unknown option: -%s\n", last)
					continue
				}
			}

			// option不需要值的情况：不接收cli值，也不接收默认值
			if opt.valueName == "" {
				ctx.parsedOpts[opt.name] = Varaint{value: true}
				if inlineVal != "" && parseOpt.strict {
					return fmt.Errorf("option %v does not accept a value", token)
				}
				continue
			}

			// 如果有直接的值，就直接解析
			if inlineVal != "" {
				ctx.parsedOpts[opt.name] = parseValue(inlineVal)
				continue
			}

			// 如果cli没有传值，尝试解析后一个token做为值,token尝试失败再尝试使用默认值
			if nextI := i + 1; nextI < len(args) {
				nextToken := args[nextI]
				// 检验下一个token
				if nextToken != FLAG_END &&
					!reLongOpt.MatchString(nextToken) &&
					!reShortOpt.MatchString(nextToken) {
					ctx.parsedOpts[opt.name] = parseValue(nextToken)
					i = nextI
					continue
				}
			}

			// 没有下一个token,或者下一个token检验失败，尝试使用默认值
			if !opt.defaultValue.IsEmpty() {
				ctx.parsedOpts[opt.name] = opt.defaultValue
				continue
			}

			// 以上都失败，如果这个值是必选的，直接结束
			if opt.valueRequired {
				return fmt.Errorf("option %v requires a value", token)
			} else {
				ctx.parsedOpts[opt.name] = Varaint{value: true}
				continue
			}
		}

		ctx.parsedArgs = append(ctx.parsedArgs, parseValue(token))
	}

	checkError := func() error {
		// 位置参数：按 _arguments 顺序填入
		if parseArgLen := len(ctx.parsedArgs); declareArgLen == 0 && parseArgLen > 0 {
			return fmt.Errorf("argument is not expected, but got %v", parseArgLen)
		} else if declareArgLen > 0 && parseArgLen > declareArgLen && !c._arguments[declareArgLen-1].multiValue {
			return fmt.Errorf("%v arguments is expected, but got %v", declareArgLen, parseArgLen)
		}
		return nil
	}

	if parseOpt.strict {
		if err := checkError(); err != nil {
			return err
		}
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
