#!/bin/bash
# this script is used to populate the database with the user defined in ./init/init.sql

cat ./init.sql | docker exec -i satd_db psql -U sleepy -d satd
