package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Say struct {
	Voice string
	Text  string
}

var _ = godotenv.Load()

var (
	bot_token_env      = "DISCORD_BOT_TOKEN"
	discord_server_id  = "DISCORD_SERVER_ID"
	allowed_channel_id = "DISCORD_ALLOWED_CHANNEL_ID"

	server_id  = os.Getenv(discord_server_id)
	channel_id = os.Getenv(allowed_channel_id)

	remove_commands = false
)

func main() {
	var err error

	bot_token := os.Getenv(bot_token_env)

	if len(bot_token) == 0 {
		log.Panicf(bot_token_env + " not set")
		os.Exit(1)
	}

	discord, _ := discordgo.New("Bot " + bot_token)

	err = discord.Open()
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
		if channel_id == "" || i.ChannelID == channel_id {
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
