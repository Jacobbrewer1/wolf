package main

import (
	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/pkg/messages"
)

func respondEphemeralError(a IApp, i *discordgo.InteractionCreate) error {
	return respondEphemeral(a, i, messages.ErrUserErrorProcessing)
}

func respondEphemeral(a IApp, i *discordgo.InteractionCreate, content string) error {
	return a.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func hasRole(member *discordgo.Member, roleID string) bool {
	for _, role := range member.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}
