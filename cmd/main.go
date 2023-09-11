package main

import (
	"log"

	"github.com/hrvadl/go_aws_test/pkg/server"
)

func main() {
	srv := server.New()
	if err := srv.Setup(); err != nil {
		log.Fatalf("server failed to setup routes: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
