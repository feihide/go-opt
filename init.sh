#! /bin/sh
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
echo "Asia/Shanghai" > /etc/timezone
while [ 1 = 1 ]
do
sleep 1;
done
