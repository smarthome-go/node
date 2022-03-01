#!/bin/bash
echo 'Pushing to Target...'
scp install.sh pi_room:~/
echo 'Success, executing'
clear
ssh pi_room 'chmod +x install.sh && bash install.sh'
