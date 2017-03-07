package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	sensorPin                   rpio.Pin
	lastReading, currentReading rpio.State
	lastState, currentState     int
)

const (
	pinSense    = 23
	stateStill  = 0
	stateMoving = 1

	// how many cycles do we go through with the same state to determine
	// that the sensor is no longer vibrating
	stateChangeThreshold = 50
	cycleDelay           = 20 * time.Millisecond
)

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	fmt.Printf("%s: Reading Sensor on Pin %d\n", time.Now().Format(time.Stamp), pinSense)

	sensor := rpio.Pin(pinSense)
	sensor.Input()

	currentState = int(sensor.Read())

	unchangedState := 0

	for {

		currentReading = sensor.Read()

		if lastReading == currentReading {
			unchangedState++

			if unchangedState > stateChangeThreshold && currentState != stateStill {
				currentState = stateStill
				fmt.Println("State Changed to Still")
				unchangedState = 0

			}
		}

		if lastReading != currentReading {
			// fmt.Printf("Sensor Change: %#v\n", currentReading)

			if currentState == stateStill && currentReading == rpio.High {
				currentState = stateMoving
				fmt.Println("Vibration Started")
			}

		}

		lastReading = currentReading
		time.Sleep(cycleDelay)
	}

}
