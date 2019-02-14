package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/euforia/metermaid"
)

func main() {
	mm, err := metermaid.New()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	mm.Stop()
}
