package Master

import (
	"fmt"

	elevio "../Hardware"
	con "../Configurations"
	network "../Network"
	fsm "../FiniteStateMachine"
)





func elevator_controller(thisElevator int, SyncChannels network.NetworkChannels, localStateChannel fsm.StateMachineChannels, elevatorcontrollers <- chan elevio.ButtonEvent, Lights chan <- [con.N_ELEVS]con.Elev) {

	var (
		elevatorList [con.N_ELEVS]con.Elev //Takes in the struct elev with info about alle elevators
		onlineElevators [con.N_ELEVS]bool //list of online elevators
		temp_Keypress con.Keypress
		temp_ButtonEvent elevio.ButtonEvent

	)
	fmt.Println("Starting elevator controller:", thisElevator)
	
	for {
		select {
		case newLocalOrder := <-elevatorcontrollers:
			id := costCalculator(thisElevator, elevatorList, newLocalOrder, onlineElevators)
			if id != -1 {
				if id == thisElevator {
					localStateChannel.NewOrder <- newLocalOrder
				} else {
					temp_Keypress = con.Keypress{DesignatedElevator: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button}
					SyncChannels.LocalOrderToExternal <- temp_Keypress
				}
			}
		case temp_Keypress = <- SyncChannels.ExternalOrderToLocal:
			temp_ButtonEvent = elevio.ButtonEvent{Button: temp_Keypress.Button, Floor: temp_Keypress.Floor}
			if elevatorList[thisElevator].State == con.Undefined {
				costID := costCalculator(thisElevator, elevatorList, temp_ButtonEvent,  onlineElevators)
				temp_Keypress.DesignatedElevator = costID
				SyncChannels.LocalOrderToExternal <- temp_Keypress
			} else {
				localStateChannel.NewOrder <- temp_ButtonEvent
			}
		case NewUpdateLocalElevator := <-localStateChannel.Elevator:

			change := false
			for floor:= 0; floor < con.N_FLOORS; floor++ {
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button ++ {
					if elevatorList[thisElevator].Queue[floor][button] && !NewUpdateLocalElevator.Queue[floor][button] {
						change = true
					}
				}
				if change {
					change = false
					for id := 0; id < con.N_ELEVS; id++ {
						if id != thisElevator {
							elevatorList[thisElevator].Queue[floor][elevio.BT_HallUp] = false
							elevatorList[thisElevator].Queue[floor][elevio.BT_HallDown] = false
						}
					}
				}
			}
			change = false
			if elevatorList[thisElevator].State != con.Undefined && NewUpdateLocalElevator.State == con.Undefined {
				elevatorList[thisElevator].State = con.Undefined
				for floor := 0; floor<con.N_FLOORS; floor++{
					for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
						if NewUpdateLocalElevator.Queue[floor][button] {
							temp_ButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
							costID := costCalculator(thisElevator, elevatorList, temp_ButtonEvent, onlineElevators)
							temp_Keypress = con.Keypress{Floor: floor, Button:button, DesignatedElevator: costID}
							SyncChannels.LocalOrderToExternal <- temp_Keypress
						}
					}
				}
			}
			elevatorList[thisElevator] = NewUpdateLocalElevator
			go func() { Lights <- elevatorList }()
			if onlineElevators[thisElevator] {
				go func() { SyncChannels.LocalElevatorToExternal <- elevatorList}()
			}
		case tempElevatorArray := <-SyncChannels.UpdateMainLogic:
			change := false
			tempQueue := elevatorList[thisElevator].Queue
			for id := 0; id < con.N_ELEVS; id++ {
				if id != thisElevator {
					for floor := 0; floor < con.N_FLOORS; floor++ {
						for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
							if elevatorList[thisElevator].Queue[floor][button] && !NewUpdateLocalElevator.Queue[floor][button] {
								change = true
							}
						}
						if change {
							change = false
							for newID := 0; newID < con.N_ELEVS; newID++ {
								if newID == thisElevator {
									elevatorList[id].Queue[floor][elevio.BT_HallUp] = false
									elevatorList[id].Queue[floor][elevio.BT_HallDown] = false
								}
								if newID != id && newID != thisElevator {
									tempElevatorArray[newID].Queue[floor][elevio.BT_HallUp] = false
									tempElevatorArray[newID].Queue[floor][elevio.BT_HallDown] = false
								}
							}
						}
					}
				}
			}
			if tempQueue != elevatorList[thisElevator].Queue {
				elevatorList[thisElevator].Queue = tempQueue
				go func () { localStateChannel.DeleteQueue <- elevatorList[thisElevator].Queue} ()
				if onlineElevators[thisElevator] {
					go func () {SyncChannels.LocalElevatorToExternal <- elevatorList }()
				}
			}
			for id := 0; id < con.N_ELEVS; id++ {
				if id == thisElevator {
					continue
				}
				elevatorList[id] = tempElevatorArray[id]
			}
			go func() { UpdateLight <- elevatorList }() 
		case updatedOnlineElevators := <- SyncChannels.onlineElevators:
			change := false
			N_ELEVS_ONLINE := 0

			for id := 0; id < con.N_ELEVS; id++ {
				if updatedOnlineElevators[id] == true {
					N_ELEVS_ONLINE++
				}
			}
			if N_ELEVS_ONLINE == 0 {
				for id := 0; id < con.N_ELEVS; id++ {
					for floor := 0; floor < con.N_FLOORS; floor++ {
						for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
							if id != thisElevator {
								change = true
								elevatorList[id].Queue[floor][button] = false
							}
						}
						elevatorList[thisElevator].Queue[floor][elevio.BT_HallUp] = false
						elevatorList[thisElevator].Queue[floor][elevio.BT_HallDown] = false
					}
				}
				if change {
					localStateChannel.DeleteQueue <- elevatorList[thisElevator].Queue
					go func() { SyncChannels.LocalElevatorToExternal <- elevatorList }()
				}
			}
			change = false
			if N_ELEVS_ONLINE > 0 {
				for id := 0; id < con.N_ELEVS; id++ {
					if onlineElevators[id] && !updatedOnlineElevators[id] {
						for floor := 0; floor < con.N_FLOORS; floor++ {
							for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
								if elevatorList[id].Queue[floor][button] {
									change = true
									elevatorList[id].Queue[floor][button] = false
									if button != elevio.BT_Cab {
										temp_ButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
										costID := costCalculator(thisElevator, elevatorList, temp_ButtonEvent, updatedOnlineElevators)
										if costID == thisElevator {
											localStateChannel.NewOrder <- temp_ButtonEvent
										}
									}
								}
							}
						}
					}
				}
			}
			if change {
				go func() { SyncChannels.LocalElevatorToExternal <- elevatorList }()
			}
			go func() { UpdateLight <- elevatorList }()
			onlineElevators = updatedOnlineElevators
		}
	}
}

func LightSetter(UpdateLight chan [con.N_ELEVS]con.Elev, thisElevator int) {
	var Order [con.N_ELEVS]bool
	for {
		select {
		case Elevator := <- UpdateLight:
			for floor := 0; floor  < con.N_FLOORS; floor++ {
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
					for id := 0; id < con.N_ELEVS; id++ {
						Order[id] = false
						if id != thisElevator && button == elevio.BT_Cab {
							continue
						}
						if Elevator[id].Queue[floor][button] {
							elevio.SetButtonLamp(button, floor, true)
							Order[id] = true
						}
					}
					if Order == [con.N_ELEVS]bool{false} {
						elevio.SetButtonLamp(button, floor, false)
					}
				}
			}
		}
	}
}

