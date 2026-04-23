package commandergo

type Command struct {
	name         string
	version      string
	description  string
	parent       *Command
	_arguments   Arguments
	_options     Options
	_subCommands Commands
	actionFn     func(ctx *Context)
}

func New(name string) *Command {
	return &Command{
		name:         name,
		parent:       nil,
		_arguments:   Arguments{},
		_options:     Options{},
		_subCommands: Commands{},
	}
}

func (c *Command) Name(name string) *Command {
	c.name = name
	return c
}

func (c *Command) Version(version string) *Command {
	c.version = version
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

func (c *Command) Parent() *Command {
	return c.parent
}

func (c *Command) Action(call func(ctx *Context)) *Command {
	c.actionFn = call
	return c
}

func (c *Command) Parse() error {
	return nil
}
