package elevenlabs

import (
	"os"

	"github.com/go-resty/resty/v2"
)

var api_key_env = "ELEVENLABS_API_KEY"
var api_key = os.Getenv(api_key_env)

var api_base = "https://api.elevenlabs.io/v1"

var default_voice_id = "YmkIUFWsIp3y1ScFnXsd"

type VoiceSettings struct {
	Stability       float32 `json:"stability"`
	SimilarityBoost float32 `json:"similarity_boost"`
	Style           float32 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

type TtsBody struct {
	Text          string        `json:"text"`
	ModelId       string        `json:"model_id"`
	VoiceSettings VoiceSettings `json:"voice_settings"`
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

func tts(text string, voice_id string, output_path string) (string, error) {
	tts_body := TtsBody{
		Text:    text,
		ModelId: "eleven_multilingual_v2",
		VoiceSettings: VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.8,
			Style:           0.0,
			UseSpeakerBoost: true,
		},
	}

	client := get_client()
	resp, err := client.R().SetBody(tts_body).Post(api_base + "/text-to-speech/" + voice_id)

	_ = os.WriteFile(output_path, resp.Body(), 0644)

	return output_path, err
}
