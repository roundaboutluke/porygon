package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Name     string `json:"name"`
	} `json:"database"`
	Discord struct {
		Token     string `json:"token"`
		ChannelID string `json:"channelID"`
		Emojis    struct {
			Valor      string `json:"valor"`
			Mystic     string `json:"mystic"`
			Instinct   string `json:"instinct"`
			Uncontested   string `json:"uncontested"`
			Normal string `json:"normal"`
			Glacial string `json:"glacial"`
			Mossy string `json:"mossy"`
			Magnetic string `json:"magnetic"`
			Rainy string `json:"rainy"`
			Sparkly string `json:"sparkly"`
		} `json:"emojis"`
	} `json:"discord"`
	API struct {
		URL    string `json:"url"`
		Secret string `json:"secret"`
	} `json:"api"`
	Coordinates struct {
		Min struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"min"`
		Max struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"max"`
	} `json:"coordinates"`
	Config struct {
	RefreshInterval int `json:"refreshInterval"`
	} `json:"config"`
	MessageID string `json:"messageID"`
}

type ApiResponse struct {
	SpawnId int `json:"spawn_id"`
}

type Query struct {
	Min struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"min"`
	Max struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"max"`
	Filters []struct {
		Pokemon []struct {
			Id int `json:"id"`
		} `json:"pokemon"`
		AtkIv struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"atk_iv"`
		DefIv struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"def_iv"`
		StaIv struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"sta_iv"`
	} `json:"filters"`
}

func apiRequest(config Config, ivMin, ivMax int) ([]ApiResponse, error) {
	query := Query{
		Min: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  config.Coordinates.Min.Latitude,
			Longitude: config.Coordinates.Min.Longitude,
		},
		Max: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  config.Coordinates.Max.Latitude,
			Longitude: config.Coordinates.Max.Longitude,
		},
		Filters: []struct {
			Pokemon []struct {
				Id int `json:"id"`
			} `json:"pokemon"`
			AtkIv struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"atk_iv"`
			DefIv struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"def_iv"`
			StaIv struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"sta_iv"`
		}{
			{
				Pokemon: func() []struct {
					Id int `json:"id"`
				} {
					pokemon := make([]struct {
						Id int `json:"id"`
					}, 1015)
					for i := range pokemon {
						pokemon[i].Id = i + 1
					}
					return pokemon
				}(),
				AtkIv: struct {
					Min int `json:"min"`
					Max int `json:"max"`
				}{
					Min: ivMin,
					Max: ivMax,
				},
				DefIv: struct {
					Min int `json:"min"`
					Max int `json:"max"`
				}{
					Min: ivMin,
					Max: ivMax,
				},
				StaIv: struct {
					Min int `json:"min"`
					Max int `json:"max"`
				}{
					Min: ivMin,
					Max: ivMax,
				},
			},
		},
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error converting query to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", config.API.URL+"/api/pokemon/v2/scan", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating API request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Golbat-Secret", config.API.Secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading API response: %w", err)
	}

	var apiResponses []ApiResponse
	err = json.Unmarshal(body, &apiResponses)
	if err != nil {
		return nil, fmt.Errorf("error parsing API response: %w", err)
	}

	return apiResponses, nil
}

func saveMessageID(config *Config, messageID string) {
	config.MessageID = messageID
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println("error marshalling config:", err)
		return
	}
	err = ioutil.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		fmt.Println("error writing config file:", err)
	}
}


