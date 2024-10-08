package main

import (
	"testing"
)

const discord_id = "1234567890"
const discord_name = "TestUser"

func Test_add_user(t *testing.T) {
	_, err := add_user(discord_id, discord_name)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_add_usage(t *testing.T) {
	err := add_usage(Usage{1, AUDIO_TYPE_TTS, AUDIO_SERVICE_TTSMONSTER, "test_prompt", "test_path.mp3"})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AddUsage(t *testing.T) {
	user := DiscordUser{discord_id, discord_name}
	AddUsage(user, Usage{1, AUDIO_TYPE_TTS, AUDIO_SERVICE_TTSMONSTER, "A bunch of blah blah blah more more more wah wah wah", "test_path.mp3"})
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
