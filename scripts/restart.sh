#!/bin/bash

shExit()
{
if [ $1 -eq 1 ]; then
    printf "\nFailed!!!\n\n"
    exit 1
fi
}

printf "\nBegin Restart server,LISTEN port 8888 \n\n"

lsof -i:8888 | grep LISTEN | awk '{print $2}' | xargs kill -s SIGINT && go run ./main.go -env fat
shExit $?

printf "\nDone.\n\n"
