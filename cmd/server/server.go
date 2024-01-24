package main

import (
	"fmt"
	"l0/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := server.NewServer()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\r")
		s.Down()
		os.Exit(0)
	}()

	if err := s.Up(); err != nil {
		panic(err)
	}
}
