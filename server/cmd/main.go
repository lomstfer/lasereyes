package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/server/internal/network"
)

func main() {
	netServer := network.NewNetworkServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for {
			eventStruct := netServer.CheckForEvents()
			switch msg := eventStruct.(type) {
			case msgfromclient.ConnectionInfo:
				fmt.Println("mess", msg)
			case msgfromclient.MoveInput:
				fmt.Println(msg.Input)
			}
		}
	}()

	<-sigChan
	netServer.Stop()
}
