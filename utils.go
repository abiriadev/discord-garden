package main

import (
	"fmt"
	"strings"

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
	"grass0",
	"grass1",
	"grass2",
	"grass3",
	"grass4",
}

func optMap(i *discordgo.InteractionCreate) map[string]string {
	optMap := make(map[string]string)
	for _, v := range i.ApplicationCommandData().Options {
		optMap[v.Name] = v.StringValue()
	}
	return optMap
}

func makeErrorResponse(err error) *discordgo.InteractionResponse {
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

func embedifyGarden(data []int) *discordgo.MessageEmbed {
	var buf strings.Builder

	for i, r := range data {
		buf.WriteString(
			fmt.Sprintf(":%s: <@%s>: %d\n", numberEmojiList[i], r.Id, r.Point),
		)
	}

	return &discordgo.MessageEmbed{
		Title:       rngText,
		Description: buf.String(),
		Color:       primaryColor,
	}
}
