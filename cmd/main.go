package main

import kickService "github.com/AgoraIO-Community/agora-activefence-kicker/service"

func main() {
	s := kickService.NewService()
	// Stop is called on another thread, but waits for an interrupt
	go s.Stop()
	s.Start()
}
