import RPi.GPIO as GPIO
from time import sleep
import requests
from datetime import datetime

PinIn = 26
GPIO.setmode(GPIO.BCM)
GPIO.setup(PinIn, GPIO.IN)


def ConvertHex(BinVal):
    tmpB2 = int(str(BinVal), 2)
    return hex(tmpB2)


def getData():  # Pulls data from sensor
    num1s = 0  # Number of consecutive 1s
    command = []  # Pulses and their timings
    binary = 1  # Decoded binary command
    previousValue = 0  # The previous pin state
    value = GPIO.input(PinIn)  # Current pin state
    while value:  # Waits until pin is pulled low
        sleep(0.0001) # This sleep decreases CPU utilization immensely
        value = GPIO.input(PinIn)
        
    startTime = datetime.now()  # Sets start time
    while True:
        if value != previousValue:  # Waits until change in state occurs
            now = datetime.now()  # Records the current time
            pulseLength = now - startTime  # Calculate time in between pulses
            startTime = now  # Resets the start time
            # Adds pulse time to array (previous val acts as an alternating 1 / 0 to show whether time is the on time or off time)
            command.append((previousValue, pulseLength.microseconds))
        # Interrupts code if an extended high period is detected (End Of Command)
        if value:
            num1s += 1
        else:
            num1s = 0
        if num1s > 10000:
            break
        # Reads values again
        previousValue = value
        value = GPIO.input(PinIn)
    # Covers data to binary
    for (typ, tme) in command:
        if typ == 1:
            if tme > 1000:  # According to NEC protocol a gap of 1687.5 microseconds represents a logical 1 so over 1000 should make a big enough distinction
                binary = binary * 10 + 1
            else:
                binary *= 10
    if len(str(binary)) > 34:  # Sometimes the binary has two rouge characters on the end
        binary = int(str(binary)[:34])
    return binary


def runTest():  # Actually runs the test
    # Takes samples
    data = getData()
    command = ConvertHex(data)
    print("Hex value: " + str(command))  # Shows results on the screen
    return command


# Main program loop
isOn = False
numMap = {"0x1001": "n4", "0x801": "s2", "0x1002": "s3", "0x802": "s5", "0x300f740bf": "s4",
          "0x300f7d02f": "n4", "0x300f7f00f": "s2", "0x300f7c837": "s3", "0x300f7e817": "s5"}
powerMap = {"n4": False, "s2": False, "s3": False, "s5": False}
while True:
    finalData = runTest()
    if str(finalData) in numMap:
        try:
            powerMap[numMap[str(finalData)]] = not powerMap[numMap[str(finalData)]]
            print("a code")
            url = 'http://192.168.178.111:8123/api/power/set?username=mik&password=test'
            body = {"switch": numMap[str(
                finalData)], "powerOn": powerMap[numMap[str(finalData)]]}
            x = requests.post(url, json=body)
            print(x.text)
            print(powerMap)
        except:
            pass
GPIO.cleanup()
