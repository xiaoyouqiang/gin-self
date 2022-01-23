#!/bin/bash

shExit()
{
if [ $1 -eq 1 ]; then
    printf "\nFailed!!!\n\n"
    exit 1
fi
}

printf "\nStart create file\n\n"
time go run -v ./cmd/mysql_cmd/main.go  -addr $1 -user $2 -pass $3 -name $4 -tables $5
shExit $?