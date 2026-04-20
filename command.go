package commandergo

import (
	"fmt"
	"regexp"
)

type Commands map[string]*Command

func (o Commands) has(name string) bool {
	_, ok := o[name]
	return ok
}

var CONMMAND_FLAG_PATTERN = regexp.MustCompile(`^\s*([a-zA-Z][a-zA-Z\d]+)\s*((?:\[[a-zA-Z][a-zA-Z\d]+(?:\.\.\.)?\])|(?:<[a-zA-Z][a-zA-Z\d]+(?:\.\.\.)?>))?\s*$`)

func newCommandWithNameAndArg(nameAndArg string) (*Command, error) {
	if !CONMMAND_FLAG_PATTERN.MatchString(nameAndArg) {
		return nil, fmt.Errorf("invalid option nameAndArg :%v", nameAndArg)
	}
	group := CONMMAND_FLAG_PATTERN.FindStringSubmatch(nameAndArg)
	tryGroup := func(indexs []int) string {
		for _, index := range indexs {
			if group[index] != "" {
				return group[index]
			}
		}
		return ""
	}
	name := tryGroup([]int{1})
	if name == "" {
		panic(fmt.Errorf("invalid option nameAndArg :%v", nameAndArg))
	}
	arg := tryGroup([]int{2})

	return New(name).Arguments(arg, "", nil), nil
}

func (c *Command) command(nameAndArg, desc string) *Command {
	cmd, err := newCommandWithNameAndArg(nameAndArg)
	if err != nil {
		panic(err)
	}
	if c._subCommands.has(cmd.name) {
		panic("command already exists")
	}

	cmd.Description(desc).parent = c
	c._subCommands[cmd.name] = cmd
	return cmd
}
