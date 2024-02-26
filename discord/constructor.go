package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"os"
	"porygon/config"
	"porygon/database"
	"strconv"
	"text/template"
)

type GatheredStats struct {
	Pokemon database.PokeStats
	RaidEgg []database.RaidStats
	Gym     database.GymStats
	Reward  []database.TypeCountStats
	Lure    []database.TypeCountStats
	Rocket  []database.TypeCountStats
	Event   []database.TypeCountStats

	Pokestop int
	Route    int

	HundoActiveCount int
	NundoActiveCount int
}

func humanizeValue(value int) string {
	return humanize.Comma(int64(value))
}

func convertToEmoji(level int, config map[string]string) string {
	replacement, ok := config[strconv.Itoa(level)]
	if !ok {
		return fmt.Sprintf("Rocket %d", level)
	}
	return replacement
}

func GenerateFields(gathered GatheredStats, config config.Config) []*discordgo.MessageEmbedField {
	currentTemplateFile, err := os.ReadFile("templates/current.override.json")
	if err != nil {
		currentTemplateFile, err = os.ReadFile("templates/current.json")
		if err != nil {
			panic(err)
		}
	}

	tmpl, err := template.New("message").Funcs(template.FuncMap{
		"Humanize":    humanizeValue,
		"LevelEmoji":  func(level int) string { return convertToEmoji(level, config.LevelEmoji) },
		"RewardEmoji": func(level int) string { return convertToEmoji(level, config.RewardEmoji) },
		"LureEmoji":   func(level int) string { return convertToEmoji(level, config.LureEmoji) },
		"RocketEmoji": func(level int) string { return convertToEmoji(level, config.RocketEmoji) },
		"EventEmoji":  func(level int) string { return convertToEmoji(level, config.EventEmoji) },
	}).Parse(string(currentTemplateFile))
	if err != nil {
		panic(err)
	}

	var resultJSON bytes.Buffer
	if err := tmpl.Execute(&resultJSON, gathered); err != nil {
		panic(err)
	}

	var fields []*discordgo.MessageEmbedField
	if err := json.Unmarshal(resultJSON.Bytes(), &fields); err != nil {
		panic(err)
	}

	return fields
}
