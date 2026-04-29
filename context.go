package commandergo

type Context struct {
	cmd        *Command
	parsedArgs []Varaint
	parsedOpts map[string]Varaint
}

func newContext(cmd *Command) *Context {
	return &Context{
		cmd:        cmd,
		parsedArgs: []Varaint{},
		parsedOpts: make(map[string]Varaint),
	}
}

func (ctx *Context) Args() []Varaint         { return ctx.parsedArgs }
func (ctx *Context) Opt(name string) Varaint { return ctx.parsedOpts[name] }
func (ctx *Context) Command() *Command       { return ctx.cmd }
