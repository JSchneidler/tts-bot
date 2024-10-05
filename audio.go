package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var dca_path_env = "DCA_PATH"
var dca_path = os.Getenv(dca_path_env)

var CHUNK_SIZE = 8

var buffer = make([][]byte, 0)

// ffmpeg -i file.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca
func convert(input_path string) (output_path string, err error) {
	args := []string{
		"-i", input_path,
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	}
	ffmpeg := exec.Command("ffmpeg", args...)
	ffmpeg_out, err := ffmpeg.Output()
	if err != nil {
		log.Println("ffpmeg failed.", err)
	}

	dca := exec.Command(dca_path, "")
	dca.Stdin = bytes.NewReader(ffmpeg_out)
	dca_out, err := dca.Output()
	if err != nil {
		log.Println("dca failed.", err)
	}

	i := strings.Index(input_path, ".")
	path := input_path[:i] + ".dca"
	os.WriteFile(path, dca_out, 0644)
	return path, err
}

func loadSound(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, channelID string, mp3_path string) error {
	output_path, err := convert(mp3_path)
	err = loadSound(output_path)

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(server_id, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
