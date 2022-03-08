#!/bin/bash
echo 'Pushing to Target...'
scp readIR.py pi_box:~/
echo 'Success, executing'
clear
ssh pi_box 'python3 readIR.py'

