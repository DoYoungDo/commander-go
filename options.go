package commandergo

import (
	"fmt"
	"regexp"
)

type Options []*Option

func (o Options) has(name string) bool {
	for _, opt := range o {
		if opt.name == name {
			return true
		}
	}
	return false
}

func (o Options) get(name string) (*Option, bool) {
	for _, opt := range o {
		if opt.name == name {
			return opt, true
		}
	}
	return nil, false
}

type Option struct {
	name         string
	alias        string
	desc         string
	defaultValue Varaint

	valueRequired bool
	valueName     string
}

var OPTION_FLAG_PATTERN = regexp.MustCompile(`^\s*(?:(?:-([a-zA-Z])(?:(?:\s+)|(?:\s*,\s*))\-\-([a-zA-Z-]+)\s+(?:\[\s*([a-zA-Z]+)\s*\]|<\s*([a-zA-Z]+)\s*>))|(?:-([a-zA-Z])(?:(?:\s+)|(?:\s*,\s*))\-\-([a-zA-Z-]+))|(?:\-\-([a-zA-Z-]+)\s+(?:(?:\[\s*([a-zA-Z]+)\s*\])|(?:\<\s*([a-zA-Z]+)\s*\>)))|(?:\-\-([a-zA-Z-]+)))\s*$`)

/**
* ps: --a | -a --abc | --a [valueName...] | -a --abc <valueName...>
 */
func NewOption(flag string) (*Option, error) {
	if !OPTION_FLAG_PATTERN.MatchString(flag) {
		return nil, fmt.Errorf("invalid option flag :%v", flag)
	}
	group := OPTION_FLAG_PATTERN.FindStringSubmatch(flag)
	tryGroup := func(indexs []int) string {
		for _, index := range indexs {
			if group[index] != "" {
				return group[index]
			}
		}
		return ""
	}

	name := func() string {
		n := tryGroup([]int{2, 6, 7, 10})
		if n != "" {
			return n
		}
		panic(fmt.Errorf("invalid option flag :%v", flag))
	}()
	alias := tryGroup([]int{1, 5})
	valueName := tryGroup([]int{3, 4, 8, 9})
	valueRequired := group[4] != "" || group[9] != ""

	return &Option{
		name:          name,
		alias:         alias,
		valueName:     valueName,
		valueRequired: valueRequired,
	}, nil
}

func (c *Command) options(flag, desc string, defaultValue any) *Command {
	option, err := NewOption(flag)
	if err != nil {
		panic(err)
	}

	if c._options.has(option.name) {
		panic(fmt.Errorf("option %s already exists", option.name))
	}

	option.desc = desc
	option.defaultValue = Varaint{value: defaultValue}
	c._options = append(c._options, option)

	return c
}
