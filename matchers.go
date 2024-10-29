package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ContainsInteractionResponseMatcher struct {
	msg string
}

func ContainsInteractionResponse(msg string) ContainsInteractionResponseMatcher {
	return ContainsInteractionResponseMatcher{msg: msg}
}

func (m ContainsInteractionResponseMatcher) Matches(input interface{}) bool {
	response, ok := input.(*discordgo.InteractionResponse)
	if !ok {
		return false
	}

	return strings.Contains(response.Data.Content, m.msg)
}

func (m ContainsInteractionResponseMatcher) String() string {
	return fmt.Sprintf("to contain msg '%v'", m.msg)
}

func (m ContainsInteractionResponseMatcher) Got(input interface{}) string {
	response, ok := input.(*discordgo.InteractionResponse)
	if !ok {
		return "not an interaction response"
	}

	return response.Data.Content
}
