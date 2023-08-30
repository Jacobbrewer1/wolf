package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jacobbrewer1/discordgo"
	"github.com/Jacobbrewer1/wolf/pkg/custom"
	"github.com/Jacobbrewer1/wolf/pkg/entities"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

const (
	// OpenTicketButtonID is the ID for the open ticket button.
	OpenTicketButtonID = "open_ticket_button"

	// ClaimTicketButtonID is the ID for the claim ticket button.
	ClaimTicketButtonID = "claim_ticket_button"

	// CloseTicketButtonID is the ID for the close ticket button.
	CloseTicketButtonID = "close_ticket_button"

	// ReopenTicketButtonID is the ID for the reopen ticket button.
	ReopenTicketButtonID = "reopen_ticket_button"

	// DeleteTicketButtonID is the ID for the delete ticket button.
	DeleteTicketButtonID = "delete_ticket_button"
)

const (
	// ClaimEmoji is the emoji that will be used for the claim button. (Ticket)
	ClaimEmoji = "\U0001F3AB"

	// CloseEmoji is the emoji that will be used for the claim button. (Padlock)
	CloseEmoji = "\U0001F510"

	// ReopenEmoji is the emoji that will be used for the claim button. (Open padlock)
	ReopenEmoji = "\U0001F513"

	// DeleteEmoji is the emoji that will be used for the claim button. (Cross)
	DeleteEmoji = "\u274C"
)

const (
	// TicketCmdName is the command for claiming a ticket.
	TicketCmdName = "ticket"

	// ClaimCmdName is the sub command for claiming a ticket.
	ClaimCmdName = "claim"

	// CloseCmdName is the sub command for closing the verification process.
	CloseCmdName = "close"

	// DeleteCmdName is the sub command for deleting the verification process.
	DeleteCmdName = "delete"

	// ReopenCmdName is the sub command for reopening the verification process.
	ReopenCmdName = "reopen"
)

var (
	// ticketCmd is the command for controlling tickets.
	ticketCmd = &discordgo.ApplicationCommand{
		Name:        TicketCmdName,
		Type:        discordgo.ChatApplicationCommand,
		Description: "This is the command for controlling tickets.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        ClaimCmdName,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This claims the ticket for the channel that the command was executed in.",
			},
			{
				Name:        CloseCmdName,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This closes the ticket for the channel that the command was executed in.",
			},
			{
				Name:        DeleteCmdName,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This deletes the ticket for the channel that the command was executed in.",
			},
			{
				Name:        ReopenCmdName,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "This reopens the ticket for the channel that the command was executed in.",
			},
		},
	}

	// NewTicketMessage is the message that is sent when a new ticket is created.
	NewTicketMessage = &discordgo.MessageSend{
		Content: `Your ticket has been created.
Please provide any additional info you deem relevant to help us answer faster.`,
		Embed:           nil,
		TTS:             false,
		Files:           nil,
		AllowedMentions: nil,
		Flags:           0,
		// Add the four buttons to the message.
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    fmt.Sprintf("%s Claim", ClaimEmoji),
						Style:    discordgo.PrimaryButton,
						Disabled: false,
						Emoji:    discordgo.ComponentEmoji{},
						URL:      "",
						CustomID: ClaimTicketButtonID,
					},
					discordgo.Button{
						Label:    fmt.Sprintf("%s Close", CloseEmoji),
						Style:    discordgo.SecondaryButton,
						Disabled: false,
						Emoji:    discordgo.ComponentEmoji{},
						URL:      "",
						CustomID: CloseTicketButtonID,
					},
					discordgo.Button{
						Label:    fmt.Sprintf("%s Reopen", ReopenEmoji),
						Style:    discordgo.SuccessButton,
						Disabled: true,
						Emoji:    discordgo.ComponentEmoji{},
						URL:      "",
						CustomID: ReopenTicketButtonID,
					},
					discordgo.Button{
						Label:    fmt.Sprintf("%s Delete", DeleteEmoji),
						Style:    discordgo.DangerButton,
						Disabled: false,
						Emoji:    discordgo.ComponentEmoji{},
						URL:      "",
						CustomID: DeleteTicketButtonID,
					},
				},
			},
		},
	}
)

