package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	hook "github.com/robotn/gohook"
)

type SoundPack struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	KeyDefineType  string                 `json:"key_define_type"`
	IncludesNumpad bool                   `json:"includes_numpad"`
	Sound          string                 `json:"sound"`
	Defines        map[string]interface{} `json:"defines"`
	Tags           []string               `json:"tags"`
}

var (
	already_pressed = map[uint16]bool{}
)

func main() {
	var soundPackID string
	var volume float64

	flag.StringVar(&soundPackID, "k", "1200000000001", "Sound pack ID (e.g., 1203000000018)")
	flag.Float64Var(&volume, "v", 1.0, "Volume level (0.0 to 1.0)")
	flag.Parse()

	if soundPackID == "" {
		log.Fatalf("Sound pack ID (-k) is required")
	}

	if volume < 0.0 || volume > 1.0 {
		log.Fatalf("Volume must be between 0.0 and 1.0")
	}

	soundPackUrls := []string{
		fmt.Sprintf("https://mechvibes.com/sound-packs/sound-pack-%s", soundPackID),
		fmt.Sprintf("https://mechvibes.com/sound-packs/custom-sound-pack-%s", soundPackID),
	}
	working_url_index := 0

	packageJSONPath := filepath.Join("./data", soundPackID, "config.json")

	var soundPack SoundPack
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {

		for i := 0; i < len(soundPackUrls); i++ {
			packageJSONURL := soundPackUrls[i] + "/dist/config.json"
			if err := downloadFile(packageJSONURL, packageJSONPath); err != nil {
				if len(soundPackUrls)-1 == i {
					log.Fatalf("Failed to download package.json: %v", err)
				}
			} else {
				working_url_index = i
				break
			}
		}

		fmt.Printf("Downloaded package.json to %s\n", packageJSONPath)
	}

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		log.Fatalf("Failed to read package.json: %v", err)
	}

	if err := json.Unmarshal(data, &soundPack); err != nil {
		log.Fatalf("Failed to parse package.json: %v", err)
	}

	fmt.Printf("Loaded sound pack: %s (%s)\n", soundPack.Name, soundPack.ID)

	uniqueSounds := getUniqueSounds(soundPack)
	soundMap := make(map[string]*beep.Buffer)
	fmt.Printf("Unique sounds: %v\n", uniqueSounds)

	for _, sound := range uniqueSounds {
		soundURL := fmt.Sprintf("%s/dist/%s", soundPackUrls[working_url_index], sound)
		soundPath := filepath.Join("./data", soundPackID, sound)

		if err := downloadFile(soundURL, soundPath); err != nil {
			fmt.Printf("Failed to download sound %s: %v\n", sound, err)
			continue
		}
		fmt.Printf("Sound ready: %s to %s\n", sound, soundPath)

		soundMap[sound] = loadSound(soundPath)
	}

	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	eventChan := hook.Start()
	defer hook.End()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	fmt.Printf("Listening for keyboard and mouse events...\n")
	fmt.Printf("Sound pack ID: %s\n", soundPackID)
	fmt.Printf("Volume: %.2f\n", volume)

	for {
		select {
		case event := <-eventChan:
			if event.Kind == hook.KeyDown {

				if already_pressed[event.Rawcode] {
					continue
				}

				already_pressed[event.Rawcode] = true

				switch v := soundPack.Defines[appMap[event.Rawcode]].(type) {
				case string:

					buffer, has := soundMap[v]
					if has {
						playSound(buffer, 0, 100, volume)
					}

				case []interface{}:
					if len(v) == 2 {
						startMs := int(v[0].(float64))
						durationMs := int(v[1].(float64))

						buffer, has := soundMap[uniqueSounds[0]]

						if has {
							playSound(buffer, startMs, durationMs, volume)
						}
					}
				}
				fmt.Printf("Current key %d\n", event.Rawcode)

			} else if event.Kind == hook.MouseDown {
				fmt.Printf("Mouse clicked: %v\n", event.Button)

			} else if event.Kind == hook.KeyUp {
				already_pressed[event.Rawcode] = false
			}
		case <-sigChan:
			fmt.Println("Exiting...")
			return
		}
	}
}

func getUniqueSounds(soundPack SoundPack) []string {
	uniqueSounds := make(map[string]bool)

	if soundPack.Sound != "" {
		uniqueSounds[soundPack.Sound] = true
	}

	for _, value := range soundPack.Defines {
		switch v := value.(type) {
		case string:
			uniqueSounds[v] = true
		}
	}

	var sounds []string
	for sound := range uniqueSounds {
		sounds = append(sounds, sound)
	}

	return sounds
}

func downloadFile(url, path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func loadSound(filename string) *beep.Buffer {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open sound file: %v", err)
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch {
	case len(filename) > 4 && filename[len(filename)-4:] == ".mp3":
		streamer, format, err = mp3.Decode(f)
	case len(filename) > 4 && filename[len(filename)-4:] == ".wav":
		streamer, format, err = wav.Decode(f)
	case len(filename) > 4 && filename[len(filename)-4:] == ".ogg":
		streamer, format, err = vorbis.Decode(f)
	default:
		log.Fatalf("Unsupported file format: %s", filename)
	}

	if err != nil {
		log.Fatalf("Failed to decode sound file: %v", err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return buffer
}

func playSound(buffer *beep.Buffer, startMs, durationMs int, volume float64) {
	startSamples := buffer.Format().SampleRate.N(time.Duration(startMs) * time.Millisecond)
	durationSamples := buffer.Format().SampleRate.N(time.Duration(durationMs) * time.Millisecond)

	streamer := buffer.Streamer(startSamples, startSamples+durationSamples)

	pitchRatio := 1.0 + (0.1 * (float64(time.Now().UnixNano()%3) - 1))
	resampled := beep.ResampleRatio(4, pitchRatio, streamer)

	volumeStreamer := &volumeCtrl{resampled, volume}

	speaker.Play(volumeStreamer)
}

type volumeCtrl struct {
	streamer beep.Streamer
	volume   float64
}

func (v *volumeCtrl) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = v.streamer.Stream(samples)
	for i := range samples[:n] {
		samples[i][0] *= v.volume
		samples[i][1] *= v.volume
	}
	return n, ok
}

func (v *volumeCtrl) Err() error {
	if s, ok := v.streamer.(beep.StreamCloser); ok {
		return s.Err()
	}
	return nil
}
