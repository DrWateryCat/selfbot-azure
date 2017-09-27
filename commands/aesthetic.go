package commands

import (
	"strconv"
)

type Aesthetic struct{}

func (a *Aesthetic) message(ctx *Context) {
	if len(ctx.Args) < 2 {
		ctx.QuickSendEm("Not enough arguments!")
		return
	}

	spaceAmount, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		ctx.QuickSendEm(err.Error())
	}
	words := ctx.Args[1:]

	retstr := ""

	for _, word := range words {
		for i := 0; i < len(word); i++ {
			retstr += string(word[i])
			for j := 0; j < spaceAmount; j++ {
				retstr += " "
			}
		}
	}

	ctx.QuickSendEm(retstr)
	return
}

func (a *Aesthetic) description() string             { return "Makes your text more A E S T H E T I C " }
func (a *Aesthetic) usage() string                   { return "[spaces] <words you want to be A E S T H E T I C>" }
func (a *Aesthetic) detailed() string                { return a.description() }
func (a *Aesthetic) subcommands() map[string]Command { return make(map[string]Command) }
