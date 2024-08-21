package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
)

var dca_path_env = "DCA_PATH"
var dca_path = os.Getenv(dca_path_env)

var CHUNK_SIZE = 8

// var buffer = make([][]byte, 0)

// ffmpeg -i file.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca
func convert(mp3_data []byte) (dca_data []byte, err error) {
	args := []string{
		"-i", "pipe:",
		//"-f", "s16le",
		"-f", "mp3",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	}
	ffmpeg := exec.Command("ffmpeg", args...)
	ffmpeg.Stdin = bytes.NewReader(mp3_data)
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

	return dca_out, err
}

// func convert(mp3_path string) (dca_path string, err error) {
// 	encodeSession, err := dca.EncodeFile(mp3_path, dca.StdEncodeOptions)
// 	if err != nil {
// 		log.Println(err)
// 		return "", err
// 	}
// 	defer encodeSession.Cleanup()

// 	output, err := os.Create(mp3_path + ".dca")
// 	if err != nil {
// 		log.Println(err)
// 		return "", err
// 	}

// 	copied, err := io.Copy(output, encodeSession)
// 	if err != nil {
// 		log.Println(err)
// 		return "", err
// 	}
// 	log.Println(copied)

// 	return output.Name(), nil
// }

// func loadSound() error {
// 	//file, err := os.Open("4d6aa57c-ab0e-43c9-ac4a-869f9f3de1e6.mp3.dca")
// 	file, err := os.Open("test.dca")
// 	if err != nil {
// 		fmt.Println("Error opening dca file :", err)
// 		return err
// 	}

// 	var opuslen int16

// 	for {
// 		// Read opus frame length from dca file.
// 		err = binary.Read(file, binary.LittleEndian, &opuslen)

// 		// If this is the end of the file, just return.
// 		if err == io.EOF || err == io.ErrUnexpectedEOF {
// 			err := file.Close()
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		}

// 		if err != nil {
// 			fmt.Println("Error reading from dca file :", err)
// 			return err
// 		}

// 		// Read encoded pcm from dca file.
// 		InBuf := make([]byte, opuslen)
// 		err = binary.Read(file, binary.LittleEndian, &InBuf)

// 		// Should not be any end of file errors
// 		if err != nil {
// 			fmt.Println("Error reading from dca file :", err)
// 			return err
// 		}

// 		// Append encoded pcm data to the buffer.
// 		buffer = append(buffer, InBuf)
// 	}
// }

func playSound(s *discordgo.Session, channelID string, dca_data []byte) (err error) {
	// loadSound()

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
	// for _, buff := range buffer {
	// 	vc.OpusSend <- buff
	// }
	for i := 0; i < len(dca_data); i += CHUNK_SIZE {
		end := i + CHUNK_SIZE
		if end > len(dca_data) {
			end = len(dca_data)
		}
		vc.OpusSend <- dca_data[i:end]
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
