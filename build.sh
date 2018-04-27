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
cd $project/src
go build -o  ../bin/$name  .

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../bin/$name"_mac" .

ps -ef | grep ./$name | grep -v grep | awk '{print $2}' | xargs -r kill -9 
cd ../bin
./$name
