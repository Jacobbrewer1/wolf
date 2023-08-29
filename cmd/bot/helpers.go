package main

import (
	"math/rand"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/pkg/messages"
)

func respondSlashError(a IApp, i *discordgo.InteractionCreate) error {
	return respondSlashEphemeral(a, i, messages.ErrUserErrorProcessing)
}

func respondSlashEphemeral(a IApp, i *discordgo.InteractionCreate, content string) error {
	return a.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
