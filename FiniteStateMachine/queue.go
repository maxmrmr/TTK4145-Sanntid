package queue


import (
	def "Configurations"
	"fmt"
	"log"
	"time"
)


// 4 indikerer etasjer
// 3 indikerer type knapp som er trykket

type queue struct {
	matrix [def.N_FLOORS][def.N_BUTTONS]orderStatus
}

//orderStatus defines the status of an order
type orderStatus struct {
	active bool
	addr string `json:"-"`
	timer *time.Timer `json:"-"`

}
var inactive = orderStatus{active: false, addr: "", timer: nil}

var local queue
var remote queue

var updateLocal = make(chan bool)
var takeBackup = make(chan bool, 10)
var OrderTimeoutChan = make(chan def.Keypress)
var newOrder chan bool


func Init(newOrderTemp chan bool, outgoing_message chan deg.Message) {
	newOrder = newOrderTemp
	go update_queue()
	runBackup(outgoing_message)
	log.Println("Queue is initialized.")
}

func add_to_localQueue(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, "", nil})
	newOrder <- true
}

//add order to remote queue, start timer. If order times out  the order will be taken care of
func add_to_remoteQueue(floor, button int, addr string) {
	alreadyExist := RemoteOrderExist(floor, button)
	remote.setOrder(floor, button, orderStatus{true, addr, nil})
	if !alreadyExist {
		go remote.startTimer(floor, button)
	}
	updateLocal <- true
}


func RemoteOrderExist(floor int, button int) bool {
	return remote.isOrder(floor, button)
}

func LocalOrderExist(floor int, button int) bool {
	return local.isOrder(floor, button)
}


void update_queue(){
  for {
	  <-updateLocal
	  for int floor: = 0; floor < def.N_FLOORS; floor++) {

    		for (button: = 0; button <= def.N_FLOORS; button++) {
				  if remote.isOrder(floor,button) {
					  if button != def.BUTTON_COMMAND && remote.matrix[floor][button].addr == def.Laddr {
						  if !local.isOrder(floor, button){
							  local.setOrder(floor, button, orderStatus{true, "", nil})
							  newOrder <- true
						  }
					  }

				  }
        	}
      	}
    }
  
}

func (q *queue) set_queue_dir(floor int, dir int) int {
    // burde ikke trengs, men ok
	if (floor==-1 && q.check_queue_empty()==0){
		return def.DIRN_STOP;
	}
    // burde ikke trengs, men ok
	if(floor == 0){
	  elev_set_motor_direction(def.DIRN_UP); // opp
	  return def.DIRN_UP;
	} else if(floor == N_FLOORS-1){
	  elev_set_motor_direction(def.def.); // down
	  return def.def.;
	}

	if (dir == DIRN_UP){
		if (q.queue_order_above(floor)){
      elev_set_motor_direction(def.DIRN_UP);
      return def.DIRN_UP;
    } else if (q.queue_order_belove(floor)){
      elev_set_motor_direction(def.def.);
      return def.def.;
    }
  }
	if (dir==def.){
    if (q.queue_order_belove(floor)){
      elev_set_motor_direction(def.def.);
      return def.def.;
    } else if (q.queue_order_above(floor)){
      elev_set_motor_direction(def.DIRN_UP);
      return def.DIRN_UP;
    }
	}
  if (dir==DIRN_STOP){
    if (queue_order_above(floor)){
      return def.DIRN_UP;
    } else if (queue_order_belove(floor)){
      return def.DIRN_DOWN;
    }
  }
  def.CloseConnectionChan <- true
  def.Restart.Run()
  log.Printf("Set_Queue_Dir(): called with invalid direction %d, returning stop\n", dir)
  return def.DIRN_STOP;
}

int queue_order_above(int curr_floor){
  for (int floor = curr_floor; floor < N_FLOORS; floor ++){
    for (int button = 0; button<=BUTTON_COMMAND; button++){
      if(queue[floor][button] == 1){
        return 1;
			}
    }
  }
  return 0;
}

func (q *queue) queue_order_below(floor int) bool{
	for (floor := 0; floor < elevator.curr_floor; floor ++){
	  for (button := 0; button<=BUTTON_COMMAND; button++){
		if(elevator.queue[floor][button] == 1){
		  return true; 
			}  
	  }  
	}
	  return false; 
}

// sjekker at køen vår er tom
func (q *queue) check_queue_empty() bool{
  for (int floor := 0; floor < def.N_FLOORS; floor ++){
    for (int button := 0; button < def.N_BUTTONS; button ++){
      if(q.matrix[floor][button].active){ // da er den i så fall ikke tom
        return 0;
      }
    }
  }
  return 1;
}

// sjekker om vi skal stoppe på vei opp eller ned til ønsket etasje, hvis vi stopper fjernes bestillingen
int queue_elev_run_stop(int floor, int dir){
  if (dir==def.){
    if (queue[floor][BUTTON_CALL_DOWN] == 1 || queue[floor][BUTTON_COMMAND] == 1){
      return 1;
    } else if ((queue_order_belove(floor)==0) && queue[floor][BUTTON_CALL_UP]){
			return 1;
		}
  } else if (dir ==DIRN_UP){
    if (queue[floor][BUTTON_CALL_UP] == 1 || queue[floor][BUTTON_COMMAND] == 1){ //test lagt til etter ||
      return 1;
    } else if (floor==2 && queue[floor][BUTTON_CALL_DOWN]==1){
			return 1;
		} else if ((queue_order_above(floor)==0) && queue[floor][BUTTON_CALL_DOWN]){
			return 1;
  	}
  return 0;

	}
	return 0;
}

void remove_floor_from_queue(int floor){
	queue[floor][BUTTON_CALL_UP] = 0;
	queue[floor][BUTTON_COMMAND] = 0;
	queue[floor][BUTTON_CALL_DOWN]=0;
}

void remove_light(int floor){
	if (floor==-1){
		return;
	}
	if (floor<=N_FLOORS-2){ //-2 fordi tall 3 er 4. etasje, men vi vil ikke ha med 4. etasje
	elev_set_button_lamp(BUTTON_CALL_UP,floor,0);
	}
	if (floor>=1){
	elev_set_button_lamp(BUTTON_CALL_DOWN,floor,0);
	}
	elev_set_button_lamp(BUTTON_COMMAND,floor,0);
}

void remove_all_queue(){
  for (int floor = 0; floor < N_FLOORS; floor ++){
    for (elev_button_type_t button = BUTTON_CALL_UP; button < N_BUTTONS; button ++){
      queue[floor][button] = 0;
    }
  }
  for (int floor = 0; floor < N_FLOORS;  floor ++) {
    for (elev_button_type_t button = BUTTON_CALL_UP; button <= BUTTON_COMMAND; button++){
      if(!((floor == 0 && button == BUTTON_CALL_DOWN) || (floor == (N_FLOORS-1) && button == BUTTON_CALL_UP))){
        elev_set_button_lamp(button, floor, 0);
      }
    }
  }
}