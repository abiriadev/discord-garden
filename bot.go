package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abiriadev/discord-garden/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	InfluxHost        string `envconfig:"INFLUX_HOST"        required:"true"`
	InfluxToken       string `envconfig:"INFLUX_TOKEN"       required:"true"`
	InfluxOrg         string `envconfig:"INFLUX_ORG"         required:"true"`
	InfluxBucket      string `envconfig:"INFLUX_BUCKET"      required:"true"`
	InfluxMeasurement string `envconfig:"INFLUX_MEASUREMENT" required:"true"`
	DiscordGuildId    string `envconfig:"DISCORD_GUILD_ID"   required:"true"`
	DiscordToken      string `envconfig:"DISCORD_TOKEN"      required:"true"`
}

var handlers = []any{
	func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	},
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	influxclient := lib.NewInfluxClient(lib.InfluxClientConfig{
		Host:        config.InfluxHost,
		Token:       config.InfluxToken,
		Org:         config.InfluxOrg,
		Bucket:      config.InfluxBucket,
		Measurement: config.InfluxMeasurement,
	})

	s, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		panic(err)
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	for _, handler := range handlers {
		s.AddHandler(handler)
	}

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		influxclient.Record(m.Author.ID, 1, m.Timestamp)
		log.Printf("save msg @%s: %s\n", m.Author.Username, m.Content)
	})

	s.AddHandler(
		map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"ranking": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				optMap := make(map[string]string)
				for _, v := range i.ApplicationCommandData().Options {
					optMap[v.Name] = v.StringValue()
				}

				rng, ok := optMap["range"]
				if !ok {
					rng = "all"
				}

				rank, err := influxclient.Rank(rng)

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Title: "title",
						Embeds: []*discordgo.MessageEmbed{
							embedifyRank(rank, rng),
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
		},
	)

	defer s.Close()
	if err := s.Open(); err != nil {
		panic(err)
	}

	fmt.Println("running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
