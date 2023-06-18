package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, m *discordgo.InteractionCreate)
}
