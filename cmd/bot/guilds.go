package main

import (
	"fmt"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/cmd/bot/monitoring"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"golang.org/x/exp/slog"
)

func (a *App) guildJoinedHandler() func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(_ *discordgo.Session, g *discordgo.GuildCreate) {
		a.Log().Info(fmt.Sprintf("Joined guild %s", g.Name))

		if err := a.registerSlashCommands(); err != nil {
			a.Log().Error("Error registering slash commands", slog.String(logging.KeyError, err.Error()))
		}

		// Increment the total number of guilds.
		monitoring.TotalDiscordGuilds.Inc()
	}
}

func (a *App) guildLeaveHandler() func(s *discordgo.Session, g *discordgo.GuildDelete) {
	return func(_ *discordgo.Session, g *discordgo.GuildDelete) {
		a.Log().Info(fmt.Sprintf("Left guild %s", g.Name))

		if err := a.unregisterSlashCommands(); err != nil {
			a.Log().Error("Error unregistering slash commands", slog.String(logging.KeyError, err.Error()))
		}

		// Decrement the total number of guilds.
		monitoring.TotalDiscordGuilds.Dec()
	}
}
