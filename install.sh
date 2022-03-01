#!/bin/bash
VERSION=V.1.0
echo -e 'Installing Smarthome Hardware Gateway \e[34m ' ${VERSION} '\e[0m ...'

sudo apt update
sudo apt upgrade

sudo apt install apache2 python3 python3-pip screen -y
echo 'Installed apache2, python3 and pip3'
sudo a2enmod proxy proxy_http
echo 'Enabled proxy and proxy_http'

pip3 install flask rpi-rf Rpi.GPIO
echo 'Installed flask, rpi-rf and Rpi.GPIO'

echo 'Attempting to configure apache2'

echo '<VirtualHost *:80>
          ProxyPass / http://127.0.0.1:4243/
          ProxyPassReverse / http://127.0.0.1:4243/
          ProxyPass /error/ !
      </VirtualHost>' | sudo tee '/etc/apache2/sites-enabled/000-default.conf'

sudo a2ensite 000-default.conf
sudo systemctl reload apache2
sudo apt autoremove
sudo apt autoclean
