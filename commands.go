package main

import (
	"log"
	"os"

	"elevenlabs"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var (
	defaultMemberPermission int64 = discordgo.PermissionManageChannels
	commands                      = []*discordgo.ApplicationCommand{
		{
			Name:        "info",
			Description: "Bot info",
		},
		{
			Name:                     "test",
			Description:              "Test TTS",
			DefaultMemberPermissions: &defaultMemberPermission,
		},
		{
			Name:                     "say",
			Description:              "Text to Speech",
			DefaultMemberPermissions: &defaultMemberPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "text",
					Description: "What to say",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "voice",
					Description: "Which voice to use",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Default",
							Value: "default",
						},
						{
							Name:  "Second",
							Value: "second",
						},
					},
				},
			},
		},
		{
			Name:                     "sfx",
			Description:              "Sound Effects",
			DefaultMemberPermissions: &defaultMemberPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "text",
					Description: "Description of sound to generate",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	command_handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Info. WIP.",
				},
			})
		},
		"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mp3_data, err := os.ReadFile("4d6aa57c-ab0e-43c9-ac4a-869f9f3de1e6.mp3")
			if err != nil {
				log.Println("Failed to load mp3.", err)
				return
			}

			dca_data, err := convert(mp3_data)
			if err != nil {
				log.Println("Failed to convert mp3.", err)
				return
			}

			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			playSound(s, vs.ChannelID, dca_data)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Test.",
				},
			})
		},
		"say": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			text := options[0].StringValue()

			voice_id := ""
			if len(options) > 1 {
				voice_id = options[1].StringValue()
			}

			log.Printf("%s(%s) used /say: (%s) %s", i.Member.User.GlobalName, i.Member.User.Username, voice_id, text)

			uuid := uuid.New()

			mp3_data, err := elevenlabs.Sfx(text)
			if err == nil {
				mp3_path := uuid.String() + ".mp3"
				os.WriteFile(mp3_path, mp3_data, 0644)
				log.Println(mp3_path)
			}

			dca_data, err := convert(mp3_data)
			if err == nil {
				dca_path := uuid.String() + ".dca"
				os.WriteFile(dca_path, dca_data, 0644)
				log.Println(dca_path)
			}

			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			playSound(s, vs.ChannelID, dca_data)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Say. WIP.",
				},
			})
		},
		"sfx": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			text := options[0].StringValue()

			log.Printf("%s(%s) used /sfx: %s", i.Member.User.GlobalName, i.Member.User.Username, text)

			uuid := uuid.New()

			mp3_data, err := elevenlabs.Sfx(text)
			if err == nil {
				mp3_path := uuid.String() + ".mp3"
				os.WriteFile(mp3_path, mp3_data, 0644)
				log.Println(mp3_path)
			}

			dca_data, err := convert(mp3_data)
			if err == nil {
				dca_path := uuid.String() + ".dca"
				os.WriteFile(dca_path, dca_data, 0644)
				log.Println(dca_path)
			}

			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			playSound(s, vs.ChannelID, dca_data)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "SFX. WIP.",
				},
			})
		},
	}
)
