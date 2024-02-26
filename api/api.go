package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"porygon/config"
)

type apiResponse struct {
	SpawnId int `json:"spawn_id"`
}

type query struct {
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

func ApiRequest(config config.Config, ivMin, ivMax int) ([]apiResponse, error) {
	query := query{
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
	if config.API.Secret != "" {
		req.Header.Set("X-Golbat-Secret", config.API.Secret)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading API response: %w", err)
	}

	var apiResponses []apiResponse
	err = json.Unmarshal(body, &apiResponses)
	if err != nil {
		return nil, fmt.Errorf("error parsing API response: %w", err)
	}

	return apiResponses, nil
}
