#!/bin/bash


#Get the dependant package

go get gopkg.in/natefinch/lumberjack.v2

# The go build command to build the files
# use ./build.sh linx to build linux target
# use ./build.sh mac to build mac target
# be default ./build.sh build linux target

case "$1" in 
    "linux")
    env GOOS=linux GOARCH=amd64 go build -o main main.go preprocessor.go handleclient.go ;;

    "mac")

    env GOOS=darwin GOARCH=amd64 go build -o main main.go preprocessor.go handleclient.go ;;

esac
