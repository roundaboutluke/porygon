package discord

import (
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionAdministrator

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "list-emotes",
			Description:              "List emotes",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:                     "create-emotes",
			Description:              "Create emotes",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"list-emotes":   listEmotes,
		"create-emotes": createEmotes,
	}
)

func listEmotes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var emotesList strings.Builder

	guildEmotes, _ := s.GuildEmojis(i.GuildID)

	emotesList.WriteString("```")
	for _, emote := range guildEmotes {
		emotesList.WriteString(fmt.Sprintf("<:%s:%s>\n", emote.Name, emote.ID))
	}
	emotesList.WriteString("```")

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: emotesList.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func createEmotes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var output strings.Builder
	output.WriteString("```")

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	emotesDir := "emojis/override"
	files, err := filepath.Glob(filepath.Join(emotesDir, "*.png"))

	if len(files) > 0 {
		output.WriteString("Using `emojis/override` directory as emotes source.\n")
	} else {
		output.WriteString("Using `emojis` directory as emotes source.\n")
		emotesDir = "emojis"
		files, err = filepath.Glob(filepath.Join(emotesDir, "*.png"))
	}

	if err != nil {
		log.Println("Error reading emotes directory:", err)
		return
	}

	// fetch existing emotes
	guildEmotes, _ := s.GuildEmojis(i.GuildID)
	existingEmotes := make(map[string]bool)
	for _, emote := range guildEmotes {
		existingEmotes[emote.Name] = true
	}

	// check and upload every emote we have under emotesDir
	for _, file := range files {
		emoteName := strings.TrimSuffix(filepath.Base(file), ".png")

		if _, exists := existingEmotes[emoteName]; exists {
			output.WriteString(fmt.Sprintf("%s - already there\n", emoteName))
			continue
		}
		emoteFile, err := os.ReadFile(filepath.Join(emotesDir, filepath.Base(file)))
		encodedImage := base64.StdEncoding.EncodeToString(emoteFile)
		dataURI := fmt.Sprintf("data:image/png;base64,%s", encodedImage)

		_, err = s.GuildEmojiCreate(i.GuildID, &discordgo.EmojiParams{
			Name:  emoteName,
			Image: dataURI,
		})
		if err != nil {
			output.WriteString(fmt.Sprintf("%s - upload error: %s\n", emoteName, err))
			continue
		}
		output.WriteString(fmt.Sprintf("%s - success\n", emoteName))
	}
	output.WriteString("```")

	_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: output.String(),
	})
}
