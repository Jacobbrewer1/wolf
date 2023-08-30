package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/pkg/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// setupCmd is the command for all configuration commands.
	setupCmdName = "setup"

	// enableTicketingCmdName is the command for all ticketing configuration commands.
	enableTicketingCmdName = "ticketing_enable"

	// disableTicketingCmdName is the command for all ticketing configuration commands.
	disableTicketingCmdName = "ticketing_disable"

	// channelCmdName is the text for the channel command.
	channelCmdName = "channel"

	// roleCmdName is the text for the role command.
	roleCmdName = "role"
)

var (
	// setupCmd is the command for all configuration commands.
	setupCmd = &discordgo.ApplicationCommand{
		Name:        setupCmdName,
		Type:        discordgo.ChatApplicationCommand,
		Description: "This is the command for all configuration commands.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        enableTicketingCmdName,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This will in the channel you specify.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        channelCmdName,
						Type:        discordgo.ApplicationCommandOptionChannel,
						Description: "This is the channel you want to enable ticketing in.",
						Required:    true,
					},
					{
						Name:        roleCmdName,
						Type:        discordgo.ApplicationCommandOptionRole,
						Description: "This is the role you want to handle tickets.",
						Required:    true,
					},
				},
			},
			{
				Name:        disableTicketingCmdName,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This will disable ticketing for your server.",
			},
		},
	}
)

func setupCmdController(a IApp, i *discordgo.InteractionCreate) (commandProcessor, error) {
	// Ensure the user is an administrator.
	if i.Member.Permissions&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
		err := respondEphemeral(a, i, "You must be an administrator to use this command")
		if err != nil {
			return nil, nil
		}
		return nil, nil
	}

	// Extract the sub command.
	subCmd := i.ApplicationCommandData().Options[0].Name

	switch subCmd {
	case enableTicketingCmdName:
		return enableTicketingCmdController, nil
	case disableTicketingCmdName:
		return disableTicketingCmdController, nil
	default:
		return nil, fmt.Errorf("unhandled sub command %s", subCmd)
	}
}

// enableTicketingCmdController is the controller for the enable ticketing command.
func enableTicketingCmdController(a IApp, i *discordgo.InteractionCreate) error {
	// Extract the channel provided.
	channel := i.ApplicationCommandData().Options[0].Options[0].ChannelValue(a.Session())

	// Extract the role provided.
	role := i.ApplicationCommandData().Options[0].Options[1].RoleValue(a.Session(), i.GuildID)

	// Ensure the channel is a text channel.
	if channel.Type != discordgo.ChannelTypeGuildText {
		return respondEphemeral(a, i, "You must provide a text channel for ticketing.")
	}

	gd := a.GuildDal(context.Background())

	// Get the guild.
	guild, err := gd.GetGuildByID(i.GuildID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("error getting guild: %w", err)
	}

	if guild == nil {
		guild = &entities.Guild{
			ID: i.GuildID,
		}
	}

	// Enable ticketing for the guild.
	guild.Ticketing.Enabled = true

	// Set the ticketing channel.
	guild.Ticketing.ChannelID = channel.ID

	// Set the ticketing role.
	guild.Ticketing.RoleID = role.ID

	msg := new(discordgo.Message)

	// Check to see if the ticketing message still exists.
	if guild.Ticketing.OpenMessageID != "" {
		// Get the ticketing message.
		msg, err = a.Session().ChannelMessage(channel.ID, guild.Ticketing.OpenMessageID)
		// If the message does not exist, set the message ID to an empty string.
		if err != nil {
			var restErr *discordgo.RESTError
			ok := errors.As(err, &restErr)
			if ok && restErr.Message.Code == discordgo.ErrCodeUnknownMessage {
				guild.Ticketing.OpenMessageID = ""
			} else {
				return fmt.Errorf("error getting ticketing message: %w", err)
			}
		}

		// If the message does not exist, set the message ID to an empty string.
		if msg == nil {
			guild.Ticketing.OpenMessageID = ""
		}
	}

	// If the ticketing message ID is empty, send a new message.
	if guild.Ticketing.OpenMessageID == "" {
		// Send the ticketing message to the channel.
		msg, err = sendOpenTicketMessage(a, channel)
		if err != nil {
			return fmt.Errorf("error sending open ticket message: %w", err)
		}
	}

	// Set the ticketing message ID.
	guild.Ticketing.OpenMessageID = msg.ID

	// Save the guild.
	if err := gd.SaveGuild(guild); err != nil {
		return fmt.Errorf("error saving guild: %w", err)
	}

	// Respond to the interaction saying that ticketing has been enabled in channel <channel>.
	if err := respondEphemeral(a, i, fmt.Sprintf("Ticketing has been enabled in channel <#%s>", channel.ID)); err != nil {
		return fmt.Errorf("error responding to interaction: %w", err)
	}

	return nil
}

// disableTicketingCmdController is the controller for the disable ticketing command.
func disableTicketingCmdController(a IApp, i *discordgo.InteractionCreate) error {
	gd := a.GuildDal(context.Background())

	// Get the guild.
	guild, err := gd.GetGuildByID(i.GuildID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("error getting guild: %w", err)
	}

	if guild == nil {
		guild = &entities.Guild{
			ID: i.GuildID,
		}
	}

	// Disable ticketing for the guild.
	guild.Ticketing.Enabled = false

	// Save the guild.
	if err := gd.SaveGuild(guild); err != nil {
		return fmt.Errorf("error saving guild: %w", err)
	}

	// Respond to the interaction saying that ticketing has been disabled.
	if err := respondEphemeral(a, i, "Ticketing has been disabled"); err != nil {
		return fmt.Errorf("error responding to interaction: %w", err)
	}

	return nil
}
