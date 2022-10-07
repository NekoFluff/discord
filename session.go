//go:generate mockgen -package=discord -source=session.go -destination=session_mock.go
package discord

import "github.com/bwmarrin/discordgo"

type Session interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
	ChannelMessageSendEmbed(string, *discordgo.MessageEmbed) (*discordgo.Message, error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) (err error)
}
