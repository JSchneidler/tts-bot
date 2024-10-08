# TTS Bot
A TTS bot written in Go supporting a few TTS APIs.

Currently supports:
- ElevenLabs
- TTS Monster

# Dependencies
- Go
- FFmpeg
- [DCA](https://github.com/bwmarrin/dca)

# Environment
Supports .env file or standard environment variables. Currently all these must be set.  

`DCA_PATH` - Path to DCA executable

`DISCORD_BOT_TOKEN` - Bot token for Discord auth
`DISCORD_SERVER_ID` - Discord server ID
`DISCORD_ALLOWED_CHANNEL_ID` - Discord channel ID to restrict messages to (currently must be set)

`ELEVENLABS_API_KEY` - ElevenLabs API key (xi-api-key)
`TTSMONSTER_API_KEY` - TTS Monster API key