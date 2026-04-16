package commandergo

import (
	"fmt"
	"regexp"
)

type Options map[string]*Option

func (o Options) has(name string) bool {
	_, ok := o[name]
	return ok
}

type Option struct {
	name         string
	alias        string
	desc         string
	defaultValue Varaint

	multiValue    bool
	valueRequired bool
	valueName     string
}

var OPTION_FLAT_PATTERN = regexp.MustCompile(`^\s*(?:(?:-([a-zA-Z])(?:(?:\s+)|(?:\s*,\s*))\-\-([a-zA-Z-]+)\s+(?:\[([a-zA-Z]+)(\.\.\.)?\]|<([a-zA-Z]+)(\.\.\.)?>))|(?:-([a-zA-Z])(?:(?:\s+)|(?:\s*,\s*))\-\-([a-zA-Z-]+))|(?:\-\-([a-zA-Z-]+)\s+(?:(?:\[([a-zA-Z]+)(\.\.\.)?\])|(?:\<([a-zA-Z]+)(\.\.\.)?\>)))|(?:\-\-([a-zA-Z-]+)))\s*$`)

/**
* ps: --a | -a --abc | --a [valueName...] | -a --abc <valueName...>
 */
func NewOption(flag string) (*Option, error) {
	if !OPTION_FLAT_PATTERN.MatchString(flag) {
		return nil, fmt.Errorf("invalid option flag :%v", flag)
	}
	group := OPTION_FLAT_PATTERN.FindStringSubmatch(flag)
	tryGroup := func(indexs []int) string {
		for _, index := range indexs {
			if group[index] != "" {
				return group[index]
			}
		}
		return ""
	}

	name := func() string {
		n := tryGroup([]int{2, 8, 9, 14})
		if n != "" {
			return n
		}
		panic(fmt.Errorf("invalid option flag :%v", flag))
	}()
	alias := tryGroup([]int{1, 7})
	valueName := tryGroup([]int{3, 5, 10, 12})
	multiValue := group[4] != "" || group[6] != "" || group[11] != "" || group[13] != ""
	valueRequired := group[5] != "" || group[12] != ""

	return &Option{
		name:          name,
		alias:         alias,
		valueName:     valueName,
		multiValue:    multiValue,
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
	c._options[option.name] = option

	return c
}
