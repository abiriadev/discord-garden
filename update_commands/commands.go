package main

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ranking",
		Description: "Server ranking",
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
