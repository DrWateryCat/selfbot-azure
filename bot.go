package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Moonlington/FloSelfbot/commands"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

func logwarning(e error) {
	if e != nil {
		log.Println(e)
	}
}

func logerror(e error) {
	if e != nil {
		log.Println(e)
	}
}

// EditConfigfile edits the config file using the conf passed
func EditConfigfile(conf *commands.Config) {
	f, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logwarning(err)
	defer f.Close()
	logwarning(toml.NewEncoder(f).Encode(conf))
}

var (
	conf           *commands.Config
	commandhandler *commands.CommandHandler
)

func createConfig() *commands.Config {
	var (
		tempprefix string
		temptoken  string
	)

	fmt.Println("\n​​To find your user token (desktop app and browser):\n1. Type Ctrl-Shift-i\n2. Go to the Application page\n3. Under Storage, select Local Storage, and then discordapp.com\n4. Find the token row and copy the value that is in quotes.")
	fmt.Print("\nInput your User Token here: ")
	fmt.Scanln(&temptoken)
	fmt.Print("\nInput your desired prefix here: ")
	fmt.Scanln(&tempprefix)

	tempconfig := &commands.Config{Token: temptoken, Prefix: tempprefix}
	EditConfigfile(tempconfig)

	return tempconfig
}

func main() {
	commands.Mgtoggle = false

	commands.AFKMode = false

	_, err := toml.DecodeFile("config.toml", &conf)
	if os.IsNotExist(err) {
		fmt.Println("No config file found, so let's make one!")
		conf = createConfig()
	}
	EditConfigfile(conf)

	dg, err := discordgo.New(conf.Token)

	logwarning(err)

	_, err = dg.User("@me")

	if err != nil {
		fmt.Println("Something went wrong with logging in, check twice if your token is correct.\nYou can do so by editing/deleting config.toml")
		fmt.Println("Press CTRL-C to exit.")
		<-make(chan struct{})
		return
	}

	dg.AddHandlerOnce(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageEdit)
	dg.AddHandler(messageDelete)
	dg.AddHandler(messageReactionAdd)
	dg.AddHandler(messageReactionRemove)
	dg.AddHandler(guildMemberChunk)

	commandhandler = &commands.CommandHandler{Commands: make(map[string]commands.Command)}

	commandhandler.AddCommand("ping", &commands.Ping{})
	commandhandler.AddCommand("setgame", &commands.SetGame{})
	commandhandler.AddCommand("getgame", &commands.GetGame{})
	commandhandler.AddCommand("me", &commands.Me{})
	commandhandler.AddCommand("embed", &commands.Embed{})
	commandhandler.AddCommand("eval", &commands.Eval{})
	commandhandler.AddCommand("clean", &commands.Clean{})
	commandhandler.AddCommand("quote", &commands.Quote{})
	commandhandler.AddCommand("config", &commands.Configcommand{})
	commandhandler.AddCommand("multigame", &commands.MultiGame{})
	commandhandler.AddCommand("status", &commands.Status{})
	commandhandler.AddCommand("react", &commands.React{})
	commandhandler.AddCommand("ban", &commands.Ban{})
	commandhandler.AddCommand("softban", &commands.Softban{})
	commandhandler.AddCommand("kick", &commands.Kick{})
	commandhandler.AddCommand("poll", &commands.Poll{})
	if conf.SketchyMode {
		commandhandler.AddCommand("afk", &commands.Afk{})
	}

	err = dg.Open()

	logwarning(err)

	if conf.SketchyMode {
		fmt.Println("You have turned on sketchy mode, this enables a few more features of the selfbot.\nBUT be aware that these features can get your account banned! Be careful or turn sketchy mode off.")
	}
	fmt.Println("FloSelfbot is now running.")
	fmt.Println("Type", conf.Prefix+"help", "to see all commands!")

	if conf.MultigameToggled {
		commands.Mgtoggle = true
		go commands.MultiGameFunc(dg, conf)
	}

	go bufferLoop(dg)

	fmt.Println("Press CTRL-C to exit.")
	<-make(chan struct{})
	return
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	for _, g := range r.Guilds {
		s.RequestGuildMembers(g.ID, "", 0)
	}
}

func guildMemberChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
	for _, g := range s.State.Guilds {
		if g.ID == c.GuildID {
			g.Members = append(g.Members, c.Members...)
		}
	}
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if conf.LogMode {
		timestamp := time.Now().UTC()
		logMessageNoAuthor(s, timestamp, m.UserID, m.MessageID, m.ChannelID, "REA", m.Emoji.Name, m.Emoji.APIName())
	}
}

func messageReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if conf.LogMode {
		timestamp := time.Now().UTC()
		logMessageNoAuthor(s, timestamp, m.UserID, m.MessageID, m.ChannelID, "RED", m.Emoji.Name, m.Emoji.APIName())
	}
}

func messageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if conf.LogMode {
		timestamp := time.Now().UTC()
		logMessage(s, timestamp, m.Message.Author, m.ID, m.ChannelID, "EDI", m.ContentWithMentionsReplaced())
	}
}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if conf.LogMode {
		timestamp := time.Now().UTC()
		logMessage(s, timestamp, m.Message.Author, m.ID, m.ChannelID, "DEL", m.ContentWithMentionsReplaced())
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if conf.LogMode {
		timestamp, _ := m.Timestamp.Parse()
		logMessage(s, timestamp, m.Message.Author, m.ID, m.ChannelID, "MSG", m.ContentWithMentionsReplaced())
		if len(m.Attachments) != 0 {
			for _, a := range m.Attachments {
				logMessage(s, timestamp, m.Message.Author, m.ID, m.ChannelID, "ATT", a.URL)
			}
		}
	}

	if commands.AFKMode {
		for _, u := range m.Mentions {
			if u.ID == s.State.User.ID {
				commands.AFKMessages = append(commands.AFKMessages, m)
				emcolor := s.State.UserColor(s.State.User.ID, m.ChannelID)
				em := &discordgo.MessageEmbed{Color: emcolor, Title: fmt.Sprintf("**%s** Is AFK!", s.State.User.Username)}
				if commands.AFKstring != "" {
					em.Description = commands.AFKstring
					s.ChannelMessageSendEmbed(m.ChannelID, em)
				} else {
					s.ChannelMessageSendEmbed(m.ChannelID, em)
				}
			}
		}
	}

	// Ignore all messages created by other users
	if m.Author.ID != s.State.User.ID {
		return
	}

	if len(m.Content) > 0 && strings.HasPrefix(strings.ToLower(m.Content), conf.Prefix) {
		// Setting values for the commands
		var ctx *commands.Context
		args := strings.Fields(m.Content[len(conf.Prefix):])
		invoked := args[0]
		args = args[1:]
		argstr := m.Content[len(conf.Prefix)+len(invoked):]
		if argstr != "" {
			argstr = argstr[1:]
		}
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			channel, _ = s.State.PrivateChannel(m.ChannelID)
			ctx = &commands.Context{Conf: conf, Invoked: invoked, Argstr: argstr, Args: args, Channel: channel, Guild: nil, Mess: m, Sess: s}
		} else {
			guild, _ := s.State.Guild(channel.GuildID)
			ctx = &commands.Context{Conf: conf, Invoked: invoked, Argstr: argstr, Args: args, Channel: channel, Guild: guild, Mess: m, Sess: s}
		}
		p, _ := s.UserChannelPermissions(s.State.User.ID, m.ChannelID)
		if channel.Recipient == nil {
			if p&discordgo.PermissionEmbedLinks != discordgo.PermissionEmbedLinks {
				logerror(errors.New("THE SELFBOT DOES NOT WORK IN CHANNELS WHERE YOU DONT HAVE EMBEDLINKS PERMISSION"))
				return
			}
		}

		commandhandler.HandleCommands(ctx)
	}
}
