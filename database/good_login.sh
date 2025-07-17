#!/bin/bash

curl -v -k -X POST https://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin", "password":"awnoidroppeditinthewater"}'
