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
			id := costCalculator(thisElevator, newLocalOrder, elevatorList, onlineElevators)

		}
	}

}

