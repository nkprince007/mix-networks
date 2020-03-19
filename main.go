package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/nkprince007/mix-networks/server"
)

func main() {
	done := make(chan os.Signal, 1)
	go signal.Notify(done, os.Interrupt, os.Kill)

	srv := server.Server{Addr: ":8080"}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Server error on %v, %v\n", srv.Addr, err)
		}
	}()

	<-done
	if err := srv.ShutDown(); err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
	os.Exit(0)
}
