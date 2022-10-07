package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session      *discordgo.Session
	Commands     map[string]Command
	DeveloperIDs []string
}

func NewBot(token string) *Bot {
	session, err := createSession(token)
	if err != nil {
		log.Fatal(err)
	}

	bot := &Bot{
		Session:      session,
		Commands:     make(map[string]Command),
		DeveloperIDs: []string{},
	}

	bot.Session.AddHandler(bot.handleInteractionCreate)
	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	return bot
}

func (bot *Bot) Stop() {
	// Cleanly close down the Discord session.
	bot.Session.Close()
}

func createSession(Token string) (s *discordgo.Session, err error) {
	// Create a new Discord session using the provided bot token.
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	return
}

func (bot *Bot) AddCommands(cmds ...Command) {
	for _, cmd := range cmds {
		bot.Commands[cmd.Command.Name] = cmd
	}
}

func (bot *Bot) RegisterCommands() {
	for _, cmd := range bot.Commands {
		cmd, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, "", &cmd.Command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", cmd.Name, err)
		}
	}
}

// This function will be called every time a new
// message is created on any channel that the authenticated bot has access to.
func (bot *Bot) handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if cmd, ok := bot.Commands[i.ApplicationCommandData().Name]; ok {
		cmd.Handler(s, i)
	}
}

func (bot *Bot) SendChannelMessage(channelName string, message string) {
	for _, guild := range bot.Session.State.Guilds {
		// Get channels for this guild (a.k.a discord server)
		channels, _ := bot.Session.GuildChannels(guild.ID)

		for _, c := range channels {
			// Ensure the channel is a guild text channel and not a voice or DM channel
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			// Check if the channel name matches target name
			if c.Name != channelName {
				continue
			}

			// Send a message to the discord channel
			_, err := bot.Session.ChannelMessageSend(
				c.ID,
				message,
			)
			if err != nil {
				log.Println("An error occurred while sending a message to a discord server")
				log.Println(err)
			}
		}
	}
}

func (bot *Bot) SendDeveloperMessage(message string) {
	for _, developerId := range bot.DeveloperIDs {
		ch, err := bot.Session.UserChannelCreate(developerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = bot.Session.ChannelMessageSend(ch.ID, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (bot *Bot) SendEmbedMessage(channelName string, message *discordgo.MessageEmbed) {
	for _, guild := range bot.Session.State.Guilds {
		// Get channels for this guild (a.k.a discord server)
		channels, _ := bot.Session.GuildChannels(guild.ID)

		for _, c := range channels {
			// Ensure the channel is a guild text channel and not a voice or DM channel
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			// Check if the channel name matches target name
			if c.Name != channelName {
				continue
			}

			// Send a message to the discord channel
			_, err := bot.Session.ChannelMessageSendEmbed(
				c.ID,
				message,
			)
			if err != nil {
				log.Println("An error occurred while sending a message to a discord server")
				log.Println(err)
			}
		}
	}
}
