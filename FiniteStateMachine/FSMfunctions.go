package FiniteStateMachine
import (
	"fmt"

	"../Configurations"
) 

//FIXME: lag Elev variabel med curr_floor, queue[floor][button], floor, dir




func queue_order_below(elevator Elev) bool{
	for (floor := 0; floor < elevator.curr_floor; floor ++){
	  for (button := 0; button<=BUTTON_COMMAND; button++){
		if(elevator.queue[floor][button] == 1){
		  return true; 
			  }  
	  }  
	}
	  return false; 
}

func queue_order_above(elevator Elev) bool{
	for (floor := elevator.curr_floor; floor < N_FLOORS; floor ++){
	  for (button := 0; button<=BUTTON_COMMAND; button++){
		if(elevator.queue[floor][button] == 1){
		  return true;
			  }
	  }
	}
	return false; 
}

func set_queue_dir(elevator Elev) Direction{
    // burde ikke trengs, men ok
	if (floor==-1 && check_queue_empty()==0){
		return DIRN_STOP;
	}
    // burde ikke trengs, men ok
	if(elevator.dir == 0){
	  elev_set_motor_direction(DIRN_UP); // opp
	  return DIRN_UP;
	} else if(elevator.floor == N_FLOORS-1){
	  elev_set_motor_direction(DIRN_DOWN); // down
	  return DIRN_DOWN;
	}

	if (elevator.dir == DIRN_UP){
		if (queue_order_above(elevator.floor)){
      elev_set_motor_direction(DIRN_UP);
      return DIRN_UP;
    } else if (queue_order_belove(floor)){
      elev_set_motor_direction(DIRN_DOWN);
      return DIRN_DOWN;
    }
  }
	if (elevator.dir==DIRN_DOWN){
    if (queue_order_belove(elevator.floor)){
      elev_set_motor_direction(DIRN_DOWN);
      return DIRN_DOWN;
    } else if (queue_order_above(elevator.floor)){
      elev_set_motor_direction(DIRN_UP);
      return DIRN_UP;
    }
	}
  if (elevator.dir==DIRN_STOP){
    if (queue_order_above(elevator.floor)){
      return DIRN_UP;
    } else if (queue_order_belove(elevator.floor)){
      return DIRN_DOWN;
    }
  }
  return DIRN_STOP;
}

func queue_elev_run_stop(elevator Elev) int{
	if (elevator.dir==DIRN_DOWN){
	  if (elevator.queue[elevator.floor][BUTTON_CALL_DOWN] == 1 || elevator.queue[elevator.floor][BUTTON_COMMAND] == 1){
		return 1;
	  } else if ((queue_order_belove(elevator.floor)==0) && queue[elevator.floor][BUTTON_CALL_UP]){
			  return 1;
		  }
	} else if (elevator.dir ==DIRN_UP){
	  if (elevator.queue[elevator.floor][BUTTON_CALL_UP] == 1 || elevator.queue[elevator.floor][BUTTON_COMMAND] == 1){ //test lagt til etter ||
		return 1;
	  } else if (elevator.floor==2 && elevator.queue[elevator.floor][BUTTON_CALL_DOWN]==1){
			  return 1;
		  } else if ((queue_order_above(elevator.floor)==0) && elevator.queue[elevator.floor][BUTTON_CALL_DOWN]){
			  return 1;
		}
	return 0;
  
	  }
	  return 0;
  }