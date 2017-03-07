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
	stateChangeThreshold = 20
	cycleDelay           = 5 * time.Millisecond
	minimumDuration      = 2000 // in Millisecond
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

	lastReading = sensor.Read()

	var lastStateChange time.Time
	unchangedState := 0

	for {

		rightNow := time.Now()
		currentReading = sensor.Read()

		if lastReading != currentReading {

			// the pin state has changed!
			unchangedState = 0
			if currentState == stateStill {
				currentState = stateMoving
				lastStateChange = rightNow
				// fmt.Printf("Started Moving: %s\n", rightNow.Format(time.Stamp))
			}

		} else {

			unchangedState++

			if unchangedState > stateChangeThreshold {
				// we've been in the same same condition long enough to assume we're not moving

				if currentState == stateMoving {
					// we were moving but we've stopped now
					currentState = stateStill

					duration := rightNow.Sub(lastStateChange).Nanoseconds() / 1000000

					fmt.Printf("Vibration Ended After %d ms\n", duration)
					lastStateChange = rightNow

					if duration > minimumDuration {
						fmt.Println("LONG VIBRATION DETECTED")
					}

				}

				unchangedState = 0
			}
		}

		lastReading = currentReading
		time.Sleep(cycleDelay)
	}

}
