package main

import (
	"clickclack/sound"
	"flag"
	"log"
	"os"
	"os/signal"

	hook "github.com/robotn/gohook"
)

func main() {
	var soundPackID string
	var clickLeft string
	var clickRight string
	var volume float64

	flag.StringVar(&soundPackID, "k", "", "Sound pack ID (e.g., 1203000000018)")
	flag.StringVar(&clickLeft, "lc", "", "Left mouse click sound path.")
	flag.StringVar(&clickRight, "rc", "", "Right mouse click sound path.")
	flag.Float64Var(&volume, "v", 0.0, "Volume level (0.0 to 1.0)")
	flag.Parse()

	if volume < 0.0 || volume > 1.0 {
		log.Fatalf("Volume must be between 0.0 and 1.0")
	}

	signal.Notify(sound.SigChan, os.Interrupt)

	conf, err := sound.InitConfig(volume, soundPackID, clickLeft, clickRight)
	if err != nil {
		panic(err)
	}

	eventChan := hook.Start()
	defer hook.End()

	sound.CreateSound(conf, eventChan)
}
