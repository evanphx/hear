package main

import (
	"fmt"
	"log"
	"os"
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

	fmt.Printf("recognizing...\n")

	gcp, err := hear.NewGCPSpeechConv(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatal(err)
	}

	words, err := gcp.Convert(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("=> %s\n", words)
}
