package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abiriadev/discord-garden/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

var primaryColor = 0x39d353
var errorColor = 0xD35039

var numberEmojiList = []string{
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
	"keycap_ten",
}

var grassEmojiList = []string{
	"<:grass0:1219046954798026943>",
	"<:grass1:1219046956840652871>",
	"<:grass2:1219046959113965688>",
	"<:grass3:1219046961475489802>",
	"<:grass4:1219046963366985809>",
}

func optMap(
	i *discordgo.InteractionCreate,
) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, v := range i.ApplicationCommandData().Options {
		optMap[v.Name] = v
	}
	return optMap
}

func makeErrorResponse(err error) *discordgo.InteractionResponse {
	if err == nil {
		err = errors.New("unknown error")
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title: "title",
			Embeds: []*discordgo.MessageEmbed{
				&discordgo.MessageEmbed{
					Title:       "Error",
					Description: fmt.Sprintf("```golang\n%s\n```", err.Error()),
					Color:       errorColor,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Open new issue on repository",
					Style: discordgo.DangerButton,
					URL:   "https://github.com/abiriadev/discord-garden/issues/new",
				},
			},
		},
	}
}

func embedifyRank(data []lib.RankRecord, rng string) *discordgo.MessageEmbed {
	var buf strings.Builder

	for i, r := range data {
		buf.WriteString(
			fmt.Sprintf(":%s: <@%s>: %d\n", numberEmojiList[i], r.Id, r.Point),
		)
	}

	rngText := lo.Switch[string, string](rng).
		Case("all", "Total Ranking").
		Case("weekly", "Weekly Ranking").
		Case("monthly", "Monthly Ranking").
		Default("Ranking")

	return &discordgo.MessageEmbed{
		Title:       rngText,
		Description: buf.String(),
		Color:       primaryColor,
	}
}

func embedifyGarden(data []int, username string) *discordgo.MessageEmbed {
	res := lib.ApplyHistogram(data, len(grassEmojiList)-1, lib.BinaryMeanHistogram{})

	eRes := lo.Map(res, func(v, _ int) string {
		return grassEmojiList[v]
	})

	var buf strings.Builder

	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			buf.WriteString(eRes[j*5+i])
		}
		buf.WriteString("\n")
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Garden", username),
		Description: buf.String(),
		Color:       primaryColor,
	}
}

func embedifyStatus(client lib.InfluxClient) (*discordgo.MessageEmbed, error) {
	influx, error := client.Status()
	if error != nil {
		return nil, error
	}

	return &discordgo.MessageEmbed{
		Title: "Status",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Influx",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "name",
				Value:  influx.Health.Name,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "version",
				Value:  lo.FromPtrOr(influx.Health.Version, "null"),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "commit",
				Value:  lo.FromPtrOr(influx.Health.Commit, "null"),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "status",
				Value:  string(influx.Health.Status),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "started",
				Value:  influx.Ready.Started.Local().Format(time.RFC3339),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "up",
				Value:  lo.FromPtrOr(influx.Ready.Up, "null"),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "message",
				Value:  lo.FromPtrOr(influx.Health.Message, "null"),
				Inline: true,
			},
		},
		Color: primaryColor,
	}, nil
}
