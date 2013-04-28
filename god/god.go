package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func cleanup() {
	fmt.Println("starting cleanup")
	<-time.After(10 * time.Second)
	fmt.Println("cleanup complete")
}

func main() {
	defer cleanup()
	fmt.Println("waiting for a signal")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
	<-sigc
	fmt.Println("shutting down.")
}
