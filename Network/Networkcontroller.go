package NetworkController

import (
	"fmt"
	"strconv"
	"time"

	con "../Configurations"
	peers "../Network/peers"
)

type NetworkChannels struct {
	//FIXME: trenger annet navn for UpdateMainLogic

	UpdateMainLogic			chan [con.N_ElEVS]con.Elev
	OnlineElevators 		chan [con.N_ELEVS]bool
	ExternalOrderToLocal	chan con.Keypress

	LocalOrderToExternal 	chan con.Keypress
	LocalElevatorToExternal	chan [con.N_ELEVS]con.Elev


	OutgoingMsg				chan con.Message
	OutgoingOrder			chan con.Keypress
	PeersTransmitEnable 	chan bool


	IncomingMsg 			chan con.Message
	IncomingOrder 			chan con.Keypress
	PeerUpdate				chan peers.PeerUpdate
}


func NetworkController(thisElevator int, channel NetworkChannels) {
	var (
		msg 				con.Message
		onlineList			[con.N_ELEVS]bool
		outgoingOrder		con.Keypress
		incomingOrder		con.Keypress
	)

	PrimaryOrderTimer := time.NewTicker(100 * time.Millisecond)
	orderTicker := time.Newticker(10 * time.Millisecond)
	broadcastMsgTicker := time.Newticker(40 * time.Millisecond)
	deleteIncomingOrderTicker := time.NewTicker(1 * time.Second)
	orderTicker.Stop()

	msg.ID = thisElevator
	channel.PeersTransmitEnable <- true
	queue := make([]con.Keypress, 0)

	for {
		select {
		case msg.Elevator = <-channel.LocalElevatorToExternal:
		case ExternalOrder := <-channel.LocalOrderToExternal:
			queue = append(queue, ExternalOrder)

		case inOrder := <- channel.IncomingOrder:
			if inOrder.DesignatedElevator == thisElevator && incomingOrder != inOrder {
				incomingOrder = inOrder
				channel.ExternalOrderToLocal
			}
		case inMSG := <-channel.IncomingMsg:
			if inMSG.ID != thisElevator && inMSG.Elevator[inMSG.ID] != msgElevator[inMSG.ID] {
				msg.Elevator[inMSG.ID] = inMSG.Elevator[inMSG.ID]

				// FIXME: må finne nytt ord til UpdateMainLogic
				channel.UpdateMainLogic <- msg.Elevator
			}
		case broadcastMsgTicker.C:
			if onlineList[thisElevator] {
				channel.OutgoingMsg <- msg
			}
		case <-PrimaryOrderTimer.C:
			if len(queue) > 0 {
				outgoingOrder = queue[0]
				queue = queue[1:]
				orderTicker = time.NewTicker(10 * time.Millisecond)
			} else {
				orderTicker.Stop()
			}
		case <-orderTicker.C:
			channel.OutgoingOrder <- outgoingOrder
		}

	case <-deleteIncomingOrderTicker.C:
		incomingOrder = con.Keypress{Floor: -1}
	case peerUpdate := <- channel.PeerUpdate:
		if len(peerUpdate.Peers) == 0 {
			for id := 0; id < con.N_ELEVS; id++ {
				onlineList[id] = false
			}
		}
		if len(peerUpdate.New) > 0 {
			newElev, _ := strconv.Atoi(peerUpdate.New)
			onlineList[newElev] = true
		}
		if len(peerUpdate.Lost) > 0 {
			lostElev, _ := strconv.Atoi(peerUpdate.Lost[0])
			onlineList[lostElev] = false
		}
		go func() { channel.OnlineElevators <- onlineList }()

		fmt.Println("Number peers. ", len(peerUpdate.Peers))
		fmt.Println("New peers. ", peerUpdate.New)
		fmt.Println("Lost peers", peerUpdate.Lost)
	}
}