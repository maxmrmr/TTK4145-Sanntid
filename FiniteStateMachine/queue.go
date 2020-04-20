package FiniteStateMachine

import ( 
	con "../Configurations"
	elevio"../Hardware"
) 



func queue_order_below(elevator con.Elev) bool {
	for floor := 0; floor < elevator.Floor; floor ++ {
	  	for button := elevio.BT_HallUp; button<=elevio.BT_Cab; button++ {
			if(elevator.Queue[floor][button] == 1){
		  		return true; 
			}  
	  }  
	}
	return false; 
}

func queue_order_above(elevator con.Elev) bool{
	for floor := elevator.Floor; floor < con.N_FLOORS; floor ++ {
	  	for button := elevio.BT_HallUp; button<=elevio.BT_Cab; button++ {
			if(elevator.Queue[floor][button] == 1){
		  		return true;
			}
	  	}
	}
	return false; 
}

func set_queue_dir(elevator con.Elev) elevio.MotorDirection{
    // burde ikke trengs, men ok
	if (floor==-1 && check_queue_empty(elevator)==0){
		return elevio.MD_Stop; //stop
	}
    
	if(elevator.Floor == con.N_FLOORS-1){
	  return elevio.MD_Down; //return down if queue is not empty and at top floor
	}

	if (elevator.Dir == elevio.MD_Up) {
		if (queue_order_above(elevator)){
      		return elevio.MD_Up; //Up
    	} else if (queue_order_belove(elevator)){
      		return elevio.MD_Down; //Down
    	}
  	}
	if (elevator.Dir==elevio.MD_Down) {
    	if (queue_order_belove(elevator)) {
      		return elevio.MD_Down; //Down
    	} else if (queue_order_above(elevator)){
      		return elevio.MD_Up; //Up
    	}
	}
  	if (elevator.Dir==elevio.MD_Stop){
    	if (queue_order_above(elevator)){
      		return elevio.MD_Up;
    	} else if (queue_order_belove(elevator)){
      		return elevio.MD_Down;
    	}
  	}
  	return elevio.MD_Stop;
}

func queue_elev_run_stop(elevator con.Elev) bool{
	if (elevator.Dir == elevio.MD_Down){
	  	if (elevator.Queue[elevator.Floor][elevio.BT_HallDown] == 1 || elevator.Queue[elevator.Floor][elevio.BT_Cab] == 1){
			return true;
	  	} else if ((queue_order_belove(elevator)==0) && elevator.Queue[elevator.Floor][BT_HallUp]){
			return true;
		}
	} else if (elevator.Dir == elevio.MD_Up){
	  	if (elevator.Queue[elevator.Floor][elevio.BT_HallUp] == 1 || elevator.Queue[elevator.Floor][elevio.BT_Cab] == 1){ //test lagt til etter ||
			return true;
	  	} else if (elevator.Floor==2 && elevator.Queue[elevator.Floor][elevio.BT_HallDown]==1){
			return true;
		} else if ((queue_order_above(elevator)==0) && elevator.Queue[elevator.Floor][elevio.BT_HallDown]){
			return true;
		}
	return false;
  
	}
	return false;
}

func check_queue_empty(elevator con.Elev) bool {
	for floor := 0; floor < con.N_FLOORS; floor ++ {
		for button = BT_HallUp; button <= elevio.BT_Cab; button++ {
			if (elevator.Queue[floor][button]==1){ //then it is not empty
				return false;
			}
		}
	}
	return true
}

func queue_order_at_floor(elevator con.Elev){
	for button := BT_HallUp; button <= elevio.BT_Cab; button++ {
		if (elevator.Queue[elevator.Floor][button]==1){ //then it is not empty
			return true;
		}
	}
}

func remove_all_queue(elevator con.Elev) {
	for floor := 0; floor < con.N_FLOORS; floor++ {
		for button = BT_HallUp; button <= elevio.BT_Cab; button ++{
			elevator.Queue[floor][button]==0
			
		}
	}
	for floor := 0; floor < con.N_FLOORS; floor++ {
		for button = BT_HallUp; button <= elevio.BT_Cab; button++ {
			if (!((floor == 0 && button == elevio.BT_HallDown) || (floor == (con.N_FLOORS-1) && button == elevio.BT_HallUp))){
				elevio.SetButtonLamp(button, floor, 0)
			}
			
		}
	}
	return true
}