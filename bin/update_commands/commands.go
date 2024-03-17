package main

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "status",
		Description: "Show bot status",
	},
	{
		Name:        "ranking",
		Description: "Server ranking",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "range",
				Description: "Range of the ranking",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "weekly",
						Value: "weekly",
					},
					{
						Name:  "monthly",
						Value: "monthly",
					},
				},
			},
		},
	},
	{
		Name:        "garden",
		Description: "Show my grass garden",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to show garden",
			},
		},
	},
}
