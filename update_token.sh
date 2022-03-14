#!/bin/bash

#  Usage:
# ./update_token.sh "old-token" "new-token"

curl --request PUT \
  --url "http://localhost:8081/token?token=$1" \
  --header 'Content-Type: application/json' \
  --data "{
	\"token\": \"$2\"
}"