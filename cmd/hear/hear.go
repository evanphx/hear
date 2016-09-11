package main

import (
	"fmt"
	"log"
	"time"

	"github.com/evanphx/hear"
)

func main() {
	hear.InitAudio()
	defer hear.FreeAudio()

	opts := hear.ListenOpts{
		QuietDuration:    1 * time.Second,
		AlreadyListening: true,
	}

	fmt.Printf("speak now\n")

	buf, err := hear.ListenIntoBuffer(opts)
	if err != nil {
		log.Fatal(err)
	}

	gcp, err := hear.NewGCPSpeechConv("/Users/evan/.gcloud/hear.json")
	if err != nil {
		log.Fatal(err)
	}

	words, err := gcp.Convert(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("=> %s\n", words)
}
