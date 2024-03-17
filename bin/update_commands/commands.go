package main

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
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
	},
	{
		Name:        "query",
		Description: "Random query",
	},
}
