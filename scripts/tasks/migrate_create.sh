#!/bin/bash

mkdir -p ./db/migrations/
migrate create -ext ".sql"  -dir "./db/migrations/" -seq init
