package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	sensorPin                                      rpio.Pin
	lastReading, currentReading                    rpio.State
	lastState, currentState                        int
	led1, led2, led3, led4, led5, led6, led7, led8 rpio.Pin
	pinSlots                                       []rpio.Pin
)

const (
	pinSense    = 23
	stateStill  = 0
	stateMoving = 1

	pinLED1 = 14
	pinLED2 = 15
	pinLED3 = 16
	pinLED4 = 17
	pinLED5 = 18
	pinLED6 = 19
	pinLED7 = 20
	pinLED8 = 21

	maxLit = 8

	// how many cycles do we go through with the same state to determine
	// that the sensor is no longer vibrating
	stateChangeThreshold = 20
	cycleDelay           = 5 * time.Millisecond
	minimumDuration      = 1000 // in Millisecond
)

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	fmt.Printf("%s: Reading Sensor on Pin %d\n", time.Now().Format(time.Stamp), pinSense)

	setup()

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
				// blink(led8, 1)
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
						renderLED(8)
						time.Sleep(100 * time.Millisecond)
						renderLED(0)

					} else {
						renderLED(3)
						time.Sleep(100 * time.Millisecond)
						renderLED(0)
					}

				}

				unchangedState = 0
			}
		}

		lastReading = currentReading
		time.Sleep(cycleDelay)
	}

}

func setup() {

	delay := 50 * time.Millisecond

	led1 = rpio.Pin(pinLED1)
	led2 = rpio.Pin(pinLED2)
	led3 = rpio.Pin(pinLED3)
	led4 = rpio.Pin(pinLED4)
	led5 = rpio.Pin(pinLED5)
	led6 = rpio.Pin(pinLED6)
	led7 = rpio.Pin(pinLED7)
	led8 = rpio.Pin(pinLED8)

	pinSlots = []rpio.Pin{led1, led2, led3, led4, led5, led6, led7, led8}
	led1.Output()
	led2.Output()
	led3.Output()
	led4.Output()
	led5.Output()
	led6.Output()
	led7.Output()
	led8.Output()

	renderLED(0)
	time.Sleep(delay)
	renderLED(1)
	time.Sleep(delay)
	renderLED(2)
	time.Sleep(delay)
	renderLED(3)
	time.Sleep(delay)
	renderLED(4)
	time.Sleep(delay)
	renderLED(5)
	time.Sleep(delay)
	renderLED(6)
	time.Sleep(delay)
	renderLED(7)
	time.Sleep(delay)
	renderLED(8)
	time.Sleep(250 * time.Millisecond)
	renderLED(0)
}

func renderLED(current int) {

	if current > maxLit {
		current = maxLit
	}

	for i := 0; i < maxLit; i++ {

		thePin := pinSlots[i]
		thePin.Low()

		if i < current {
			thePin.High()
		}
	}

}

func blink(pin rpio.Pin, times int) {

	for i := 0; i < times; i++ {
		pin.Low()
		time.Sleep(20 * time.Millisecond)
		pin.High()
		time.Sleep(20 * time.Millisecond)
		pin.Low()
	}
}
