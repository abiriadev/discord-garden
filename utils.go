package main

import (
	"fmt"
	"strings"

	"github.com/abiriadev/discord-garden/lib"
	"github.com/bwmarrin/discordgo"
)

var primaryColor = 0x39d353

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

func embedifyRank(data []lib.RankRecord, rng string) *discordgo.MessageEmbed {
	var buf strings.Builder

	for i, r := range data {
		buf.WriteString(
			fmt.Sprintf(":%s: <@%s>: %d\n", numberEmojiList[i], r.Id, r.Point),
		)
	}

	var rngText string
	switch rng {
	case "weekly":
		rngText = "Weekly Ranking"
		break
	case "monthly":
		rngText = "Monthly Ranking"
		break
	case "all":
		rngText = "Total Ranking"
	default:
		rngText = "Ranking"
	}

	return &discordgo.MessageEmbed{
		Title:       rngText,
		Description: buf.String(),
		Color:       primaryColor,
	}
}
