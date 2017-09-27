package commands

type Boolin struct{}

func (b *Boolin) message(ctx *Context) {
	if len(ctx.Args) < 2 {
		ctx.QuickSendEm(":b:ool :b:achine :b:roke: not enough arguments")
		return
	}

	cmd := ctx.Args[0]
	if cmd == "all" {
		words := ctx.Args[1:]

		amt := 0

		for _, w := range words {
			amt += len(w)
		}

		retstr := ""
		for i := 0; i < amt; i++ {
			retstr += ":b:"
		}

		ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.Message.ID)
		_, err := ctx.QuickSendEm(retstr)
		if err != nil {
			ctx.QuickSendEm(":b:ool :b:achine :b:roke, error: " + err.Error())
		}
		return
	} else if cmd == "other" {
		words := ctx.Args[1:]

		retstr := ""
		for _, word := range words {
			retstr += ":b:" + word[1:] + " "
		}

		ctx.Sess.ChannelMessageDelete(ctx.Mess.ChannelID, ctx.Mess.Message.ID)
		ctx.QuickSendEm(retstr)
		return
	}

}

func (b *Boolin) description() string             { return ":b:ools your :b:essage up ni:b::b:a" }
func (b *Boolin) usage() string                   { return "[all or other] <message>" }
func (b *Boolin) detailed() string                { return b.description() }
func (b *Boolin) subcommands() map[string]Command { return make(map[string]Command) }
