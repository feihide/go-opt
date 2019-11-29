! /bin/sh
#nohup ./run.sh > log  2>&1 &
#tail -f log
project=$1
if [ "$project" = "" ]; then
    project="app"
fi
name=$2
if [ "$name" = "" ]; then
    name="main"
fi
cd /work/opt/app
svn up bin
svn up template

ps -ef | grep ./$name | grep -v grep | awk '{print $2}' | xargs -r  kill -9
