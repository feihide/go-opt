#! /bin/sh
project=$1
if [ "$project" = "" ]; then
    project="app"
fi
name=$2
if [ "$name" = "" ]; then
    name="main"
fi

ps -ef | grep ./$name"_mac" | grep -v grep | awk '{print $2}' | xargs  kill -9 
cd $project/bin
./$name"_mac" -port="8600"
