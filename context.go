package commandergo

type Context struct {
	cmd        *Command
	parsedArgs map[string]Varaint
	parsedOpts map[string]Varaint
}

func newContext(cmd *Command) *Context {
	return &Context{
		cmd:        cmd,
		parsedArgs: make(map[string]Varaint),
		parsedOpts: make(map[string]Varaint),
	}
}

func (ctx *Context) Arg(name string) Varaint { return ctx.parsedArgs[name] }
func (ctx *Context) Opt(name string) Varaint { return ctx.parsedOpts[name] }
func (ctx *Context) Command() *Command       { return ctx.cmd }
