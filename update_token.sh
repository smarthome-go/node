#!/bin/bash

echo -e "Enter \x1b[1;31mcurrent\x1b[1;0m token:"
read -rs OLD_TOKEN

echo -e "Enter \x1b[1;32mnew\x1b[1;0m token:"
read -rs NEW_TOKEN

curl --request PUT \
  --url "http://localhost:8081/token?token=$OLD_TOKEN" \
  --header 'Content-Type: application/json' \
  --data "{
	\"token\": \"$NEW_TOKEN\"
}"