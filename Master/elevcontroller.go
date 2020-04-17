package ElevController

import (
	"fmt"
	. "github.com/maxmrmr/TTK4145-Sanntid/blob/master/Configurations"
	hw "https://github.com/maxmrmr/TTK4145-Sanntid/blob/master/Hardware"
)

func costCalculator(order buttonPressed, elevatorList[N_ELEVS], thisElevator int, onlineElevators [N_ELEVS]bool) int {
	minCost := (N_BUTTONS * N_FLOORS)*N_ELEVS
	bestElevator := thisElevator
	for elevator :=  0; elevator < N_ELEVS; elevator++ {
		if !onlineElevators[elevtor] {
			continue //disregarding offline elevators
		}
		cost := order.Floor - elevatorList[elevator].Floor
	}
}