func sendOpenTicketMessage(a IApp, channel *discordgo.Channel) (*discordgo.Message, error) {
	const messageText = `How can we help?
Welcome to our tickets channel. If you have any questions or inquiries, please click on the button below to contact the staff by opening a ticket!`

	// The ticket emoji is the emoji that will be used for the button. (Envelope with arrow)
	const ticketEmoji = "\U0001F4E9"

	// Create the button with the ticket emoji.
	button := discordgo.Button{
		Label:    fmt.Sprintf("%s Open Ticket", ticketEmoji),
		Style:    discordgo.PrimaryButton,
		Disabled: false,
		Emoji:    discordgo.ComponentEmoji{},
		URL:      "",
		CustomID: OpenTicketButtonID,
	}

	// Create the message.
	message := discordgo.MessageSend{
		Content:         messageText,
		Embed:           nil,
		TTS:             false,
		Files:           nil,
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Flags:           0,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					button,
				},
			},
		},
	}

	// Send the message.
	msg, err := a.Session().ChannelMessageSendComplex(channel.ID, &message)
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return msg, nil
}

// createTicket is the function for creating a ticket.
func createTicket(a IApp, i *discordgo.InteractionCreate) error {
	// Get the guild configuration.
	guild, err := a.GuildDal(context.Background()).GetGuildByID(i.GuildID)
	if err != nil {
		return fmt.Errorf("error getting guild configuration: %w", err)
	}

	category := new(discordgo.Channel)

	// Ensure that the category exists for created tickets.
	category, err = a.Session().Channel(guild.Ticketing.CreatedTicketsCategoryID)
	if err != nil {
		er := new(discordgo.RESTError)
		if errors.As(err, &er) && (er.Message.Code == discordgo.ErrCodeUnknownChannel || er.Message.Code == discordgo.ErrCodeGeneralError) { // General is thrown when a 404 is returned.
			a.Log().Warn("Created tickets category does not exist, creating it now")

			category, err = a.Session().GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
				Name: "Created Tickets",
				Type: discordgo.ChannelTypeGuildCategory,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					// Deny @everyone from seeing the ticket.
					{
						ID:    i.GuildID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: 0,
						Deny:  discordgo.PermissionAll,
					},
					// The creator of the ticket can see the ticket.
					{
						ID:    i.Member.User.ID,
						Type:  discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionAllText,
						Deny:  discordgo.PermissionMentionEveryone,
					},
					// Add the ticket role.
					{
						ID:    guild.Ticketing.RoleID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionAllText,
						Deny:  discordgo.PermissionMentionEveryone,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("error creating category: %w", err)
			}

			// Save the guild configuration.
			guild.Ticketing.CreatedTicketsCategoryID = category.ID
			if err := a.GuildDal(context.Background()).SaveGuild(guild); err != nil {
				return fmt.Errorf("error saving guild configuration: %w", err)
			}
		} else {
			return fmt.Errorf("error getting category: %w", err)
		}
	} else if category != nil && category.ID != guild.Ticketing.CreatedTicketsCategoryID {
		// Update the guild configuration.
		guild.Ticketing.CreatedTicketsCategoryID = category.ID
		if err := a.GuildDal(context.Background()).SaveGuild(guild); err != nil {
			return fmt.Errorf("error saving guild configuration: %w", err)
		}
	}

	// Get the latest ticket.
	latestTicket, err := a.TicketDal(context.Background()).GetLatestTicket(i.GuildID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("error getting latest ticket: %w", err)
	}

	ticketID := 1
	// Get the ticket ID.
	if latestTicket != nil {
		ticketID = latestTicket.ID + 1
	}

	// Create the ticket.
	ticket := &entities.Ticket{
		ID:        ticketID,
		GuildID:   i.GuildID,
		UserID:    i.Member.User.ID,
		Username:  i.Member.User.Username,
		CreatedAt: custom.Datetime(time.Now().UTC()),
	}

	// Create the ticket channel only the ticket role and the creator can see.
	ticketChannel, err := a.Session().GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:  ticket.Name(),
		Type:  discordgo.ChannelTypeGuildText,
		Topic: fmt.Sprintf("Ticket created by %s", i.Member.User.Username),
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			// Deny @everyone from seeing the ticket.
			{
				ID:    i.GuildID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: 0,
				Deny:  discordgo.PermissionAll,
			},
			// The creator of the ticket can see the ticket.
			{
				ID:    i.Member.User.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionAllText,
				Deny:  discordgo.PermissionMentionEveryone,
			},
			// Add the ticket role.
			{
				ID:    guild.Ticketing.RoleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionAllText,
				Deny:  discordgo.PermissionMentionEveryone,
			},
		},
		ParentID:         category.ID,
		NSFW:             false,
		Position:         0,
		Bitrate:          0,
		UserLimit:        0,
		RateLimitPerUser: 0,
	})
	if err != nil {
		return err
	}

	// Set the ticket channel ID.
	ticket.ChannelID = ticketChannel.ID

	// Save the ticket.
	if err := a.TicketDal(context.Background()).SaveTicket(ticket); err != nil {
		return fmt.Errorf("error saving ticket: %w", err)
	}

	go func() {
		err := setupNewTicketChannel(a, ticket)
		if err != nil {
			a.Log().Error("Error setting up new ticket channel", slog.String(logging.KeyError, err.Error()))
		}
	}()

	// Respond to the interaction saying that the ticket has been created in channel <channel>.
	// This message is an embedded ephemeral message with all the information about the ticket.
	err = a.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Ticket Created",
					Description: fmt.Sprintf("<@%s>, you created a ticket and it has been moved to the **Created Tickets** category.", i.Member.User.ID),
					Color:       0x00ff00,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Ticket Name",
							Value:  ticket.Name(),
							Inline: true,
						},
						{
							Name:   "Ticket Channel",
							Value:  fmt.Sprintf("<#%s>", ticket.ChannelID),
							Inline: true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error responding to interaction: %w", err)
	}
	return nil
}

