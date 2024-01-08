package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"

	"Porygon/config"
	"Porygon/pokemon"
)

type GatheredStats struct {
	ScannedCount      int
	HundoValue        string
	NundoValue        string
	ShinyCount        int
	ShinySpeciesCount int

	RaidEggStats  string
	GymStats      string
	PokestopStats string
	RewardStats   string
	LureStats     string
	RocketStats   string

	KecleonStats      string
	ShowcaseStats     string
	ActiveRoutesStats string

	HundoActiveCount int
	NundoActiveCount int
}

func GenerateFields(config config.Config, gathered GatheredStats) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   pokemon.FormatEmoji(config.Discord.Emojis.Scanned) + " Scanned",
			Value:  humanize.Comma(int64(gathered.ScannedCount)),
			Inline: false,
		},
		{
			Name:   pokemon.FormatEmoji(config.Discord.Emojis.Hundo) + " Hundos",
			Value:  gathered.HundoValue,
			Inline: false,
		},
		{
			Name:   pokemon.FormatEmoji(config.Discord.Emojis.Nundo) + " Nundos",
			Value:  gathered.NundoValue,
			Inline: false,
		},
		{
			Name:   pokemon.FormatEmoji(config.Discord.Emojis.Shinies) + " Shinies",
			Value:  fmt.Sprintf("Species: %d | Total: %s", gathered.ShinySpeciesCount, humanize.Comma(int64(gathered.ShinyCount))),
			Inline: false,
		},
		{
			Name:   "Gym Statistics",
			Value:  gathered.GymStats,
			Inline: false,
		},
	}

	if gathered.RaidEggStats != "" {
		newFields := make([]*discordgo.MessageEmbedField, len(fields)+1)

		copy(newFields, fields[:5])

		newFields[5] = &discordgo.MessageEmbedField{
			Name:   "Active Raids",
			Value:  gathered.RaidEggStats,
			Inline: false,
		}

		copy(newFields[6:], fields[5:])

		fields = newFields
	}

	fields = append(fields, []*discordgo.MessageEmbedField{
		{
			Name:   "PokéStops Scanned",
			Value:  gathered.PokestopStats,
			Inline: false,
		},
		{
			Name:   "Quest Rewards",
			Value:  gathered.RewardStats,
			Inline: false,
		},
		{
			Name:   "Active Lures",
			Value:  gathered.LureStats,
			Inline: false,
		},
		{
			Name:   "Active Rockets",
			Value:  gathered.RocketStats,
			Inline: false,
		},
		{
			Name:   "Active Kecleon",
			Value:  gathered.KecleonStats,
			Inline: false,
		},
		{
			Name:   "Showcases",
			Value:  gathered.ShowcaseStats,
			Inline: false,
		},
		{
			Name:   "Routes",
			Value:  gathered.ActiveRoutesStats,
			Inline: false,
		},
	}...)
	return fields
}
