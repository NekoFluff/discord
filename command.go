package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command discordgo.ApplicationCommand
	Handler interface{} // Matching the session AddHandler function, the fist argument is a *discordgo.Session object and the second is a *discordgo.InteractionCreate object
}
