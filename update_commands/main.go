package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var tok = os.Getenv("DISCORD_TOKEN")
var gid = os.Getenv("GID")

func main() {
	discord, err := discordgo.New("Bot " + tok)
	if err != nil {
		panic(err)
	}

	log.Println("Updating slash commands")

	for _, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, gid, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		} else {
			log.Printf("added: %v", v.Name)
		}
	}
}
