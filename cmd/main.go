package main

import "agora-activefence-kicker/service"

func main() {
	s := service.NewService()
	// Stop is called on another thread, but waits for an interrupt
	go s.Stop()
	s.Start()
}
