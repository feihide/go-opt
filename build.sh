#! /bin/sh
project=$1
if [ "$project" = "" ]; then
    project="app"
fi
name=$2
if [ "$name" = "" ]; then
    name="main"
fi

gofmt -w -l  $project/src

go build -o $project/bin/$name  $project/src/$name.go

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $project/bin/$name"_mac"  $project/src/$name.go

ps -ef | grep ./$name | grep -v grep | awk '{print $2}' | xargs -r kill -9 
cd $project/bin
./$name
