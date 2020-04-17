//Bruker main filen som lå ute i nettverksmodulen på GitHub som kladd"

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"./bcast"
	"./localip"
	"./peers"
	//"https://github.com/maxmrmr/TTK4145-Sanntid/Network/bcast"
	//"https://github.com/maxmrmr/TTK4145-Sanntid/Network/localip"
	//"https://github.com/maxmrmr/TTK4145-Sanntid/Network/peers"
)

// Defining sone struct to send over the network
// All members we want to transmit must be public. Any private members will be received as zero-values
type HelloMsg struct {
	Message string
	Iter    int
}

func main() {
	fmt.Println("Hello world!")
	//
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peers-%s-%d", localIP, os.Getpid())
	}

	//Make a channel to receive updates on the id's of the peers that are "alive" on the network
	peerUpdateCh := make(chan peers.PeerUpdate)

	//Writes id to port 15647 and
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable) //FIXME sjekk tall
	go peers.Receiver(15647, peerUpdateCh)

	//Channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	//Sender og lytter til meldinger på Aurora sin favorittkanal
	go bcast.Transmitter(33333, helloTx)
	go bcast.Receiver(33333, helloRx)

	go func() {
		helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
