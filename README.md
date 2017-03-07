# Vibration Sensor

## Hardware
For this project, I used a Qungi SW-420 Motion Sensor Module http://amzn.to/2mCcAc8
and a Raspberry Pi 3.

I connected the motion sensor to the +5V, GND and pin 23.

## Software

My approach is to detect bursts of activity on the data pin and try to determine
how long the vibration lasted.  This might require fine tuning of the sleep time between
samples or the number of samples that indicates a return to stillness.
