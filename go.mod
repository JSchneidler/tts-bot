module ttsbot

go 1.22.5

require (
	elevenlabs v0.0.0-00010101000000-000000000000
	github.com/bwmarrin/discordgo v0.28.1
	github.com/google/uuid v1.6.0
	github.com/jogramming/dca v0.0.0-20210930103944-155f5e5f0cc7
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/go-resty/resty/v2 v2.13.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jonas747/ogg v0.0.0-20161220051205-b4f6f4cf3757 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace elevenlabs => ./elevenlabs
