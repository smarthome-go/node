from flask import Flask, request, redirect
# from rpi_rf import RFDevice
# import RPi.GPIO as GPIO
from time import sleep
import logger
import json

# flask names
app = Flask(__name__)

debug = False
hardware_enabled = False
ip =   ''
port = ''
device_pin = 0
device_repeat = 0
config_file = {}
pins = {}
relay_codes = {}
outlet_codes = {}


def gpio_setup():
    GPIO.setmode(GPIO.BCM)
    for i in pins:
        GPIO.setup(pins[i], GPIO.OUT)
        GPIO.output(pins[i], GPIO.HIGH)


# toggles a certain gpio pin
def pc_gpio_handler(pin):
    if not hardware_enabled:
        logger.log('Hardware is disabled', 2)
        return
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(pin, GPIO.OUT)
    GPIO.output(pin, GPIO.LOW)
    sleep(0.5)
    GPIO.output(pin, GPIO.HIGH)
    GPIO.cleanup(pin)


def get_outlet_code(outlet, state, outlets):
    for i in outlets:
        if i == outlet:
            if state:
                return outlets[i]['on']
            else:
                return outlets[i]['off']
    return


# Handles 433mhz outlets
def handle_outlet(p_outlet, p_state, p_device_pin, repeat, p_outlets):
    if p_outlet not in p_outlets:
        logger.log('this outlet does not exist', 3)
        return
    if hardware_enabled:
        rf_device = RFDevice(p_device_pin)
        rf_device.enable_tx()
        rf_device.tx_repeat = repeat
        code = get_outlet_code(p_outlet, p_state, p_outlets)
        rf_device.tx_code(code, 1, 180, 24)
        rf_device.cleanup()
    else:
        logger.log(get_outlet_code(p_outlet, p_state, p_outlets), 1)
        logger.log(f"handled outlet {p_outlet} with {p_state}", 2)


@app.route('/', methods=["POST"])
def request_handler():
    if request.method == "POST":
        req = {}
        if request.json:
            req = request.json
        if debug:
            logger.log(req, 0)
        if config_file['security']['useToken']:
            if 'token' in req and req['token'] == config_file['security']['token']:
                pass
            else:
                return 'unauthorized', 401

        logger.log(req, 1)
        if req["content"] == "power" and 'switch' in req and 'turnOn' in req:
            if req['switch'] in outlet_codes and req['switch'][0] != 'r':
                handle_outlet(req['switch'], req['turnOn'], device_pin, device_repeat, outlet_codes)
            elif req['switch'][0] == 'r':
                if (req['switch'] in relay_codes):
                    pc_gpio_handler(relay_codes[req['switch']])
                else:
                    return 'not found.', 404

        else:
            logger.log(req, 0)
            return 'not found.', 404
    return 'success', 200


def load_json():
    global debug, ip, port, hardware_enabled, pins, config_file, outlet_codes, device_pin, device_repeat, relay_codes, rgb_enabled
    with open("config/config.json", "r") as file:
        config_file = json.load(file)
        debug = config_file['server']['debug']
        ip = config_file['network']['ip']
        port = config_file['network']['port']
        hardware_enabled = config_file['hardware']['general']['enabled']
        outlet_codes = config_file['outletCodes']
        relay_codes = config_file['hardware']['relay']['pinOut']
        device_pin = config_file['hardware']['rfDevice']['devicePin']
        device_repeat = config_file['hardware']['rfDevice']['repeat']
        rgb_enabled = config_file['hardware']['rgbLed']['enabled']

        if hardware_enabled:
            from rpi_rf import RFDevice
            import RPi.GPIO as GPIO
            GPIO.setwarnings(config_file['hardware']['general']['gpioSetWarnings'])
            pins = config_file['hardware']['rgbLed']['pinOut']
            logger.log(pins, 0)
            gpio_setup()
        else:
            logger.log('HARDWARE DISABLED, ONLY USE FOR TESTING.', 3)


if __name__ == '__main__':
    load_json()
    app.run(host=ip, port=port, debug=debug)
