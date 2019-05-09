#!/bin/bash

pid=`ps -ef | grep ghpb | head -n 1 | awk '{ print $2}'`
for i in {1..100}
do
    statu=`cat /proc/${pid}/status | grep -E "VmRSS|Threads"`
    date=`date`
    echo "$date"  >> memstat.log
    echo "$statu" >> memstat.log
    sleep 2
done

