package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abiriadev/discord-garden/lib"
	"github.com/bwmarrin/discordgo"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var host = os.Getenv("INFLUX_HOST")
var qapi, wapi = lib.InitClient(host, os.Getenv("INFLUX"), "cl", "hello")

var gid = os.Getenv("GID")

func main() {
	tok := os.Getenv("TOK")

	discord, err := discordgo.New("Bot " + tok)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(msg)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"query": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var re string

			if res, err := qapi.QueryRaw(context.Background(), "", influxdb2.DefaultDialect()); err != nil {
				re = fmt.Sprintf("error:\n```%s```", err.Error())
			} else {
				// if _, err := s.ChannelMessageSend(m.ChannelID, ); err != nil
				// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("failed to send message:\n```%s```", err.Error()))
				re = fmt.Sprintf("res:\n```%s```", res)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: re,
				},
			})
		},
		"ranking": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			rank := lib.Rank(qapi)

			var buf bytes.Buffer

			var a = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "keycap_ten"}

			for i, r := range rank {
				buf.WriteString(fmt.Sprintf(":%s: <@%s>: %d\n", a[i], r.Id, r.Point))
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						&discordgo.MessageEmbed{
							Title:       "Ranking",
							Description: buf.String(),
							Color:       0x39d353,
						},
					},
				},
			})
		},
		"garden": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			res := lib.Garden(qapi)
			var buf bytes.Buffer

			for i := 0; i < 5; i++ {
				for j := 0; j < 6; j++ {
					buf.WriteString(fmt.Sprintf("%d ", res[i*5+j]))
				}
				buf.WriteString("\n")
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						&discordgo.MessageEmbed{
							Title:       "Garden",
							Description: buf.String(),
						},
					},
				},
			})
		},
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	defer discord.Close()
	if err := discord.Open(); err != nil {
		panic(err)
	}

	fmt.Println("running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func msg(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	lib.Record(wapi, m.Author.ID, 1, m.Timestamp)
	log.Printf("save msg @%s: %s\n", m.Author.Username, m.Content)
}
