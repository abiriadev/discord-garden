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
	"github.com/mstoykov/envconfig"
)

var qapi, wapi = lib.InitClient(host, os.Getenv("INFLUX"), "cl", "hello")

type Config struct {
	InfluxHost       string `envconfig:"INFLUX_HOST"      required:"true"`
	InfluxToken      string `envconfig:"INFLUX_TOKEN"     required:"true"`
	DISCORD_GUILD_ID string `envconfig:"DISCORD_GUILD_ID" required:"true"`
	DISCORD_TOKEN    string `envconfig:"DISCORD_TOKEN"    required:"true"`
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	tok := os.Getenv("TOK")

	s, err := discordgo.New("Bot " + tok)
	if err != nil {
		panic(err)
	}

	s.AddHandler(msg)

	s.Identify.Intents = discordgo.IntentsGuildMessages

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
			optMap := make(map[string]string)
			for _, v := range i.ApplicationCommandData().Options {
				optMap[v.Name] = v.StringValue()
			}

			rng, ok := optMap["range"]
			var rngText string
			if !ok {
				rng = "all"
				rngText = "Total Ranking"
			} else {
				switch rng {
				case "weekly":
					rngText = "Weekly Ranking"
					break
				case "monthly":
					rngText = "Monthly Ranking"
					break
				}
			}

			rank := lib.Rank(qapi, rng)

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
							Title:       rngText,
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

			fmt.Println(res)

			for i := 0; i < 5; i++ {
				for j := 0; j < 6; j++ {
					buf.WriteString(fmt.Sprintf("%d ", res[i*6+j]))
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

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	defer s.Close()
	if err := s.Open(); err != nil {
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
