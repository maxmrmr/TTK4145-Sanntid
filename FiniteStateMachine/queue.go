package FiniteStateMachine

import ( 
	con "../Configurations"
	elevio "../Hardware"
) 



func queueOrderBelow(elevator con.Elev) bool {
	for floor := 0; floor < elevator.Floor; floor ++ {
	  	for button := elevio.BT_HallUp; button <=elevio.BT_Cab; button++ {
			if(elevator.Queue[floor][button] == true){
		  		return true; 
			}  
	  }  
	}
	return false; 
}

func queueOrderAbove(elevator con.Elev) bool{
	for floor := elevator.Floor; floor < con.N_FLOORS; floor ++ {
	  	for button := elevio.BT_HallUp; button<=elevio.BT_Cab; button++ {
			if(elevator.Queue[floor][button] == true){
		  		return true;
			}
	  	}
	}
	return false; 
}

func QueueSetDir(elevator con.Elev) elevio.MotorDirection{
    // burde ikke trengs, men ok
	if (floor==-1 && QueueCheckEmptyQueue(elevator)==false){
		return elevio.MD_Stop; //stop
	}
    
	if(elevator.Floor == con.N_FLOORS-1){
	  return elevio.MD_Down; //return down if queue is not empty and at top floor
	}

	if (elevator.Dir == elevio.MD_Up) {
		if (queueOrderAbove(elevator)){
      		return elevio.MD_Up; //Up
    	} else if (queueOrderBelow(elevator)){
      		return elevio.MD_Down; //Down
    	}
  	}
	if (elevator.Dir==elevio.MD_Down) {
    	if (queueOrderBelow(elevator)) {
      		return elevio.MD_Down; //Down
    	} else if (queueOrderAbove(elevator)){
      		return elevio.MD_Up; //Up
    	}
	}
  	if (elevator.Dir==elevio.MD_Stop){
    	if (queueOrderAbove(elevator)){
      		return elevio.MD_Up;
    	} else if (queueOrderBelow(elevator)){
      		return elevio.MD_Down;
    	}
  	}
  	return elevio.MD_Stop;
}

func QueueElevRunStop(elevator con.Elev) bool{
	if (elevator.Dir == elevio.MD_Down){
	  	if (elevator.Queue[elevator.Floor][elevio.BT_HallDown] == true || elevator.Queue[elevator.Floor][elevio.BT_Cab] == true){
			return true;
	  	} else if ((queueOrderBelow(elevator) == false) && elevator.Queue[elevator.Floor][elevio.BT_HallUp]){
			return true;
		}
	} else if (elevator.Dir == elevio.MD_Up){
	  	if (elevator.Queue[elevator.Floor][elevio.BT_HallUp] == true || elevator.Queue[elevator.Floor][elevio.BT_Cab] == true){ //test lagt til etter ||
			return true;
	  	} else if (elevator.Floor==2 && elevator.Queue[elevator.Floor][elevio.BT_HallDown] == true){
			return true;
		} else if ((queueOrderAbove(elevator)== false) && elevator.Queue[elevator.Floor][elevio.BT_HallDown]){
			return true;
		}
	return false;
  
	}
	return false;
}

func QueueCheckEmptyQueue(elevator con.Elev) bool {
	for floor := 0; floor < con.N_FLOORS; floor ++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
			if (elevator.Queue[floor][button] == true){ //then it is not empty
				return false;
			}
		}
	}
	return true
}

func QueueOrderAtFloor(elevator con.Elev) bool {
	for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
		if elevator.Queue[elevator.Floor][button] == true { //then it is not empty
			return true
		}
	}
	return false
}

func QueueRemoveAll(elevator con.Elev) {
	for floor := 0; floor < con.N_FLOORS; floor++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button ++{
			elevator.Queue[floor][button] == false
		}
	}
	for floor := 0; floor < con.N_FLOORS; floor++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
			if (!((floor == 0 && button == elevio.BT_HallDown) || (floor == (con.N_FLOORS-1) && button == elevio.BT_HallUp))){
				elevio.SetButtonLamp(button, floor, false)
			}
		}
	}
}