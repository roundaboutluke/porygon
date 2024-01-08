package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"

	"Porygon/config"
	"Porygon/database"
	"Porygon/pokemon"
	"Porygon/query"
)

func saveMessageIDs(filename string, messageIDs map[string]string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("error creating message IDs file:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(messageIDs)
}

func loadMessageIDs(filename string) map[string]string {
	messageIDs := make(map[string]string)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("error creating message IDs file:", err)
			return messageIDs
		}
		defer file.Close()

		json.NewEncoder(file).Encode(messageIDs)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("error opening message IDs file:", err)
			return messageIDs
		}
		defer file.Close()

		json.NewDecoder(file).Decode(&messageIDs)
	}

	return messageIDs
}

func main() {
	var config config.Config

	if err := config.ParseConfig(); err != nil {
		panic(err)
	}

	messageIDs := loadMessageIDs("messageIDs.json")

	dg, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	go func() {
		for {
			db, err := database.DbConn(config)
			if err != nil {
				fmt.Println("error connecting to MariaDB,", err)
				continue
			}
			defer db.Close()

			var scannedCount, hundoCount, nundoCount, shinyCount, shinySpeciesCount int
			err = db.QueryRow("SELECT COALESCE(SUM(count), 0) FROM pokemon_stats WHERE date = CURDATE()").Scan(&scannedCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}

			err = db.QueryRow("SELECT COALESCE(SUM(count), 0) FROM pokemon_hundo_stats WHERE date = CURDATE()").Scan(&hundoCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}

			err = db.QueryRow("SELECT COALESCE(SUM(count), 0) FROM pokemon_nundo_stats WHERE date = CURDATE()").Scan(&nundoCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}

			err = db.QueryRow("SELECT COALESCE(SUM(count), 0) FROM pokemon_shiny_stats WHERE date = CURDATE()").Scan(&shinyCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}

			err = db.QueryRow("SELECT COUNT(DISTINCT pokemon_id) FROM pokemon_shiny_stats WHERE date = CURDATE()").Scan(&shinySpeciesCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}

			raidEggStats, err := database.RaidStats(db, config)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				continue
			}

			gymStats, err := database.GymStats(db, config)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				continue
			}

			pokestopStats, err := database.PokestopStats(db, config)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				continue
			}

			rewardStats, err := database.RewardStats(db, config)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				continue
			}

			lureStats, err := database.LureStats(db, config)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				continue
			}

			var hundoActiveCount, nundoActiveCount int
			if config.Config.IncludeActiveCounts {
				hundoApiResponses, err := query.ApiRequest(config, 15, 15)
				if err != nil {
					fmt.Println(err)
					db.Close()
					continue
				}

				hundoSpawnIds := make(map[int]bool)
				for _, apiResponse := range hundoApiResponses {
					hundoSpawnIds[apiResponse.SpawnId] = true
				}
				hundoActiveCount = len(hundoSpawnIds)

				nundoApiResponses, err := query.ApiRequest(config, 0, 0)
				if err != nil {
					fmt.Println(err)
					db.Close()
					continue
				}

				nundoSpawnIds := make(map[int]bool)
				for _, apiResponse := range nundoApiResponses {
					nundoSpawnIds[apiResponse.SpawnId] = true
				}
				nundoActiveCount = len(nundoSpawnIds)
			}

			hundoValue := humanize.Comma(int64(hundoCount))
			nundoValue := humanize.Comma(int64(nundoCount))

			if config.Config.IncludeActiveCounts {
				hundoValue = fmt.Sprintf("Active: %d | Today: %s", hundoActiveCount, hundoValue)
				nundoValue = fmt.Sprintf("Active: %d | Today: %s", nundoActiveCount, nundoValue)
			}

			rocketIncidents := []pokemon.Incident{
				{ID: 1, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Grunt)},
				{ID: 2, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Leader)},
				{ID: 3, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Giovanni)},
			}

			rocketStats := ""
			for _, incident := range rocketIncidents {
				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM incident WHERE display_type = ? AND expiration > UNIX_TIMESTAMP()", incident.ID).Scan(&count)
				if err != nil {
					fmt.Println("error querying MariaDB,", err)
					db.Close()
					continue
				}

				rocketStats += fmt.Sprintf("%s %d ", incident.Emoji, count)
			}

			var kecleonCount int
			err = db.QueryRow("SELECT COUNT(*) FROM incident WHERE display_type = ? AND expiration > UNIX_TIMESTAMP()", 8).Scan(&kecleonCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}
			kecleonStats := fmt.Sprintf("%s %d ", pokemon.FormatEmoji(config.Discord.Emojis.Kecleon), kecleonCount)

			var showcaseCount int
			err = db.QueryRow("SELECT COUNT(*) FROM incident WHERE display_type = ? AND expiration > UNIX_TIMESTAMP()", 9).Scan(&showcaseCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}
			showcaseStats := fmt.Sprintf("%s %d ", pokemon.FormatEmoji(config.Discord.Emojis.Showcase), showcaseCount)

			var activeRoutesCount int
			err = db.QueryRow("SELECT COUNT(*) FROM route WHERE type = 1").Scan(&activeRoutesCount)
			if err != nil {
				fmt.Println("error querying MariaDB,", err)
				db.Close()
				continue
			}
			activeRoutesStats := fmt.Sprintf("%s %d ", pokemon.FormatEmoji(config.Discord.Emojis.Route), activeRoutesCount)

			fields := []*discordgo.MessageEmbedField{
				{
					Name:   pokemon.FormatEmoji(config.Discord.Emojis.Scanned) + " Scanned",
					Value:  humanize.Comma(int64(scannedCount)),
					Inline: false,
				},
				{
					Name:   pokemon.FormatEmoji(config.Discord.Emojis.Hundo) + " Hundos",
					Value:  hundoValue,
					Inline: false,
				},
				{
					Name:   pokemon.FormatEmoji(config.Discord.Emojis.Nundo) + " Nundos",
					Value:  nundoValue,
					Inline: false,
				},
				{
					Name:   pokemon.FormatEmoji(config.Discord.Emojis.Shinies) + " Shinies",
					Value:  fmt.Sprintf("Species: %d | Total: %s", shinySpeciesCount, humanize.Comma(int64(shinyCount))),
					Inline: false,
				},
				{
					Name:   "Gym Statistics",
					Value:  gymStats,
					Inline: false,
				},
			}

			if raidEggStats != "" {
				newFields := make([]*discordgo.MessageEmbedField, len(fields)+1)

				copy(newFields, fields[:5])

				newFields[5] = &discordgo.MessageEmbedField{
					Name:   "Active Raids",
					Value:  raidEggStats,
					Inline: false,
				}

				copy(newFields[6:], fields[5:])

				fields = newFields
			}

			fields = append(fields, []*discordgo.MessageEmbedField{
				{
					Name:   "PokéStops Scanned",
					Value:  pokestopStats,
					Inline: false,
				},
				{
					Name:   "Quest Rewards",
					Value:  rewardStats,
					Inline: false,
				},
				{
					Name:   "Active Lures",
					Value:  lureStats,
					Inline: false,
				},
				{
					Name:   "Active Rockets",
					Value:  rocketStats,
					Inline: false,
				},
				{
					Name:   "Active Kecleon",
					Value:  kecleonStats,
					Inline: false,
				},
				{
					Name:   "Showcases",
					Value:  showcaseStats,
					Inline: false,
				},
				{
					Name:   "Routes",
					Value:  activeRoutesStats,
					Inline: false,
				},
			}...)

			embed := &discordgo.MessageEmbed{
				Title:     config.Config.EmbedTitle,
				Fields:    fields,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			for _, channelID := range config.Discord.ChannelIDs {
				var msg *discordgo.Message
				var err error

				var msgID string
				var ok bool

				if msgID, ok = messageIDs[channelID]; ok {
					msg, err = dg.ChannelMessageEditEmbed(channelID, msgID, embed)
					if err != nil {
						msg, err = dg.ChannelMessageSendEmbed(channelID, embed)
					}
				} else {
					msg, err = dg.ChannelMessageSendEmbed(channelID, embed)
				}

				if err != nil {
					fmt.Println("error sending or editing message in channel", channelID, ":", err)
					continue
				} else if msgID == "" || msgID != msg.ID {
					messageIDs[channelID] = msg.ID
					saveMessageIDs("messageIDs.json", messageIDs)
				}
			}

			db.Close()
			time.Sleep(time.Duration(config.Config.RefreshInterval) * time.Second)
		}
	}()

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Porygon is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
	return
}
