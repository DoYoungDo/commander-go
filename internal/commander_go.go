package commandergo

type Action func(ctx *Context)
type ActionE func(ctx *Context) error

type Command struct {
	name         string
	version      string
	description  string
	parent       *Command
	_arguments   Arguments
	_options     Options
	_subCommands Commands
	actionFn     ActionE
	customHelpFn ActionE
}

func New(name string) *Command {
	c := &Command{
		name:         name,
		parent:       nil,
		_arguments:   Arguments{},
		_options:     Options{},
		_subCommands: Commands{},
		actionFn:     nil,
		customHelpFn: nil,
	}
	c.options("-h, --help", "display help for command", nil)
	return c
}

func NewFromJson(jsonData string) *Command {
	return newFromJson(jsonData)
}

func (c *Command) Name(name string) *Command {
	c.name = name
	return c
}

func (c *Command) CustomHelp(hp ActionE) *Command {
	c.customHelpFn = hp
	return c
}

func (c *Command) Version(version string) *Command {
	c.version = version
	c.options("-V, --version", "output the version number", nil)
	return c
}

func (c *Command) Description(description string) *Command {
	c.description = description
	return c
}

func (c *Command) Arguments(name, desc string, defaultValue any) *Command {
	return c.arguments(name, desc, defaultValue)
}

func (c *Command) Options(flag, desc string, defaultValue any) *Command {
	return c.options(flag, desc, defaultValue)
}

func (c *Command) Command(nameAndArg, desc string) *Command {
	return c.command(nameAndArg, desc)
}

func (c *Command) Child(name string) *Command {
	cmd, ok := c._subCommands.get(name)
	if ok {
		return cmd
	}
	return nil
}

func (c *Command) Parent() *Command {
	return c.parent
}

func (c *Command) Action(call Action) *Command {
	if call == nil {
		// 允许 Action(nil) 清空之前注册的 action。
		c.actionFn = nil
		return c
	}
	c.actionFn = func(ctx *Context) error {
		call(ctx)
		return nil
	}
	return c
}

func (c *Command) ActionE(call ActionE) *Command {
	c.actionFn = call
	return c
}

func (c *Command) Parse(args []string) error {
	return c.parse(args, &parseOption{strict: false})
}
func (c *Command) ParseStrict(args []string) error {
	return c.parse(args, &parseOption{strict: true})
}
