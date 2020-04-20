package FiniteStateMachine

import (
	"fmt"
	"time"

	"../Configurations"
	"../Hardware"
)

int curr_floor;
int prev_floor;
int dir;
elev_button_type_t button;

states_t state = INIT;

// StateMachineChannels contains the channels between the elevators
type StateMachineChannels struct {
	OrderCompleted  chan int
	Elecator 		chan Elev
	StateError 		chan error
	NewOrder        chan Keypress
	ArrivedAtFloor  chan int
}

type states_t enum {
	INIT,
	IDLE,
	RUN, 
	STOPPER,
	DOOROPEN
}

// RunFSM runs elevator and updates variables in stateMachineChannels

func RunFSM(ch StateMachineChannels){
	elevator := Elev{
		State: IDLE, 
		Dir: DIRN_STOP,
		Floor: hardware.elev_get_floor_sensor_signal()
		Queue: [N_FLOOR][N_BUTTONS]bool{},
	}
	// Need a timer to check if order is completed
	
	update_queue();
	curr_floor = elev_get_floor_sensor_signal(); //gjør om til en variabel, øker lesbarhet
	// sett lys etter hvilken etasje heisen er i
	if (curr_floor ! = -1){
		prev_floor = curr_floor;
		elev_set_floor_indicator(curr_floor);
	}

	//stopper og gjelder uansett
	if(elev_get_stop_signal() == 1) {
		state = STOPPER;
	  }
	
		switch (state){
	
		  case INIT:
			elev_set_motor_direction(DIRN_UP);
			dir = 1;
	
			if(curr_floor != -1){
				elev_set_motor_direction(DIRN_STOP);
				state = IDLE;
				break;
			}
	
			break;
	
		  case IDLE:
			if (check_queue_empty() == 0){
			  state = RUN;
			  break;
			}
			break;
	
		  case RUN:
			//sjekker etasjene heisen passerer
			if(elev_get_floor_sensor_signal() != -1){
			  if (queue_elev_run_stop(curr_floor, dir) == 1) { //potensiell test
							remove_floor_from_queue(curr_floor);
							  remove_light(curr_floor);
				state = DOOROPEN;
				start_timer();
				break;
			  }
						if (state!=DOOROPEN){
							dir = set_queue_dir(curr_floor,dir);
			}
				}
			break;
	
		  case STOPPER:
			//queue cleares, heisen stopper
			elev_set_stop_lamp(1);
			//int dirn = elev_get_motor_directcheck_queue_empty()==0){ion(); //setter dir til den retningen vi kjørte i før stop ble trykekt
			elev_set_motor_direction(DIRN_STOP);
			remove_all_queue();
	
	
			if ((elev_get_floor_sensor_signal()!= -1 ) && elev_get_stop_signal()){
			  while ((elev_get_floor_sensor_signal()!= -1 ) && elev_get_stop_signal()){
				elev_set_door_open_lamp(1);
			  }
			  state = DOOROPEN;
			  elev_set_stop_lamp(0);
			  break;
			  }
	
			if (!elev_get_stop_signal()){
				elev_set_stop_lamp(0);
				state=INIT;
				break;
				if(elev_get_floor_sensor_signal()==-1){
				  elev_set_motor_direction(dir); //fortsetter å kjøre i samme retning som før stopp
				}
				state = DOOROPEN;
				 break;
			  }
	
	
			break;
	
		  case DOOROPEN: //door open
	
			elev_set_motor_direction(DIRN_STOP);
	
			elev_set_door_open_lamp(1);
			if(check_timer()){
			  elev_set_door_open_lamp(0);
			  if (!check_queue_empty()){
				state = RUN;
				break;
			  } else
				state = IDLE;
				break;
			  }
			}
}




