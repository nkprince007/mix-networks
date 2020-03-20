package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/nkprince007/mix-networks/mixes"
)

func feeder(m mixes.Mix) {
	for {
		time.Sleep(time.Duration(rand.Intn(5)*100) * time.Millisecond)
		msg := mixes.Message(rand.Int())
		m.AddMessage(msg)
	}
}

func listenToSignals(done chan interface{}) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	fmt.Println("signal received: ", <-s)
	fmt.Println("shutting down...")
	done <- nil
}

func main() {
	done := make(chan interface{}, 1)
	go listenToSignals(done)

	var mix mixes.Mix = &mixes.ThresholdMix{Size: 4}
	go feeder(mix)
	for {
		select {
		case <-done:
			mix.CleanUp()
			return
		default:
			fmt.Println(mix.GetMessages())
		}
	}
}

// func main() {
// 	done := make(chan os.Signal, 1)
// 	go signal.Notify(done, os.Interrupt, os.Kill)

// 	srv := server.Server{
// 		Addr:            ":8080",
// 		IdleConnTimeout: time.Second * 15,
// 		BufferSize:      1024,
// 	}

// 	go func() {
// 		err := srv.ListenAndServe()
// 		if err != nil {
// 			log.Fatalf("Server error on %v, %v\n", srv.Addr, err)
// 		}
// 	}()

// 	<-done
// 	if err := srv.ShutDown(); err != nil {
// 		log.Fatalf("Shutdown error: %v\n", err)
// 	}
// 	os.Exit(0)
// }
