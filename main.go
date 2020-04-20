package main

import (
	"os"
	"strconv"

	con "./Configurations"
	fsm "./FiniteStateMachine"
	mstr "./Master"
	network "./Network"
	bcast "./Network/bcast"
	peers "./Network/peers"
	elevio "./Hardware"
)

func main() {

	thisElevatorString := os.Args[1]
	localhost := "localhost:" + os.Args[2]
	thisElevator, _ := strconv.Atoi(thisElevatorString)

	elevio.Init(localhost, con.N_FLOORS)
	channels := fsm.StateMachineChannels{
		Elevator:       make(chan con.Elev),
		NewOrder:       make(chan elevio.ButtonEvent, 100),
		ArrivedAtFloor: make(chan int),
		DeleteQueue:    make(chan [con.N_FLOORS][con.N_BUTTONS]bool),
	}

	network := network.NetworkChannels{
		//from network to elevator controller
		UpdateMainLogic:      make(chan [con.N_ELEVS]con.Elev, 100),
		OnlineElevators:      make(chan [con.N_ELEVS]bool),
		ExternalOrderToLocal: make(chan con.Keypress),

		//from elevator to network controller
		LocalOrderToExternal:    make(chan con.Keypress),
		LocalElevatorToExternal: make(chan [con.N_ELEVS]con.Elev),

		//network controller to network
		OutgoingMsg:         make(chan con.Message),
		OutgoingOrder:       make(chan con.Keypress),
		PeersTransmitEnable: make(chan bool),

		//network to network controller
		IncomingMsg:   make(chan con.Message, 30),
		IncomingOrder: make(chan con.Keypress),
		PeerUpdate:    make(chan peers.PeerUpdate),
	}

	var (
		newOrder    = make(chan elevio.ButtonEvent)
		updateLight = make(chan [con.N_ELEVS]con.Elev)
	)

	msgpPort := 42000
	orderPort := 43000
	peersPort := 44000

	go elevio.PollButtons(newOrder)
	go elevio.PollFloorSensor(channels.ArrivedAtFloor)

	go fsm.RunFSM(channels, thisElevator)
	go mstr.ElevatorController(thisElevator, newOrder, updateLight, channels, network)
	go mstr.LightSetter(updateLight, thisElevator)
	go network.NetworkController(thisElevator, network)

	go bcast.Transmitter(msgpPort, network.OutgoingMsg)
	go bcast.Receiver(msgpPort, network.IncomingMsg)

	go bcast.Transmitter(orderPort, network.OutgoingOrder)
	go bcast.Receiver(orderPort, network.IncomingOrder)

	go peers.Transmitter(peersPort, thisElevatorString, network.PeersTransmitEnable)
	go peers.Receiver(peersPort, network.PeerUpdate)

	select {}
}
