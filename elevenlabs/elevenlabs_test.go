package elevenlabs

import (
	"testing"
)

// func printResponse(resp *resty.Response, err error) {
// 	fmt.Println("Response Info:")
// 	fmt.Println("  Error      :", err)
// 	fmt.Println("  Status Code:", resp.StatusCode())
// 	fmt.Println("  Status     :", resp.Status())
// 	fmt.Println("  Proto      :", resp.Proto())
// 	fmt.Println("  Time       :", resp.Time())
// 	fmt.Println("  Received At:", resp.ReceivedAt())
// 	var pretty_json bytes.Buffer
// 	json.Indent(&pretty_json, []byte(resp.String()), "", "\t")
// 	fmt.Println(pretty_json.String())
// 	fmt.Println("  Body       :\n", resp)
// }

// func init() {
// 	godotenv.Load("../.env")
// }

func Test_GetVoices(t *testing.T) {
	_, err := GetVoices()

	if err != nil {
		t.Fatal(err)
	}
}

func Test_Tts(t *testing.T) {
	_, err := Tts("Test", default_voice_id)

	if err != nil {
		t.Fatal(err)
	}
}
