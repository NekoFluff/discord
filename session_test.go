package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestSession_ConcreateImplementationMatchesInterface(t *testing.T) {
	var _ Session = &discordgo.Session{}
}