func loadMessageID(config Config) string {
	return config.MessageID
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("error opening config file,", err)
		return
	}
	defer file.Close()

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println("error decoding config file,", err)
		return
	}

	messageID := loadMessageID(config)

	dg, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	go func() {
		for {
			db, err := sql.Open("mysql", config.Database.Username+":"+config.Database.Password+"@tcp("+config.Database.Host+":"+config.Database.Port+")/"+config.Database.Name)
			if err != nil {
				fmt.Println("error connecting to MariaDB,", err)
				continue
			}

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

			type Team struct {
				ID    int
				Emoji string
			}

			teams := []Team{
				{ID: 1, Emoji: "<" + config.Discord.Emojis.Valor + ">"}, 
				{ID: 2, Emoji: "<" + config.Discord.Emojis.Mystic + ">"}, 
				{ID: 3, Emoji: "<" + config.Discord.Emojis.Instinct + ">"}, 
				{ID: 0, Emoji: "<" + config.Discord.Emojis.Uncontested + ">"}, 
			}

			gymStats := ""
			for _, team := range teams {
				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE team_id = ? AND updated > UNIX_TIMESTAMP() - 4 * 60 * 60", team.ID).Scan(&count)
				if err != nil {
					fmt.Println("error querying MariaDB,", err)
					db.Close()
					continue
				}

				gymStats += fmt.Sprintf("%s %d ", team.Emoji, count)
			}

			type Lure struct {
				ID    int
				Emoji string
			}

			lures := []Lure{
				{ID: 501, Emoji: "<" + config.Discord.Emojis.Normal + ">"}, 
				{ID: 502, Emoji: "<" + config.Discord.Emojis.Glacial + ">"}, 
				{ID: 503, Emoji: "<" + config.Discord.Emojis.Mossy + ">"}, 
				{ID: 504, Emoji: "<" + config.Discord.Emojis.Magnetic + ">"}, 
				{ID: 505, Emoji: "<" + config.Discord.Emojis.Rainy + ">"}, 
				{ID: 506, Emoji: "<" + config.Discord.Emojis.Sparkly + ">"}, 
			}

			lureStats := ""
			for _, lure := range lures {
				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM pokestop WHERE lure_id = ? AND lure_expire_timestamp > UNIX_TIMESTAMP()", lure.ID).Scan(&count)
				if err != nil {
					fmt.Println("error querying MariaDB,", err)
					db.Close()
					continue
				}

				lureStats += fmt.Sprintf("%s %d ", lure.Emoji, count)
			}

			hundoApiResponses, err := apiRequest(config, 15, 15)
			if err != nil {
				fmt.Println(err)
				db.Close()
				continue
			}

			hundoSpawnIds := make(map[int]bool)
			for _, apiResponse := range hundoApiResponses {
				hundoSpawnIds[apiResponse.SpawnId] = true
			}
			hundoActiveCount := len(hundoSpawnIds)

			nundoApiResponses, err := apiRequest(config, 0, 0)
			if err != nil {
				fmt.Println(err)
				db.Close()
				continue
			}

			nundoSpawnIds := make(map[int]bool)
			for _, apiResponse := range nundoApiResponses {
				nundoSpawnIds[apiResponse.SpawnId] = true
			}
			nundoActiveCount := len(nundoSpawnIds)

			embed := &discordgo.MessageEmbed{
				Title: "Today's Pokémon Stats",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "📈 Scanned",
						Value:  fmt.Sprintf("%d", scannedCount),
						Inline: false,
					},
					{
						Name:   "💯 Hundos",
						Value:  fmt.Sprintf("Active: %d | Today: %d", hundoActiveCount, hundoCount),
						Inline: false,
					},
					{
						Name:   "🗑️ Nundos",
						Value:  fmt.Sprintf("Active: %d | Today: %d", nundoActiveCount, nundoCount),
						Inline: false,
					},
					{
						Name:   "✨ Shinies",
						Value:  fmt.Sprintf("Species: %d | Total: %d", shinySpeciesCount, shinyCount),
						Inline: false,
					},
					{
						Name:   "Gym Statistics",
						Value:  gymStats,
						Inline: false,
					},
					{
						Name:   "Active Lures",
						Value:  lureStats,
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339), 
			}

			var msg *discordgo.Message

			if messageID != "" {
				msg, err = dg.ChannelMessageEditEmbed(config.Discord.ChannelID, messageID, embed)
				if err != nil {
					fmt.Println("error editing message,", err)
					msg, err = dg.ChannelMessageSendEmbed(config.Discord.ChannelID, embed)
					if err != nil {
						fmt.Println("error sending embed to Discord channel,", err)
						db.Close()
						continue
					}
				}
			} else {
				msg, err = dg.ChannelMessageSendEmbed(config.Discord.ChannelID, embed)
				if err != nil {
					fmt.Println("error sending embed to Discord channel,", err)
					db.Close()
					continue
				}
			}

			messageID = msg.ID

			saveMessageID(&config, messageID)

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
