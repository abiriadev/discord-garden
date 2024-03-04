package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	tok := os.Getenv("TOK")

	discord, err := discordgo.New("Bot " + tok)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(msg)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	if err := discord.Open(); err != nil {
		panic(err)
	}

	fmt.Println("running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func msg(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		log.Println(m.Content)

		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
