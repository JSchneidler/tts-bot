package main

import (
	"testing"
)

var mp3_path = "4d6aa57c-ab0e-43c9-ac4a-869f9f3de1e6.mp3"

func TestConvert(t *testing.T) {
	_, err := convert(mp3_path)
	if err != nil {
		t.Fatal(err)
	}
}
