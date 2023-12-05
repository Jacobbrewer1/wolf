package main

import (
	"fmt"
	"log/slog"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
)

func (a *App) guildJoinedHandler() func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(_ *discordgo.Session, g *discordgo.GuildCreate) {
		slog.Info(fmt.Sprintf("Joined guild %s", g.Name))

		if err := a.registerSlashCommands(); err != nil {
			slog.Error("Error registering slash commands", slog.String(logging.KeyError, err.Error()))
		}

		// Increment the total number of guilds.
		TotalDiscordGuilds.Inc()
	}
}

func (a *App) guildLeaveHandler() func(s *discordgo.Session, g *discordgo.GuildDelete) {
	return func(_ *discordgo.Session, g *discordgo.GuildDelete) {
		slog.Info(fmt.Sprintf("Left guild %s", g.Name))

		if err := a.unregisterSlashCommands(); err != nil {
			slog.Error("Error unregistering slash commands", slog.String(logging.KeyError, err.Error()))
		}

		// Decrement the total number of guilds.
		TotalDiscordGuilds.Dec()
	}
}
