#! /bin/sh
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
cp FlameGraph/flamegraph.pl /usr/local/bin
echo "Asia/Shanghai" > /etc/timezone
export GOPATH=$GOPATH:/work
while [ 1 = 1 ]
do
sleep 1;
done
