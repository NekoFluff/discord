package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestSession_ConcreteImplementationMatchesInterface(t *testing.T) {
	var _ Session = &discordgo.Session{}
}
