package Configurations

import ( 
	elevio "../Hardware"
)

const (
	N_FLOORS = 4
	N_ELEVS = 3
	N_BUTTONS = 3
)

type Button int
//Button types for function elev_set_button_lamp() and elev_get_button()
const ( 
	BUTTON_CALL_UP Button = iota
	BUTTON_CALL_DOWN 
	BUTTON_COMMAND
)

type Motor_Direction int

const ( 
	DIRN_DOWN Motor_Direction = iota - 1
	DIRN_STOP
	DIRN_UP
)

type ElevState int

const (
	Undefined ElevState = iota - 1 // undefined = -1 should be in INIT
	IDLE
	RUN
	STOPPER
	DOOROPEN
)

type Elev struct {
	State ElevState
	Dir elevio.MotorDirection
	Floor int
	Queue [N_FLOORS][N_BUTTONS]bool
}

type Keypress struct {
	Button elevio.ButtonType
	Floor int
	DesignatedElev int
}

const (
	Connected int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)