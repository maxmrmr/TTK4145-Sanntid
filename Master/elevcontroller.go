package Master

import (
	"fmt"

	hw "../Hardware"
	con "../Configurations"
)

func elevator_controller(thisElevator, SyncChannels nc.NetworkChannels, localStateChannel fsm.StateChannels, elevatorcontrollers <- chan elevio.ButtonEvent, Lights chan <- [con.N_ELEVS]con.Elev) {

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
					temp_Keypress = con.Keypress{DesignatedElevtor: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button}
					SyncChannels.LocalOrderToExternal <- temp_Keypress
				}
			}
		case temp_Keypress = <- SyncChannels.ExternalOrderToLocal:
			temp_ButtonEvent = elevio.ButtonEvent{Button: temp_Keypress.Button, Floor: temp_Keypress.Floor}
			if elevatorList[thisElevator].State == con.Undefined {
				costID := costCalculator(thisElevator, temp_ButtonEvent, elevatorList, onlineElevators)
				temp_Keypress.DesignatedElevtor = costID
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

		}
		change = false
		if elevatorList[thisElevator].State != con.Undefined && NewUpdateLocalElevator.State == con.Undefined {
			elevatorList[thisElevator].State = con.Undefined
			for floor := 0; floor<con.N_FLOORS; floor++{
				for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
					if NewUpdateLocalElevator.Queue[floor][button] {
						temp_ButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
						costID := costCalculator(thisElevator, temp_ButtonEvent, elevatorList, onlineElevators)
						temp_Keypress = con.Keypress{Floor: floor, Button:button, DesignatedElevator: costID}
						SyncChannels.LocalOrderToExternal <- temp_Keypress
					}
				}
			}
		}
		elevatorList[thisElevator]
	}

}

