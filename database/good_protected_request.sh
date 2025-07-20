#!/bin/bash

if [ $# -eq 0 ]; then
  echo "Please provide a JWT token so I can insert it into the Bearer"
  exit 1
fi

curl -k -X GET https://localhost:3000/fetch-dashboard-info \
  -H "Authorization: Bearer $1"
