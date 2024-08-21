package main

import (
	"os"
	"testing"
)

var mp3_path = "4d6aa57c-ab0e-43c9-ac4a-869f9f3de1e6.mp3"

func TestConvert(t *testing.T) {
	mp3_data, err := os.ReadFile(mp3_path)
	if err != nil {
		t.Fatal(err)
	}

	_, err = convert(mp3_data)
	if err != nil {
		t.Fatal(err)
	}
}
