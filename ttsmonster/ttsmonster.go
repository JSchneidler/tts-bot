package ttsmonster

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

var api_key_env = "TTSMONSTER_API_KEY"
var api_key = os.Getenv(api_key_env)

var api_base = "https://api.console.tts.monster"

type ttsBody struct {
	VoiceId string `json:"voice_id"`
	Message string `json:"message"`
}

type ttsResponse struct {
	Status int    `json:"status"`
	Url    string `json:"url"`
}

func get_client() *resty.Client {
	client := resty.New()
	client.SetHeader("Authorization", api_key)
	return client
}

func get_voices() (*resty.Response, error) {
	client := get_client()
	resp, err := client.R().Get(api_base + "/voices")
	return resp, err
}

func Tts(text string, voice_id string) (wav_data []byte, err error) {
	tts_body := ttsBody{
		Message: text,
		VoiceId: voice_id,
	}

	log.Printf("Calling TTSMonster generate: (%s) %s", voice_id, text)

	client := get_client()
	resp, err := client.R().SetBody(tts_body).Post(api_base + "/generate")

	if resp.StatusCode() != 200 {
		return []byte{}, err
	}

	log.Printf("TTSMonster generate response: %s", resp.Body())

	var response ttsResponse
	json.Unmarshal(resp.Body(), &response)

	if err != nil {
		return []byte{}, err
	}

	resp, err = client.R().Get(response.Url)

	if resp.StatusCode() != 200 {
		return []byte{}, err
	}

	return resp.Body(), nil
}
