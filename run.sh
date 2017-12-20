#! /bin/sh
project=$1
if [ "$project" = "" ]; then
    project="app"
fi
name=$2
if [ "$name" = "" ]; then
    name="main"
fi

ps -ef | grep ./$name | grep -v grep | awk '{print $2}' | xargs -r  kill -9 
cd $project/bin
./$name
