package main

import (
	"fmt"
	"wzrds/common"
	"wzrds/server/internal/network"
)

func main() {
	netServer := network.NewNetworkServer()

	for {
		eventStruct := netServer.CheckForEvents()
		switch msg := eventStruct.(type) {
		case common.ClientConnectionInfo:
			fmt.Println("mess", msg)
		}
	}
}
