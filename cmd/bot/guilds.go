package main

import (
	"fmt"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/cmd/bot/monitoring"
)

func guildJoinedHandler(a IApp) func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(_ *discordgo.Session, g *discordgo.GuildCreate) {
		a.Log().Info(fmt.Sprintf("Joined guild %s", g.Name))

		// Increment the total number of guilds.
		monitoring.TotalDiscordGuilds.Inc()
	}
}

func guildLeaveHandler(a IApp) func(s *discordgo.Session, g *discordgo.GuildDelete) {
	return func(_ *discordgo.Session, g *discordgo.GuildDelete) {
		a.Log().Info(fmt.Sprintf("Left guild %s", g.Name))

		// Decrement the total number of guilds.
		monitoring.TotalDiscordGuilds.Dec()
	}
}
