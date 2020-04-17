package FiniteStateMachine

import (
	"time"
)

func timer(timeout chan<- bool, reset <-chan bool) {
	const time_dooropen = 3*time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <- reset:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			timeout <- True
		}
	}
}