func setupNewTicketChannel(a IApp, ticket *entities.Ticket) error {
	// Get the channel.
	channel, err := a.Session().Channel(ticket.ChannelID)
	if err != nil {
		return fmt.Errorf("error getting channel: %w", err)
	}

	// Send the initial message to the channel.
	msg, err := a.Session().ChannelMessageSendComplex(channel.ID, NewTicketMessage)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	// Pin the message to the channel.
	if err := a.Session().ChannelMessagePin(channel.ID, msg.ID); err != nil {
		return fmt.Errorf("error pinning message: %w", err)
	}

	// Update the ticket with the message ID.
	ticket.SetupMessageID = msg.ID

	// Save the ticket.
	if err := a.TicketDal(context.Background()).SaveTicket(ticket); err != nil {
		return fmt.Errorf("error saving ticket: %w", err)
	}

	return nil
}

func claimTicketHandler(a IApp, i *discordgo.InteractionCreate) error {
	// Get the channel name.
	channel, err := a.Session().Channel(i.ChannelID)

	// Get the ticket.
	ticket, err := a.TicketDal(context.Background()).GetTicket(i.GuildID, channel.ID)
	if err != nil {
		return fmt.Errorf("error getting ticket: %w", err)
	}

	// Get the guild configuration.
	guild, err := a.GuildDal(context.Background()).GetGuildByID(i.GuildID)
	if err != nil {
		return fmt.Errorf("error getting guild configuration: %w", err)
	}

	// Get the member that executed the command.
	member, err := a.Session().GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		return fmt.Errorf("error getting member: %w", err)
	}

	// Ensure that the user has the ticket role.
	if !hasRole(member, guild.Ticketing.RoleID) {
		err = respondEphemeral(a, i, "You do not have the ticket role to claim tickets. [<@&"+guild.Ticketing.RoleID+">]")
		if err != nil {
			return fmt.Errorf("error responding to interaction: %w", err)
		}
		return nil
	}

	// Ensure that the ticket is not already claimed.
	if ticket.ClaimedBy != "" && ticket.ClaimedBy != i.Member.User.ID {
		err = respondEphemeral(a, i, "This ticket is already claimed by <@"+ticket.ClaimedBy+">.")
		if err != nil {
			return fmt.Errorf("error responding to interaction: %w", err)
		}
		return nil
	} else if ticket.ClaimedBy == i.Member.User.ID {
		err = respondEphemeral(a, i, "You have already claimed this ticket <@"+ticket.ClaimedBy+">")
		if err != nil {
			return fmt.Errorf("error responding to interaction: %w", err)
		}
		return nil
	}

	// Claim the ticket.
	ticket.ClaimedBy = i.Member.User.ID

	category := new(discordgo.Channel)

	// Ensure that the category exists for created tickets.
	category, err = a.Session().Channel(guild.Ticketing.ClaimedTicketsCategoryID)
	if err != nil {
		er := new(discordgo.RESTError)
		if errors.As(err, &er) && (er.Message.Code == discordgo.ErrCodeUnknownChannel || er.Message.Code == discordgo.ErrCodeGeneralError) { // General is thrown when a 404 is returned.
			a.Log().Warn("Claimed tickets category does not exist, creating it now")

			category, err = a.Session().GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
				Name: "Claimed Tickets",
				Type: discordgo.ChannelTypeGuildCategory,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					// Deny @everyone from seeing the ticket.
					{
						ID:    i.GuildID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: 0,
						Deny:  discordgo.PermissionAll,
					},
					// The creator of the ticket can see the ticket.
					{
						ID:    i.Member.User.ID,
						Type:  discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionAllText,
						Deny:  discordgo.PermissionMentionEveryone,
					},
					// Add the ticket role.
					{
						ID:    guild.Ticketing.RoleID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionAllText,
						Deny:  discordgo.PermissionMentionEveryone,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("error creating category: %w", err)
			}

			// Save the guild configuration.
			guild.Ticketing.ClaimedTicketsCategoryID = category.ID
			if err := a.GuildDal(context.Background()).SaveGuild(guild); err != nil {
				return fmt.Errorf("error saving guild configuration: %w", err)
			}
		} else {
			return fmt.Errorf("error getting category: %w", err)
		}
	} else if category != nil && category.ID != guild.Ticketing.ClaimedTicketsCategoryID {
		// Update the guild configuration.
		guild.Ticketing.ClaimedTicketsCategoryID = category.ID
		if err := a.GuildDal(context.Background()).SaveGuild(guild); err != nil {
			return fmt.Errorf("error saving guild configuration: %w", err)
		}
	}

	// Move the ticket to the claimed tickets' category.
	if _, err := a.Session().ChannelEditComplex(ticket.ChannelID, &discordgo.ChannelEdit{
		Name:     ticket.Name(),
		Position: &channel.Position,
		ParentID: category.ID,
	}); err != nil {
		return fmt.Errorf("error editing channel: %w", err)
	}

	// Save the ticket.
	if err := a.TicketDal(context.Background()).SaveTicket(ticket); err != nil {
		return fmt.Errorf("error saving ticket: %w", err)
	}

	// Set the claim button to be disabled.
	if err := setButtonDisabled(a, i, ClaimTicketButtonID, true); err != nil {
		return fmt.Errorf("error setting button disabled: %w", err)
	}

	// Respond to the interaction saying that the ticket has been claimed.
	err = a.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s>, you have claimed this ticket.", i.Member.User.ID),
		},
	})

	return nil
}

