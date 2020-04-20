package Master

import ( 
	con "../Configurations"
	elevio "../Hardware"
)


func costCalculator(thisElevator int, elevatorList [con.N_ELEVS]con.Elev, LocalOrder elevio.ButtonEvent, onlineElevators [con.N_ELEVS]bool) int {
	if LocalOrder.Button == elevio.BT_Cab {
		return thisElevator
	}
	
	var CostList [con.N_ELEVS]int

	for elev :=0; elev < con.N_ELEVS; elev++ {
		cost := LocalOrder.Floor - elevatorList[elev].Floor

		if cost == 0 && onlineElevators[elev] && elevatorList[elev] != con.Undefined && elevatorList[elev].State != con.RUN  {
			return elev
		}
		if cost == 0 && elevatorList[elev].State == con.RUN {
			cost += 4
		}
		if cost < 0 {
			cost = -cost
			if elevatorList[elev].Dir == elevio.MD_Up {
				cost += 3
			}
		} else if cost > 0 {
			if elevatorList[elev].Dir == elevio.MD_Down {
				cost += 3
			}
		}
		if elevatorList[elev].State == con.DOOROPEN {
			cost ++
		}
		CostList[elev] = cost
	}
	maxCost := 700
	bestElevator := -1
	for elev :=  0; elev < con.N_ELEVS; elev++ {
		if onlineElevators[elev] && elevatorList[elev].State != con.Undefined && CostList[elev] < maxCost {
			bestElevator = elev
			maxCost = CostList[elev]
		}
	}
	return bestElevator
}