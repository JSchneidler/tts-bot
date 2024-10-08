module ttsbot

go 1.23.1

require (
	elevenlabs v0.0.0-00010101000000-000000000000
	github.com/bwmarrin/discordgo v0.28.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.24
	ttsmonster v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-resty/resty/v2 v2.13.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace elevenlabs => ./elevenlabs

replace ttsmonster => ./ttsmonster
