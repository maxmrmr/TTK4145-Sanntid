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
		Floor: elevio.GetFloor(),
	}
	
	// Need a timer to check if order is completed
	timerDoorOpen := t.NewTimer(3 * t.Second)
	timerElevatorLost := t.NewTimer(3 * t.Second)

	timerDoorOpen.Stop()
	timerElevatorLost.Stop()
	updateExternal := false

	if (elevio.GetFloor() == -1) {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	for (elevio.GetFloor() == -1) {
		t.Sleep(10  * t.Millisecond)
	}
	elevio.SetMotorDirection(elevio.MD_Stop) 
	//readFromFile("elevOrders", thisElevator, &elevator)


	//stopper og gjelder uansett
	if ( elevio.GetStop() == true) { 
		elevator.State = con.STOPPER; // litt usikre på om det er slik funksjonaliteten skal være
	}
	for {
		select {
		case newOrder := <- ch.NewOrder:
			elevator.Queue[newOrder.Floor][newOrder.Button] = true
			switch (elevator.State) {
			case con.Undefined:
				fmt.Println("Error, elevator state for elevator running is undefines in fsm.go function")
				updateExternal = true
				break
			case con.IDLE:
				elevator.Dir = QueueSetDir(elevator)
				elevio.SetMotorDirection(elevator.Dir)
				if elevator.Dir == elevio.MD_Stop {
					elevator.State = con.DOOROPEN
					elevio.SetDoorOpenLamp(true)
					timerDoorOpen.Reset(3* t.Second)
					elevator.Queue[elevator.Floor] = [con.N_BUTTONS]bool{false} //removing order from queue
				} else {
					elevator.State = con.RUN
					timerElevatorLost.Reset( 3* t.Second)
				}
				updateExternal = true

	
			case con.RUN:
				updateExternal = true
	
		  	case con.STOPPER:
				//queue cleares, heisen stopper
				elevio.SetStopLamp(true)
				//int dirn = elev_get_motor_directcheck_queue_empty()==0){ion(); //setter dir til den retningen vi kjørte i før stop ble trykekt
				elevio.SetMotorDirection(elevio.MD_Stop)
				QueueRemoveAll(elevator)
	
	
				if ((elevio.GetFloor() !=-1) && elevio.GetStop()) {
			  		for ((elevio.GetFloor() != -1 ) && elevio.GetStop()) {
						elevio.SetDoorOpenLamp(true)
			  		}
			  	elevator.State = con.DOOROPEN
			  	elevio.SetStopLamp(false)
			  	break
			  	}
	
				if (!elevio.GetStop()){
					elevio.SetStopLamp(false)
					elevator.State = con.Undefined
					if (elevio.GetFloor()==-1){
				  		elevio.SetMotorDirection(elevator.Dir) //fortsetter å kjøre i samme retning som før stopp
					}
					elevator.State = con.DOOROPEN
				 	break
			  	}	
				break
	
		  	case con.DOOROPEN: //door open
	
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				if elevator.Floor == newOrder.Floor {
					timerDoorOpen.Reset(3 * t.Second)
					elevator.Queue[elevator.Floor] = [con.N_BUTTONS]bool{false}	
				} else {
					updateExternal = true
				}
				
			}
		case DeleteQueue := <-ch.DeleteQueue:
			elevator.Queue = DeleteQueue
		case elevator.Floor = <- ch.ArrivedAtFloor:
			elevio.SetFloorIndicator(elevator.Floor)
			if QueueElevRunStop(elevator) {
				timerElevatorLost.Stop()
				elevio.SetMotorDirection(elevio.MD_Stop)
				if !QueueOrderAtFloor(elevator) {
					elevator.State = con.IDLE
					timerDoorOpen.Reset(3 * t.Second)
				} else {
					elevio.SetDoorOpenLamp(true)
					elevator.State = con.DOOROPEN
					//FIXME: Det stod: DoorTimer.Reset(3 * t.Second). Endret til det som står under. Ok?
					timerDoorOpen.Reset(3 * t.Second)
					elevator.Queue[elevator.Floor] = [con.N_BUTTONS]bool{false}
				}
			} else if elevator.State == con.RUN {
				timerElevatorLost.Reset(3*t.Second)
			}
			updateExternal = true
		case <- timerDoorOpen.C: //if dooropen is timed out, order is done
			elevio.SetDoorOpenLamp(false)
			elevator.Dir = QueueSetDir(elevator)
			if elevator.Dir == elevio.MD_Stop {
				elevator.State = con.IDLE
				timerElevatorLost.Stop()
			} else {
				elevator.State = con.RUN
				timerElevatorLost.Reset(3*t.Second)
				elevio.SetMotorDirection(elevator.Dir)
			}
			updateExternal = true
		case <- timerElevatorLost.C:
			elevator.State = con.Undefined
			fmt.Println("Elevator connection is lost")
			timerElevatorLost.Reset(5 * t.Second)
			updateExternal = true
		 }
		 if updateExternal {
			 updateExternal = false
			 go func() {ch.Elevator <- elevator }()
		}
	}
}

	
func UpdateKeysPressed(NewOrder chan con.Keypress, receiveOrder chan elevio.ButtonEvent) {
	var key con.Keypress
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




