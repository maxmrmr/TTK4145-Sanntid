package FiniteStateMachine

import (
	"fmt"
	t "time"
	con "../Configurations"
	elevio "../Hardware"
)



// StateMachineChannels contains the channels between the elevators
type StateMachineChannels struct {
	OrderCompleted  chan int
	Elevator 		chan con.Elev
	NewOrder 		chan elevio.ButtonEvent
	ArrivedAtFloor  chan int
	DeleteQueue     chan [con.N_FLOORS][con.N_BUTTONS]bool
}



// RunFSM runs elevator and updates variables in stateMachineChannels

func RunFSM(ch StateMachineChannels, thisElevator  int){
	elevator := con.Elev{
		State: con.IDLE, 
		Dir: elevio.MD_Stop,
		Floor: elevio.getFloor(),
	}
	
	// Need a timer to check if order is completed
	timer_DoorOpen := t.NewTimer(3 * t.Second)
	timer_ElevatorLost := t.NewTimer(3 * t.Second)

	timer_DoorOpen.Stop()
	timer_ElevatorLost.Stop()
	updateExternal := false

	//update_queue();
	//curr_floor = elev_get_floor_sensor_signal(); //gjør om til en variabel, øker lesbarhet
	// sett lys etter hvilken etasje heisen er i
	if (elevio.getFloor() == -1) {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	//readFromFile("elevOrders", thisElevator, &elevator)


	//stopper og gjelder uansett
	if ( elevio.getStop() == 1) { //ville ikke fungere med elevio...
		elevator.State = con.STOPPER;
	}
	for {
		select {
		case newOrder := <- ch.NewOrder:
			elevator.Queue[newOrder.Floor][newOrder.Button] = true
			switch (elevator.State) {
			case con.Undefined:
				fmt.Println("Error, elevator state for elevator running is undefines in fsm.go function\n")
				updateExternal = true
				break
			case con.IDLE:
				elevator.Dir = set_queue_dir(elevator)
				elevio.SetMotorDirection(elevator.Dir)
				if elevator.Dir == elevio.MD_Stop {
					elevator.State = con.DOOROPEN
					elevio.SetDoorOpenLamp(true)
					timer_DoorOpen.Reset(3* t.Second)
					elevator.Queue[elevator.Floor] = [con.N_BUTTONS]bool{false} //removing order from queue
				} else {
					elevator.State = con.RUN
					timer_ElevatorLost.Reset( 3* t.Second)
				}
				updateExternal = true

	
			case con.RUN:
				updateExternal = true
	
		  	case con.STOPPER:
				//queue cleares, heisen stopper
				elevio.SetStopLamp(true)
				//int dirn = elev_get_motor_directcheck_queue_empty()==0){ion(); //setter dir til den retningen vi kjørte i før stop ble trykekt
				elevio.SetMotorDirection(elevio.MD_Stop)
				remove_all_queue(elevator)
	
	
				if ((elevio.getFloor() !=-1) && elevio.getStop()) {
			  		for ((elevio.getFloor() != -1 ) && elevio.getStop()) {
						elevio.SetDoorOpenLamp(1)
			  		}
			  	elevator.State = con.DOOROPEN
			  	elevio.SetStopLamp(0)
			  	break
			  	}
	
				if (!elevio.GetStop()){
					elevio.SetStopLamp(0)
					elevator.state=con.Undefined
					break;
					if (elevio.GetFloor()==-1){
				  		elevio.SetMotorDirection(elevator.Dir) //fortsetter å kjøre i samme retning som før stopp
					}
					elevator.state = con.DOOROPEN
				 	break
			  	}	
				break
	
		  	case con.DOOROPEN: //door open
	
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(1)
				if elevator.Floor == newOrder.Floor {
					timer_DoorOpen.Reset(3 * t.Second)
					elevator.Queue[elevator.Floor] = [con.N_BUTTONS]bool{false}	
				} else {
					updateExternal = true
				}
				
			}
		case DeleteQueue := <-ch.DeleteQueue:
			elevator.Queue = DeleteQueue
		case elevator.Floor = <- ch.ArrivedAtFloor:
			elevio.SetFloorIndicator(elevator.Floor)
			if queue_elev_run_stop(elevator) {
				timer_ElevatorLost.Stop()
				elevio.SetMotorDirection(elevio.MD_Stop)
				if !orderAtFloor(elevator) {
					elevator.State = con.IDLE
					timer_DoorOpen.Reset(3 * t.Second)
				} else {
					elevio.SetDoorOpenLamp(true)
					elevator.State = con.DOOROPEN
					DoorTimer.Reset(3 * t.Second)
					elevator.Queue[elevator.Floor] = [con.N_BUTTONS]bool{false}
				}
			} else if elevator.State == con.RUN {
				timer_ElevatorLost.Reset(3*t.Second)
			}
			updateExternal = true
		case <- timer_DoorOpen.C: //if dooropen is timed out, order is done
			elevio.SetDoorOpenLamp(false)
			elevator.Dir = set_queue_dir(elevator)
			if elevator.Dir == elevio.MD_Stop {
				elevator.State = con.IDLE
				timer_ElevatorLost.Stop()
			} else {
				elevator.State = con.RUN
				timer_ElevatorLost.Reset(3*t.Second)
				elevio.SetMotorDirection(elevator.Dir)
			}
			updateExternal = true
		case <- timer_ElevatorLost.C:
			elevator.State = con.Undefined
			fmt.Println("Elevator connection is lost")
			timer_ElevatorLost.Reset(5 * t.Second)
			updateExternal = true
		 }
		 if updateExternal {
			 updateExternal = false
			 go func() {ch.Elevator <- elevator }()
		}
	}
}

	
func UpdateKeysPressed(NewOrder chan con.Keypress, receiveOrder chan elevio.ButtonEvent) {
	var key config.Keypress
	key.DesignatedElevator = 1
	for {
		select {
		case order := <-receiveOrder:
			key.Floor = order.Floor
			key.Button = order.Button
			NewOrder <- key
		}
	}
}		  




