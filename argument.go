package commandergo

import (
	"fmt"
	"regexp"
)

type Arguments []*Argument

func (a Arguments) has(name string) bool {
	for _, arg := range a {
		if arg.name == name {
			return true
		}
	}
	return false
}

type Argument struct {
	name          string
	desc          string
	defaultValue  Varaint
	multiValue    bool
	valueRequired bool
}

var ARGUMENT_FLAG_PATTERN = regexp.MustCompile(`^\s*(?:(?:\[([a-zA-Z][a-zA-Z\d]+)(\.\.\.)?\])|(?:<([a-zA-Z][a-zA-Z\d]+)(\.\.\.)?>))\s*$`)

/**
 * ps: [name] | [name...] | <name> | <name...>
 */

func NewArgument(text string) (*Argument, error) {
	if !ARGUMENT_FLAG_PATTERN.MatchString(text) {
		return nil, fmt.Errorf("invalid argument name :%v", text)
	}

	group := ARGUMENT_FLAG_PATTERN.FindStringSubmatch(text)
	if group == nil {
		return nil, fmt.Errorf("invalid argument name :%v", text)
	}
	name := func() string {
		if group[1] != "" {
			return group[1]
		}
		if group[3] != "" {
			return group[3]
		}
		panic(fmt.Errorf("invalid argument name :%v", text))
	}()
	multiValue := group[2] != "" || group[4] != ""
	valueRequired := group[3] != ""

	return &Argument{
		name:          name,
		multiValue:    multiValue,
		valueRequired: valueRequired,
	}, nil
}

func (c *Command) arguments(name, desc string, defaultValue any) *Command {
	arg, err := NewArgument(name)
	if err != nil {
		panic(err)
	}
	if c._arguments.has(arg.name) {
		panic(fmt.Errorf("argument %s already exists", arg.name))
	}

	arg.desc = desc
	arg.defaultValue = Varaint{value: defaultValue}
	c._arguments = append(c._arguments, arg)

	return c
}