func setButtonDisabled(a IApp, i *discordgo.InteractionCreate, buttonID string, disabled bool) error {
	// Get the channel name.
	channel, err := a.Session().Channel(i.ChannelID)

	// Get the ticket.
	ticket, err := a.TicketDal(context.Background()).GetTicket(i.GuildID, channel.ID)
	if err != nil {
		return fmt.Errorf("error getting ticket: %w", err)
	}

	// Get the message.
	msg, err := a.Session().ChannelMessage(channel.ID, ticket.SetupMessageID)
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	// Get the button.
	button := new(discordgo.Button)
	for _, comp := range msg.Components {
		for _, component := range comp.(*discordgo.ActionsRow).Components {
			if component.(*discordgo.Button).CustomID == buttonID {
				button = component.(*discordgo.Button)
				break
			}
		}
	}

	// Set the button to be disabled.
	button.Disabled = disabled

	// Update the message.
	if _, err := a.Session().ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: channel.ID,
		ID:      msg.ID,
		Content: &NewTicketMessage.Content,
		Embed:   nil,
		Flags:   0,
		Components: []discordgo.MessageComponent{
			msg.Components[0],
		},
	}); err != nil {
		return fmt.Errorf("error editing message: %w", err)
	}

	return nil
}

func closeTicketHandler(a IApp, i *discordgo.InteractionCreate) error {
	// Get the channel name.
	channel, err := a.Session().Channel(i.ChannelID)

	// Get the ticket.
	ticket, err := a.TicketDal(context.Background()).GetTicket(i.GuildID, channel.ID)
	if err != nil {
		return fmt.Errorf("error getting ticket: %w", err)
	}

	// Get the guild configuration.
	guild, err := a.GuildDal(context.Background()).GetGuildByID(i.GuildID)
	if err != nil {
		return fmt.Errorf("error getting guild configuration: %w", err)
	}

	// Get the member that executed the command.
	member, err := a.Session().GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		return fmt.Errorf("error getting member: %w", err)
	}

	// Ensure that the user has the ticket role.
	if !hasRole(member, guild.Ticketing.RoleID) {
		err = respondEphemeral(a, i, "You do not have the ticket role to claim tickets. [<@&"+guild.Ticketing.RoleID+">]")
		if err != nil {
			return fmt.Errorf("error responding to interaction: %w", err)
		}
		return nil
	}

	// Ensure that the ticket is not already closed by using the category ID.
	if channel.ParentID == guild.Ticketing.ClosedTicketsCategoryID {
		err = respondEphemeral(a, i, "This ticket is already closed.")
		if err != nil {
			return fmt.Errorf("error responding to interaction: %w", err)
		}
		return nil
	}

	category := new(discordgo.Channel)

	// Ensure that the category exists for created tickets.
	category, err = a.Session().Channel(guild.Ticketing.ClosedTicketsCategoryID)
	if err != nil {
		er := new(discordgo.RESTError)
		if errors.As(err, &er) && (er.Message.Code == discordgo.ErrCodeUnknownChannel || er.Message.Code == discordgo.ErrCodeGeneralError) { // General is thrown when a 404 is returned.
			a.Log().Warn("Claimed tickets category does not exist, creating it now")

			category, err = a.Session().GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
				Name: "Closed Tickets",
				Type: discordgo.ChannelTypeGuildCategory,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					// Deny @everyone from seeing the ticket.
					{
						ID:    i.GuildID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: 0,
						Deny:  discordgo.PermissionAll,
					},
					// The creator of the ticket can see the ticket.
					{
						ID:    i.Member.User.ID,
						Type:  discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionAllText,
						Deny:  discordgo.PermissionMentionEveryone,
					},
					// Add the ticket role.
					{
						ID:    guild.Ticketing.RoleID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionAllText,
						Deny:  discordgo.PermissionMentionEveryone,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("error creating category: %w", err)
			}

			// Save the guild configuration.
			guild.Ticketing.ClosedTicketsCategoryID = category.ID
			if err := a.GuildDal(context.Background()).SaveGuild(guild); err != nil {
				return fmt.Errorf("error saving guild configuration: %w", err)
			}
		} else {
			return fmt.Errorf("error getting category: %w", err)
		}
	} else if category != nil && category.ID != guild.Ticketing.ClosedTicketsCategoryID {
		// Update the guild configuration.
		guild.Ticketing.ClosedTicketsCategoryID = category.ID
		if err := a.GuildDal(context.Background()).SaveGuild(guild); err != nil {
			return fmt.Errorf("error saving guild configuration: %w", err)
		}
	}

	// Move the ticket to the closed tickets' category.
	if _, err := a.Session().ChannelEditComplex(ticket.ChannelID, &discordgo.ChannelEdit{
		Name:     ticket.Name(),
		Position: &channel.Position,
		ParentID: category.ID,
	}); err != nil {
		return fmt.Errorf("error editing channel: %w", err)
	}

	go func() {
		// Set the close button to be disabled.
		if err := setButtonDisabled(a, i, CloseTicketButtonID, true); err != nil {
			a.Log().Error("Error setting close button disabled", slog.String(logging.KeyError, err.Error()))
		}

		// Set the reopen button to be enabled.
		if err := setButtonDisabled(a, i, ReopenTicketButtonID, false); err != nil {
			a.Log().Error("Error setting reopen button enabled", slog.String(logging.KeyError, err.Error()))
		}

		// Set the claim button to be disabled.
		if err := setButtonDisabled(a, i, ClaimTicketButtonID, true); err != nil {
			a.Log().Error("Error setting claim button disabled", slog.String(logging.KeyError, err.Error()))
		}

		// Set the delete button to be disabled.
		if err := setButtonDisabled(a, i, DeleteTicketButtonID, true); err != nil {
			a.Log().Error("Error setting delete button disabled", slog.String(logging.KeyError, err.Error()))
		}
	}()

	// Respond to the interaction saying that the ticket has been closed.
	err = a.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s>, congratulations on closing this ticket.", i.Member.User.ID),
		},
	})

	return nil
}
