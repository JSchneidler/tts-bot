package elevenlabs

import (
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

var api_key_env = "ELEVENLABS_API_KEY"
var api_key = os.Getenv(api_key_env)

var api_base = "https://api.elevenlabs.io/v1"

var default_voice_id = "YmkIUFWsIp3y1ScFnXsd"

type voiceSettings struct {
	Stability       float32 `json:"stability"`
	SimilarityBoost float32 `json:"similarity_boost"`
	Style           float32 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

type ttsBody struct {
	Text          string        `json:"text"`
	ModelId       string        `json:"model_id"`
	VoiceSettings voiceSettings `json:"voice_settings"`
}

type sfxBody struct {
	Text string `json:"text"`
}

func get_client() *resty.Client {
	client := resty.New()
	client.Header.Add("xi-api-key", api_key)
	return client
}

func get_voices() (*resty.Response, error) {
	client := get_client()
	resp, err := client.R().Get(api_base + "/voices")
	return resp, err
}

func Tts(text string, voice_id string) (mp3_data []byte, err error) {
	tts_body := ttsBody{
		Text:    text,
		ModelId: "eleven_multilingual_v2",
		VoiceSettings: voiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.8,
			Style:           0.0,
			UseSpeakerBoost: true,
		},
	}

	if voice_id == "" {
		voice_id = default_voice_id
	}

	log.Printf("Calling ElevenLabs TTS: (%s) %s", voice_id, text)

	client := get_client()
	resp, err := client.R().SetBody(tts_body).Post(api_base + "/text-to-speech/" + voice_id)

	if resp.StatusCode() == 200 {
		return resp.Body(), nil
	}
	return []byte{}, err
}

func Sfx(text string) (mp3_data []byte, err error) {
	sfx_body := sfxBody{
		Text: text,
	}

	log.Printf("Calling ElevenLabs SFX: %s", text)

	client := get_client()
	resp, err := client.R().SetBody(sfx_body).Post(api_base + "/sound-generation/")

	if resp.StatusCode() == 200 {
		return resp.Body(), nil
	}
	return []byte{}, err
}
