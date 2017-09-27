package commands

import "github.com/getwe/figlet4go"

type Ascii struct{}

func (a *Ascii) message(ctx *Context) {
	if len(ctx.Args) == 0 {
		ctx.QuickSendEm("Not enough args!")
		return
	}

	wordstr := ""

	for _, word := range ctx.Args {
		wordstr += word
		wordstr += " "
	}

	as := figlet4go.NewAsciiRender()
	renderstr, _ := as.Render(wordstr)

	ctx.QuickSendEm("```\n" + renderstr + "```")
	return
}

func (a *Ascii) description() string             { return "Makes your text more ASCII" }
func (a *Ascii) usage() string                   { return "<words you want to be ASCII>" }
func (a *Ascii) detailed() string                { return a.description() }
func (a *Ascii) subcommands() map[string]Command { return make(map[string]Command) }
