package elevenlabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
)

func printResponse(resp *resty.Response, err error) {
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	var pretty_json bytes.Buffer
	json.Indent(&pretty_json, []byte(resp.String()), "", "\t")
	fmt.Println(pretty_json.String())
	fmt.Println("  Body       :\n", resp)
}

func TestApi(t *testing.T) {
	resp, err := get_voices()
	printResponse(resp, err)
}

func TestTts(t *testing.T) {
	path, _ := tts("Test", default_voice_id, "test.mp3")
	fmt.Println("Downloaded TTS: " + path)
	// printResponse(resp, err)
}
