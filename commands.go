package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"elevenlabs"
	"ttsmonster"

	"github.com/bwmarrin/discordgo"
)

var (
	defaultMemberPermission int64 = discordgo.PermissionManageChannels

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "info",
			Description: "Bot info",
		},
		// {
		// 	Name:                     "test",
		// 	Description:              "Test TTS",
		// 	DefaultMemberPermissions: &defaultMemberPermission,
		// },
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
			stats, err := ttsmonster.GetStats()

			var content string
			if err == nil {
				const template = "Character usage: %d/%d (%.1f%%)\nRenewal Date: %s"

				percent := (float32(stats.Usage) / float32(stats.Limit)) * 100
				renewal_date := time.Unix(int64(stats.RenewalTimestamp), 0)
				content = fmt.Sprintf(template, stats.Usage, stats.Limit, percent, renewal_date.Format(time.RFC850))
			} else {
				content = "Error. See logs."
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		// "test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 	vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
		// 	playSound(s, vs.ChannelID, "f15b66cf-3cee-455f-8751-2f06910e7a39.wav")

		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
		// 		Data: &discordgo.InteractionResponseData{
		// 			Content: "Test.",
		// 		},
		// 	})
		// },
		"say": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			options := i.ApplicationCommandData().Options

			text := options[0].StringValue()

			voice_id := ""
			if len(options) > 1 {
				voice_id = options[1].StringValue()
			}

			log.Printf("%s(%s) used /say: (%s) %s", i.Member.User.GlobalName, i.Member.User.Username, voice_id, text)
			// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
			// 	Data: &discordgo.InteractionResponseData{
			// 		Content: text,
			// 	},
			// })

			wav_data, err := ttsmonster.Tts(text, voice_id)
			if err == nil {
				sound_path := saveSound(wav_data, "wav")
				user := DiscordUser{ID: i.Member.User.ID, Name: i.Member.User.GlobalName}
				usage := Usage{AudioType: AUDIO_TYPE_TTS, AudioService: AUDIO_SERVICE_TTSMONSTER, Prompt: text, AudioFilename: filepath.Base(sound_path)}
				AddUsage(user, usage)
				playSound(s, vs.ChannelID, sound_path)
			}
		},
		"sfx": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vs, _ := s.State.VoiceState(server_id, i.Member.User.ID)
			options := i.ApplicationCommandData().Options

			text := options[0].StringValue()

			log.Printf("%s(%s) used /sfx: %s", i.Member.User.GlobalName, i.Member.User.Username, text)
			// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
			// 	Data: &discordgo.InteractionResponseData{
			// 		Content: text,
			// 	},
			// })

			mp3_data, err := elevenlabs.Sfx(text)
			if err == nil {
				sound_path := saveSound(mp3_data, "mp3")
				user := DiscordUser{ID: i.Member.User.ID, Name: i.Member.User.GlobalName}
				usage := Usage{AudioType: AUDIO_TYPE_SFX, AudioService: AUDIO_SERVICE_ELEVENLABS, Prompt: text, AudioFilename: filepath.Base(sound_path)}
				AddUsage(user, usage)
				playSound(s, vs.ChannelID, sound_path)
			}
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
