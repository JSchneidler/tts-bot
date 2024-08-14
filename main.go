package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type Say struct {
	Voice string
	Text  string
}

var (
	bot_token_env      = "DISCORD_BOT_TOKEN"
	discord_server_id  = "DISCORD_SERVER_ID"
	allowed_channel_id = "DISCORD_ALLOWED_CHANNEL_ID"

	remove_commands = false

	defaultMemberPermission int64 = discordgo.PermissionManageChannels
	commands                      = []*discordgo.ApplicationCommand{
		{
			Name:        "info",
			Description: "Bot info",
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
		"say": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//options := i.ApplicationCommandData().Options
			log.Printf("%s(%s) used /say: (%s) %s", i.Member.User.GlobalName, i.Member.User.Username, "WIP", "WIP")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Say. WIP.",
				},
			})
		},
	}
)

func main() {
	bot_token := os.Getenv(bot_token_env)
	server_id := os.Getenv(discord_server_id)
	channel_id := os.Getenv(allowed_channel_id)

	if len(bot_token) == 0 {
		log.Panicf(bot_token_env + " not set")
		os.Exit(1)
	}

	discord, _ := discordgo.New("Bot " + bot_token)

	err := discord.Open()
	if err != nil {
		log.Panicf("could not open session: %s", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, server_id, v)
		if err != nil {
			log.Panicf("Cannot create '%s' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
		log.Printf("Created '%s' command", v.Name)
	}

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Logged in as", r.User.String())
		s.UpdateCustomStatus("/info")
	})

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ChannelID == channel_id {
			if h, ok := command_handlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		}
	})

	if remove_commands {
		log.Println("Removing all slash commands...")

		commands, _ := discord.ApplicationCommands(discord.State.User.ID, "")

		for _, v := range commands {
			log.Println(v.Name)
			err := discord.ApplicationCommandDelete(discord.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%s' command: %v", v.Name, err)
			}
		}
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = discord.Close()
	if err != nil {
		log.Panicf("could not close session gracefully: %s", err)
	}
}

// func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

// 	// Join the provided voice channel.
// 	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
// 	if err != nil {
// 		return err
// 	}

// 	// Sleep for a specified amount of time before playing the sound
// 	time.Sleep(250 * time.Millisecond)

// 	// Start speaking.
// 	vc.Speaking(true)

// 	// Send the buffer data.
// 	for _, buff := range buffer {
// 		vc.OpusSend <- buff
// 	}

// 	// Stop speaking
// 	vc.Speaking(false)

// 	// Sleep for a specificed amount of time before ending.
// 	time.Sleep(250 * time.Millisecond)

// 	// Disconnect from the provided voice channel.
// 	vc.Disconnect()

// 	return nil
// }
