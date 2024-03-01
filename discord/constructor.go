package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"log"
	"os"
	"porygon/config"
	"porygon/database"
	"reflect"
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
		return fmt.Sprintf("[%d]", level)
	}
	return replacement
}

func hasValues(data interface{}) bool {
	// check whatever a provided struct or type, has non-default value
	// currently supporting struct of integers and integers
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if hasValues(v.Index(i).Interface()) {
				return true
			}
		}
		return false
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i)
			switch fieldValue.Kind() {
			case reflect.Int:
				if fieldValue.Int() != 0 {
					return true
				}
			case reflect.Struct:
				if hasValues(fieldValue.Interface()) {
					return true
				}
			case reflect.Slice:
				for j := 0; j < fieldValue.Len(); j++ {
					if hasValues(fieldValue.Index(j).Interface()) {
						return true
					}
				}
			default:
				return true
			}
		}
		return false
	case reflect.Int:
		return v.Int() != 0
	default:
		return false
	}
}

func GenerateFields(gathered GatheredStats, config config.Config) []*discordgo.MessageEmbedField {
	currentTemplateFile, err := os.ReadFile("templates/current.override.json")
	if err != nil {
		currentTemplateFile, err = os.ReadFile("templates/current.json")
		if err != nil {
			log.Panicln(err)
		}
	}

	tmpl, err := template.New("message").Funcs(template.FuncMap{
		"Humanize":    humanizeValue,
		"HasValues":   hasValues,
		"LevelEmoji":  func(level int) string { return convertToEmoji(level, config.LevelEmoji) },
		"RewardEmoji": func(level int) string { return convertToEmoji(level, config.RewardEmoji) },
		"LureEmoji":   func(level int) string { return convertToEmoji(level, config.LureEmoji) },
		"RocketEmoji": func(level int) string { return convertToEmoji(level, config.RocketEmoji) },
		"EventEmoji":  func(level int) string { return convertToEmoji(level, config.EventEmoji) },
	}).Parse(string(currentTemplateFile))
	if err != nil {
		log.Panicln(err)
	}

	var resultJSON bytes.Buffer
	if err := tmpl.Execute(&resultJSON, gathered); err != nil {
		log.Panicln(err)
	}

	var fields []*discordgo.MessageEmbedField
	if err := json.Unmarshal(resultJSON.Bytes(), &fields); err != nil {
		log.Panicln(err)
	}

	return fields
}
