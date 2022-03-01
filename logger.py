import time
import os
from typing import Optional, Union
from datetime import datetime

log_level = 0

log_levels = ["debug", "info", "warn", "error", "fatal"]
log_colors = {
    "debug": 14,
    "info": 2,
    "warn": 3,
    "error": 9,
    "fatal": 1
}

log_italic = {
    "debug": False,
    "info": False,
    "warn": False,
    "error": False,
    "fatal": True
}


def style(
        bold: Optional[bool] = None,
        italic: Optional[bool] = None,
        underline: Optional[bool] = None,
        blink: Optional[bool] = None,
        strike: Optional[bool] = None,
        fg: Union[int, tuple, None] = None,
        bg: Union[int, tuple, None] = None
):
    out = ''
    if bold is not None:
        out += '1;' if bold else '22;'
    if italic is not None:
        out += '3;' if italic else '23;'
    if underline is not None:
        out += '4;' if underline else '24;'
    if blink is not None:
        out += '5;' if blink else '25;'
    if strike is not None:
        out += '9;' if strike else '29;'
    if fg is not None:
        if isinstance(fg, int):
            out += f'38;5;{fg};'
        else:
            out += f'38;2;{fg[0]};{fg[1]};{fg[2]};'
    if bg is not None:
        if isinstance(bg, int):
            out += f'48;5;{bg};'
        else:
            out += f'48;2;{bg[0]};{bg[1]};{bg[2]};'
    return f'\033[{out[:-1]}m' if out else '\033[0m'


def write(text):
    text = str(text)
    with open("smarthome.log", "a") as file:
        file.write("\n"+text)


def log(text, level):
    text = str(text)
    try:
        if level < 0:
            pass
    except TypeError:
        log(f"The loglevel must be an integer in range 0, {len(log_levels)}.", 3)
        return False
    c_time = datetime.now().strftime("%H:%M:%S")
    text = text.split("\n")
    for line in text:
        if level <= len(log_levels) and level >= log_level:
            print(f"{style(fg=log_colors[log_levels[level]], italic=log_italic[log_levels[level]])}[{c_time}]  ({log_levels[level]+' '*(5-len(log_levels[level]))})  '{line}'{style()}")
        write(f"[{c_time}]  ({log_levels[level]+' '*(5-len(log_levels[level]))})  '{line}'")
    else:
        pass


def clear():
    with open("hardware.log", "w") as file:
        file.write("# Smarthome flask app log")

clear()