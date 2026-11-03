package commandergo

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Commands []*Command

func (o Commands) has(name string) bool {
	for _, cmd := range o {
		if cmd.name == name {
			return true
		}
	}
	return false
}

func (o Commands) get(name string) (*Command, bool) {
	for _, cmd := range o {
		if cmd.name == name {
			return cmd, true
		}
	}
	return nil, false
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

	cmd := New(name)
	if arg != "" {
		cmd.Arguments(arg, "", nil)
	}
	return cmd, nil
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
	c._subCommands = append(c._subCommands, cmd)
	return cmd
}

type JsonArgument struct {
	Name string  `json:"name"`
	Desc *string `json:"desc"`

	DefaultValue  any  `json:"defaultValue"`
	ValueRequired bool `json:"valueRequired"`
	MultiValue    bool `json:"multiValue"`
}
type JsonOption struct {
	Name  string  `json:"name"`
	Alias *string `json:"alias"`
	Desc  *string `json:"desc"`

	Value *struct {
		ValueName     string `json:"valueName"`
		DefaultValue  any    `json:"defaultValue"`
		ValueRequired bool   `json:"valueRequired"`
	} `json:"value"`
}
type JsonCommand struct {
	Name        string          `json:"name"`
	Version     *string         `json:"version"`
	Description *string         `json:"description"`
	Arguments   *[]JsonArgument `json:"arguments"`
	Options     *[]JsonOption   `json:"options"`
	Commands    *[]JsonCommand  `json:"commands"`
}

func newFromJson(jsonData string) *Command {
	jsonCmd := &JsonCommand{}
	decoder := json.NewDecoder(strings.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(jsonCmd); err != nil {
		panic(err)
	}
	var initArgument func(cmd *Command, jsonArg *JsonArgument)
	var initOption func(cmd *Command, jsonOpt *JsonOption)
	var initCmd func(cmd *Command, jsonCmd *JsonCommand)
	initArgument = func(cmd *Command, jsonArg *JsonArgument) {
		name, desc := func() (string, string) {
			var n, d string = "", ""
			if jsonArg.ValueRequired {
				n = fmt.Sprintf("<%s%v>", jsonArg.Name, IfEl(jsonArg.MultiValue, "...", ""))
			} else {
				n = fmt.Sprintf("[%s%v]", jsonArg.Name, IfEl(jsonArg.MultiValue, "...", ""))
			}

			if jsonArg.Desc != nil {
				d = *jsonArg.Desc
			}

			return n, d
		}()
		cmd.Arguments(name, desc, normalizeJSONValue(jsonArg.DefaultValue))
	}
	initOption = func(cmd *Command, jsonOpt *JsonOption) {
		flag, desc, defaultValue := func() (string, string, any) {
			var f, d string = "", ""
			var v any = nil

			if jsonOpt.Alias != nil {
				f = fmt.Sprintf("-%s,", *jsonOpt.Alias)
			}
			f = fmt.Sprintf("%s--%s", f, jsonOpt.Name)
			if jsonOpt.Value != nil {
				if jsonOpt.Value.ValueRequired {
					f = fmt.Sprintf("%s <%s>", f, jsonOpt.Value.ValueName)
				} else {
					f = fmt.Sprintf("%s [%s]", f, jsonOpt.Value.ValueName)
				}

				v = normalizeJSONValue(jsonOpt.Value.DefaultValue)
			}

			if jsonOpt.Desc != nil {
				d = *jsonOpt.Desc
			}

			return f, d, v
		}()
		cmd.Options(flag, desc, defaultValue)
	}
	initCmd = func(cmd *Command, jsonCmd *JsonCommand) {
		if jsonCmd.Version != nil {
			cmd.Version(*jsonCmd.Version)
		}
		if jsonCmd.Description != nil {
			cmd.Description(*jsonCmd.Description)
		}
		if jsonCmd.Arguments != nil {
			for _, arg := range *jsonCmd.Arguments {
				initArgument(cmd, &arg)
			}
		}
		if jsonCmd.Options != nil {
			for _, opt := range *jsonCmd.Options {
				initOption(cmd, &opt)
			}
		}
		if jsonCmd.Commands != nil {
			for _, child := range *jsonCmd.Commands {
				c := cmd.Command(child.Name, "")
				initCmd(c, &child)
			}
		}
	}
	c := New(jsonCmd.Name)
	initCmd(c, jsonCmd)
	return c
}

func normalizeJSONValue(v any) any {
	switch value := v.(type) {
	case json.Number:
		if i, err := value.Int64(); err == nil {
			return int(i)
		}
		if f, err := value.Float64(); err == nil {
			return f
		}
	case []any:
		res := make([]any, len(value))
		for i, item := range value {
			res[i] = normalizeJSONValue(item)
		}
		return res
	case map[string]any:
		res := make(map[string]any, len(value))
		for k, item := range value {
			res[k] = normalizeJSONValue(item)
		}
		return res
	}

	return v
}
