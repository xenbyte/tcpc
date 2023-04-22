package main

import (
	"fmt"
	"log"
	"time"

	"github.com/xenbyte/tcpc/tcpc"
)

var (
	LocalPort  = ":3000"
	RemotePort = ":4000"
)

func main() {

	channelLocal, err := tcpc.New[string](LocalPort, RemotePort)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second)
		channelLocal.SendChan <- "Message from TCPC"
	}()

	channelRemote, err := tcpc.New[string](RemotePort, LocalPort)

	if err != nil {
		log.Fatal(err)
	}

	msg := <-channelRemote.RecvChan
	fmt.Println("message: ", msg)
}
