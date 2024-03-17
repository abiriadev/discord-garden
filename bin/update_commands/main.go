package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var tok = os.Getenv("DISCORD_TOKEN")
var gid = os.Getenv("GID")

func main() {
	s, err := discordgo.New("Bot " + tok)
	if err != nil {
		panic(err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	log.Println("Opening session")
	s.Open()

	log.Println("Updating slash commands")

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, gid, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		} else {
			log.Printf("added: %v", v.Name)
		}
	}

	log.Println("Closing session")
	s.Close()
}
