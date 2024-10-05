package main

import (
	"log"
	"os"

	"elevenlabs"
	"ttsmonster"

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
					Choices:     toCommandOptions(ttsmonster.Voices),
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
			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			playSound(s, vs.ChannelID, "f15b66cf-3cee-455f-8751-2f06910e7a39.wav")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Test.",
				},
			})
		},
		"say": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			options := i.ApplicationCommandData().Options

			text := options[0].StringValue()

			voice_id := ""
			if len(options) > 1 {
				voice_id = options[1].StringValue()
			}

			log.Printf("%s(%s) used /say: (%s) %s", i.Member.User.GlobalName, i.Member.User.Username, voice_id, text)

			uuid := uuid.NewString()

			wav_data, err := ttsmonster.Tts(text, voice_id)
			if err == nil {
				wav_path := uuid + ".wav"
				os.WriteFile(wav_path, wav_data, 0644)
				log.Println(wav_path)
				playSound(s, vs.ChannelID, wav_path)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Say. WIP.",
				},
			})
		},
		"sfx": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			options := i.ApplicationCommandData().Options

			text := options[0].StringValue()

			log.Printf("%s(%s) used /sfx: %s", i.Member.User.GlobalName, i.Member.User.Username, text)

			uuid := uuid.NewString()

			mp3_data, err := elevenlabs.Sfx(text)
			if err == nil {
				mp3_path := uuid + ".mp3"
				os.WriteFile(mp3_path, mp3_data, 0644)
				log.Println(mp3_path)
				playSound(s, vs.ChannelID, mp3_path)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "SFX. WIP.",
				},
			})
		},
	}
)

func toCommandOptions(voices []ttsmonster.Voice) []*discordgo.ApplicationCommandOptionChoice {
	commandOptions := []*discordgo.ApplicationCommandOptionChoice{}

	for _, voice := range voices {
		commandOptions = append(commandOptions, &discordgo.ApplicationCommandOptionChoice{
			Name:  voice.Name,
			Value: voice.VoiceID,
		})
	}

	return commandOptions
}
