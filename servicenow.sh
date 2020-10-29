#!/bin/bash

ENV_FILE=".env"

if [ -f $ENV_FILE ]; then
  # shellcheck disable=SC2059
  printf "File $ENV_FILE found! Starting.\n"
  # shellcheck disable=SC2046
  export $(cat /app/.env | sed 's/#.*//g')
  go get -u github.com/chromedp/chromedp
  # shellcheck disable=SC1097
  go run /app/servicenow-instance-wakeup.go -headless=false -username="${USERNAME}" -password="${PASSWORD}" -debug=true
else
  # shellcheck disable=SC2059
  printf "File $ENV_FILE not found!\n"
fi
