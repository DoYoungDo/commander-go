package commandergo

type Commands map[string]*Command

func (o Commands) has(name string) bool {
	_, ok := o[name]
	return ok
}

func newCommandWithNameAndArg(nameAndArg string) (*Command, error) {
	// return New(nameAndArg)
	return nil, nil
}

func (c *Command) command(nameAndArg, desc string) *Command {
	cmd, err := newCommandWithNameAndArg(nameAndArg)
	if err != nil {
		panic(err)
	}
	cmd.Description(desc).parent = c

	c._subCommands[cmd.name] = cmd
	return cmd
}
