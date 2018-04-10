#! /bin/sh
while :
do
        if test `ps -ef |grep ./run.sh |grep -v grep |wc -l` -eq 0
        then
                echo " `date` now,starting server!">>  restart.log
                echo 'start'
                nohup ./run.sh > log  2>&1 &
        fi
        sleep 5
done
