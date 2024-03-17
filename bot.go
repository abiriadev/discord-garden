package main

import (
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

	for _, handler := range Handlers {
		s.AddHandler(handler)
	}

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		influxclient.Record(m.Author.ID, 1, m.Timestamp)
		log.Printf("save msg @%s: %s\n", m.Author.Username, m.Content)
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
