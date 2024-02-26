package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"

	"porygon/api"
	"porygon/config"
	"porygon/database"
	"porygon/discord"
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

func gatherStats(db *sqlx.DB, config config.Config) (discord.GatheredStats, error) {
	start := time.Now()
	var err error
	var gathered discord.GatheredStats

	gathered.Pokemon, err = database.GetPokeStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.RaidEgg, err = database.GetRaidStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Gym, err = database.GetGymStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Pokestop, err = database.GetPokestopStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Reward, err = database.GetRewardStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Lure, err = database.GetLureStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Rocket, err = database.GetRocketStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Event, err = database.GetEventStats(db)
	if err != nil {
		return gathered, err
	}

	gathered.Route, err = database.GetRoutesStats(db)
	if err != nil {
		return gathered, err
	}

	// probs break this out into query? again idk how to handle passing the config well just yet
	if config.Config.IncludeActiveCounts {
		hundoApiResponses, err := api.ApiRequest(config, 15, 15)
		if err != nil {
			return gathered, err
		}

		hundoSpawnIds := make(map[int]bool)
		for _, apiResponse := range hundoApiResponses {
			hundoSpawnIds[apiResponse.SpawnId] = true
		}
		gathered.HundoActiveCount = len(hundoSpawnIds)

		nundoApiResponses, err := api.ApiRequest(config, 0, 0)
		if err != nil {
			return gathered, err
		}

		nundoSpawnIds := make(map[int]bool)
		for _, apiResponse := range nundoApiResponses {
			nundoSpawnIds[apiResponse.SpawnId] = true
		}
		gathered.NundoActiveCount = len(nundoSpawnIds)
	}

	elapsed := time.Since(start)
	fmt.Printf("Fetched stats in %s\n", elapsed)
	return gathered, nil
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
				time.Sleep(time.Duration(config.Config.ErrorRefreshInterval) * time.Second)
				continue
			}
			gathered, err := gatherStats(db, config)
			db.Close()
			if err != nil {
				fmt.Println("failed to fetch stats,", err)
				time.Sleep(time.Duration(config.Config.ErrorRefreshInterval) * time.Second)
				continue
			}

			fields := discord.GenerateFields(gathered, config)

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
